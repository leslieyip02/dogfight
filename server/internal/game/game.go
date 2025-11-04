package game

import (
	"context"
	"log"
	"server/internal/game/collision"
	"server/internal/game/constants"
	"server/internal/game/entities"
	"server/pb"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	MAX_ENTITY_COUNT = 256
)

// A Game stores the game's state and handles its logic.
type Game struct {
	Incoming chan []byte
	Outgoing chan []byte
	mu       sync.Mutex

	// Game tate.
	entities  map[string]entities.Entity
	usernames map[string]string
	spawner   entities.Spawner

	// Game state deltas.
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

// AddPlayer spawns a new Player into the game.
func (g *Game) AddPlayer(id string, username string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, err := g.spawner.SpawnPlayer(id, username)
	if err != nil {
		return err
	}

	g.entities[id] = player
	g.usernames[id] = username
	return nil
}

// RemovePlayer removes a Player from the game.
func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.removed = append(g.removed, id)
	delete(g.usernames, id)
}

// respawnPlayer adds a new Player into the game, checking if id is already
// being used and looks up a username based on id.
func (g *Game) respawnPlayer(id string) {
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
		// TODO: handle error
		log.Fatalf("could not spawn player")
	}
	g.entities[id] = player
}

// GetPbEntities unwraps the game's entities into their underlying EntityData
// for serialization.
func (g *Game) GetPbEntities() []*pb.EntityData {
	entities := make([]*pb.EntityData, len(g.entities))
	i := 0
	for _, entity := range g.entities {
		entities[i] = entity.GetEntityData()
		i++
	}
	return entities
}

// GetTimestamp returns the timestamp as a float. The cast is necessary because
// JavaScript's Number.MAX_SAFE_INTEGER can't handle int64.
func (g *Game) GetTimestamp() float64 {
	return float64(time.Now().UnixMilli())
}

// GetSnapshot serializes the game state.
func (g *Game) GetSnapshot() *pb.Event {
	return &pb.Event{
		Type: pb.EventType_EVENT_TYPE_SNAPSHOT,
		Data: &pb.Event_SnapshotEventData_{
			SnapshotEventData: &pb.Event_SnapshotEventData{
				Timestamp: g.GetTimestamp(),
				Entities:  g.GetPbEntities(),
			},
		},
	}
}

// GetDelta serializes the changes in the game state between each call to
// GetDelta.
func (g *Game) GetDelta() *pb.Event {
	return &pb.Event{
		Type: pb.EventType_EVENT_TYPE_DELTA,
		Data: &pb.Event_DeltaEventData_{
			DeltaEventData: &pb.Event_DeltaEventData{
				Timestamp: g.GetTimestamp(),
				Updated:   g.GetPbEntities(),
				Removed:   g.removed,
			},
		},
	}
}

// Run starts the game loop.
func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second / constants.FPS)
	defer ticker.Stop()

	for _, entity := range g.spawner.InitEntities() {
		g.entities[entity.GetId()] = entity
		g.updated[entity.GetId()] = entity
	}

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
				g.respawnPlayer(data.GetId())

			case pb.EventType_EVENT_TYPE_INPUT:
				data := event.GetInputEventData()
				g.input(data)
			}
		}
	}
}

// input passes input event data to the corresponding Player.
func (g *Game) input(data *pb.Event_InputEventData) {
	entity, found := g.entities[data.GetId()]
	if !found {
		return
	}

	if player, ok := entity.(*entities.Player); ok {
		player.Input(data.GetMouseX(), data.GetMouseY(), data.MousePressed)
	}
}

// update is called once per tick and computes all updates.
//
// More specifically, it
//   - updates positions
//   - resolves collisions
//   - adds new entities
//   - removes expired entities
//   - broacasts the updated delta
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

	data := g.GetDelta()
	g.broadcast(data)

	clear(g.updated)
	g.removed = g.removed[:0]
}

// updateEntities updates entities, and marks expired entities for deletion.
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

// resolveCollisions checks and handles collisions for all entities.
func (g *Game) resolveCollisions() {
	collision.ResolveCollisionsLineSweep(&g.entities, g.handleCollision)
}

// handleCollision updates the entities with id1 and id2 and marks them for
// removal if needed.
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

// pollNewEntities polls all new entities that have been created and adds them
// into the game.
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

// broadcast sends a message to the room.
//
// TODO: maybe this shouldn't be here
func (g *Game) broadcast(data *pb.Event) error {
	message, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	g.Outgoing <- message
	return nil
}
