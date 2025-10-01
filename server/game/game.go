package game

import (
	"context"
	"encoding/json"
	"server/game/entities"
	"sync"
	"time"
)

type Game struct {
	Incoming chan []byte
	Outgoing chan []byte
	mu       sync.Mutex

	// state
	entities     map[string]entities.Entity
	usernames    map[string]string
	frameCounter int64

	// state delta
	updated map[string]entities.Entity
	removed []string
}

func NewGame() Game {
	return Game{
		Incoming:     make(chan []byte),
		Outgoing:     make(chan []byte),
		mu:           sync.Mutex{},
		entities:     make(map[string]entities.Entity),
		usernames:    map[string]string{},
		frameCounter: 0,
		updated:      make(map[string]entities.Entity),
		removed:      []string{},
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

func (g *Game) remove(id string) {
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
				case QuitEventType:
					var data QuitEventData
					json.Unmarshal(event.Data, &data)
					g.remove(data.ID)

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

	// TODO: move this
	g.frameCounter++
	if g.frameCounter%1800 == 0 {
		g.addPowerup()
		g.frameCounter = 0
	}

	g.updateEntities()
	g.resolveCollisions()
	for _, id := range g.removed {
		delete(g.entities, id)
		delete(g.updated, id)
	}
	g.broadcast()

	clear(g.updated)
	g.removed = g.removed[:0]
}

func (g *Game) updateEntities() {
	newEntities := []entities.Entity{}
	for id, entity := range g.entities {
		if entity.Update() {
			g.updated[id] = entity
		}

		if entity.GetIsExpired() {
			g.removed = append(g.removed, id)
		}

		newEntities = append(newEntities, entity.PollNewEntities()...)
	}

	for _, newEntity := range newEntities {
		g.entities[newEntity.GetID()] = newEntity
	}
}

func (g *Game) broadcast() {
	message, err := CreateMessage(DeltaEventType, g.GetDelta())
	if err != nil {
		return
	}
	g.Outgoing <- message
}

func (g *Game) resolveCollisions() {
	// TODO: use line sweep to lower time complexity to O(n log(n))
	for i, entityA := range g.entities {
		for j, entityB := range g.entities {
			if i >= j {
				continue
			}

			if g.checkCollision(entityA, entityB) {
				g.handleCollision(entityA, entityB)
			}
		}
	}
}

func (g *Game) checkCollision(a entities.Entity, b entities.Entity) bool {
	return a.GetBoundingBox().DidCollide(b.GetBoundingBox())
}

func (g *Game) handleCollision(a entities.Entity, b entities.Entity) {
	// TODO: replace with something more robust
	typeA := a.GetType()
	typeB := b.GetType()

	switch {
	case typeA == entities.PlayerEntityType && typeB == entities.PlayerEntityType:
	case typeA == entities.PlayerEntityType && typeB == entities.ProjectileEntityType:
	case typeA == entities.ProjectileEntityType && typeB == entities.PlayerEntityType:
		g.removed = append(g.removed, a.GetID())
		g.removed = append(g.removed, b.GetID())

	case typeA == entities.PlayerEntityType && typeB == entities.PowerupEntityType:
		player := a.(*entities.Player)
		powerup := b.(*entities.Powerup)
		player.Powerup = powerup
		g.removed = append(g.removed, b.GetID())

	case typeA == entities.PowerupEntityType && typeB == entities.PlayerEntityType:
		powerup := a.(*entities.Powerup)
		player := b.(*entities.Player)
		player.Powerup = powerup
		g.removed = append(g.removed, a.GetID())
	}
}

func (g *Game) addPowerup() error {
	powerup, err := entities.NewPowerup(entities.MultishotPowerupAbility)
	if err != nil {
		return err
	}

	g.entities[powerup.ID] = powerup
	g.updated[powerup.ID] = powerup
	return nil
}
