package game

import (
	"context"
	"log"
	"server/game/collision"
	"server/game/entities"
	"server/pb"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
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

	player, err := g.spawner.SpawnPlayer(id, username)
	if err != nil {
		log.Fatalf("could not spawn player")
	}
	g.entities[id] = player
	g.usernames[id] = username

	data := &pb.Event{
		Type: pb.EventType_EVENT_TYPE_JOIN,
		Data: &pb.Event_JoinEventData_{
			JoinEventData: &pb.Event_JoinEventData{
				Id:       id,
				Username: username,
			},
		},
	}
	message, err := proto.Marshal(data)
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

	player, err := g.spawner.SpawnPlayer(id, username)
	if err != nil {
		log.Fatalf("could not spawn player")
	}
	g.entities[id] = player
}

func (g *Game) GetPBEntities() []*pb.EntityData {
	entities := make([]*pb.EntityData, len(g.entities))
	i := 0
	for _, entity := range g.entities {
		entities[i] = entity.GetEntityData()
		i++
	}
	return entities
}

func (g *Game) GetTimestamp() float64 {
	// JavaScript's Number.MAX_SAFE_INTEGER can't handle int64
	return float64(time.Now().UnixMilli())
}

func (g *Game) GetSnapshot() *pb.Event {
	return &pb.Event{
		Type: pb.EventType_EVENT_TYPE_SNAPSHOT,
		Data: &pb.Event_SnapshotEventData_{
			SnapshotEventData: &pb.Event_SnapshotEventData{
				Timestamp: g.GetTimestamp(),
				Entities:  g.GetPBEntities(),
			},
		},
	}
}

func (g *Game) GetDelta() *pb.Event {
	return &pb.Event{
		Type: pb.EventType_EVENT_TYPE_DELTA,
		Data: &pb.Event_DeltaEventData_{
			DeltaEventData: &pb.Event_DeltaEventData{
				Timestamp: g.GetTimestamp(),
				Updated:   g.GetPBEntities(),
				Removed:   g.removed,
			},
		},
	}
}

func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second / entities.FPS)

	for _, entity := range g.spawner.InitEntities() {
		g.entities[entity.GetId()] = entity
		g.updated[entity.GetId()] = entity
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

				var event pb.Event
				proto.Unmarshal(message, &event)

				switch event.Type {
				case pb.EventType_EVENT_TYPE_RESPAWN:
					data := event.GetRespawnEventData()
					g.respawn(data.GetId())

				case pb.EventType_EVENT_TYPE_INPUT:
					data := event.GetInputEventData()
					g.input(data)
				}
			}
		}
	}()
}

func (g *Game) input(data *pb.Event_InputEventData) {
	entity, found := g.entities[data.GetId()]
	if !found {
		return
	}

	if player, ok := entity.(*entities.Player); ok {
		player.Input(data.GetMouseX(), data.GetMouseY(), data.MousePressed)
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
	collision.ResolveCollisionsLineSweep(&g.entities, g.handleCollision)
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
		g.removed = append(g.removed, e1.GetId())
	}
	if e2.RemoveOnCollision(e1) {
		g.removed = append(g.removed, e2.GetId())
	}
}

func (g *Game) pollNewEntities() {
	for _, entity := range g.entities {
		for _, newEntity := range entity.PollNewEntities() {
			g.entities[newEntity.GetId()] = newEntity
			g.updated[newEntity.GetId()] = newEntity
		}
	}

	if len(g.entities) > MAX_ENTITY_COUNT {
		return
	}
	for _, newEntity := range g.spawner.PollNewEntities() {
		g.entities[newEntity.GetId()] = newEntity
		g.updated[newEntity.GetId()] = newEntity
	}
}

func (g *Game) broadcast() {
	message, err := proto.Marshal(g.GetDelta())
	if err != nil {
		log.Printf("broadcast failed: %v", err)
		return
	}
	g.Outgoing <- message
}
