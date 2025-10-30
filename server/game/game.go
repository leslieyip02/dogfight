package game

import (
	"context"
	"encoding/json"
	"server/game/collision"
	"server/game/entities"
	"sync"
	"time"
)

const (
	MAX_ENTITY_COUNT = 256
)

type Game struct {
	Incoming chan []byte
	Outgoing chan []byte
	mu       sync.Mutex

	// state
	entities  map[string]entities.Entity
	usernames map[string]string
	spawner   entities.Spawner

	// state delta
	updated map[string]entities.Entity
	removed []string
}

func NewGame() *Game {
	return &Game{
		Incoming:  make(chan []byte),
		Outgoing:  make(chan []byte),
		mu:        sync.Mutex{},
		entities:  make(map[string]entities.Entity),
		usernames: map[string]string{},
		spawner:   entities.NewSpawner(),
		updated:   make(map[string]entities.Entity),
		removed:   []string{},
	}
}

func (g *Game) AddPlayer(id string, username string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.entities[id] = entities.NewPlayer(id, username)
	g.usernames[id] = username

	message, err := CreateMessage(JoinEventType, JoinEventData{
		ID:       id,
		Username: username,
	})
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
}

func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.removed = append(g.removed, id)
	delete(g.usernames, id)
}

func (g *Game) respawn(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, found := g.entities[id]
	if found {
		return
	}

	username, found := g.usernames[id]
	if !found {
		return
	}

	g.entities[id] = entities.NewPlayer(id, username)
}

func (g *Game) GetSnapshot() SnapshotEventData {
	return SnapshotEventData{
		Timestamp: time.Now().UnixNano(),
		Entities:  g.entities,
	}
}

func (g *Game) GetDelta() DeltaEventData {
	return DeltaEventData{
		Timestamp: time.Now().UnixNano(),
		Updated:   g.updated,
		Removed:   g.removed,
	}
}

func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second / entities.FPS)

	for _, entity := range g.spawner.InitEntities() {
		g.entities[entity.GetID()] = entity
		g.updated[entity.GetID()] = entity
	}

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				g.update()

			case message := <-g.Incoming:
				g.Outgoing <- message

				var event Event
				json.Unmarshal(message, &event)

				switch event.Type {
				case RespawnEventType:
					var data RespawnEventData
					json.Unmarshal(event.Data, &data)
					g.respawn(data.ID)

				case InputEventType:
					var data InputEventData
					json.Unmarshal(event.Data, &data)
					g.input(data)
				}
			}
		}
	}()
}

func (g *Game) input(data InputEventData) {
	entity, found := g.entities[data.ID]
	if !found {
		return
	}

	if player, ok := entity.(*entities.Player); ok {
		player.Input(data.MouseX, data.MouseY, data.MousePressed)
	}
}

func (g *Game) update() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.updateEntities()
	g.resolveCollisions()
	g.pollNewEntities()
	for _, id := range g.removed {
		delete(g.entities, id)
		delete(g.updated, id)
	}
	g.broadcast()

	clear(g.updated)
	g.removed = g.removed[:0]
}

func (g *Game) updateEntities() {
	for id, entity := range g.entities {
		if entity.Update() {
			g.updated[id] = entity
		}

		if entity.GetIsExpired() {
			g.removed = append(g.removed, id)
		}
	}
}

func (g *Game) resolveCollisions() {
	collision.ResolveCollisions(&g.entities, g.handleCollision)
}

func (g *Game) handleCollision(id1 *string, id2 *string) {
	if id1 == id2 {
		return
	}

	e1 := g.entities[*id1]
	e2 := g.entities[*id2]
	e1.UpdateOnCollision(e2)
	e2.UpdateOnCollision(e1)

	if e1.RemoveOnCollision(e2) {
		g.removed = append(g.removed, e1.GetID())
	}
	if e2.RemoveOnCollision(e1) {
		g.removed = append(g.removed, e2.GetID())
	}
}

func (g *Game) pollNewEntities() {
	for _, entity := range g.entities {
		for _, newEntity := range entity.PollNewEntities() {
			g.entities[newEntity.GetID()] = newEntity
			g.updated[newEntity.GetID()] = newEntity
		}
	}

	if len(g.entities) > MAX_ENTITY_COUNT {
		return
	}
	for _, newEntity := range g.spawner.PollNewEntities() {
		g.entities[newEntity.GetID()] = newEntity
		g.updated[newEntity.GetID()] = newEntity
	}
}

func (g *Game) broadcast() {
	message, err := CreateMessage(DeltaEventType, g.GetDelta())
	if err != nil {
		return
	}
	g.Outgoing <- message
}
