package game

import (
	"context"
	"encoding/json"
	"math"
	"time"
)

type Game struct {
	// messages
	Incoming chan []byte
	Outgoing chan []byte

	// state
	entities     map[string]Entity
	frameCounter int64

	// state delta
	updated map[string]Entity
	removed []string
}

func NewGame() Game {
	return Game{
		Incoming:     make(chan []byte),
		Outgoing:     make(chan []byte),
		entities:     make(map[string]Entity),
		frameCounter: 0,
		updated:      make(map[string]Entity),
		removed:      []string{},
	}
}

func (g *Game) AddPlayer(id string, username string) error {

	g.entities[id] = NewPlayer(id, username)
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

func (g *Game) GetSnapshot() SnapshotEventData {
	return SnapshotEventData{
		Entities: g.entities,
	}
}

func (g *Game) GetDelta() DeltaEventData {
	return DeltaEventData{
		Updated: g.updated,
		Removed: g.removed,
	}
}

func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(FRAME_DURATION)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				g.update()

			case message := <-g.Incoming:
				var event Event
				json.Unmarshal(message, &event)

				switch event.Type {
				case QuitEventType:
					var data QuitEventData
					json.Unmarshal(event.Data, &data)
					delete(g.entities, data.ID)

					g.Outgoing <- message

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

	if player, ok := entity.(*Player); ok {
		player.input(data)
	}
}

func (g *Game) update() {
	clear(g.updated)
	clear(g.removed)

	// TODO: move this?
	g.frameCounter++
	if g.frameCounter%POWERUP_SPAWN_INTERVAL == 0 {
		g.addPowerup()
		g.frameCounter = 0
	}

	g.updateEntities()
	g.resolveCollisions()
	for _, id := range g.removed {
		delete(g.entities, id)
	}

	g.broadcast()
}

func (g *Game) updateEntities() {
	for id, entity := range g.entities {
		if entity.Update(g) {
			g.updated[id] = entity
		}

		if entity.GetIsExpired() {
			g.removed = append(g.removed, id)
		}
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

func (g *Game) checkCollision(a Entity, b Entity) bool {
	// TODO: replace with something more robust
	dx := a.GetPosition().X - b.GetPosition().X
	dy := a.GetPosition().Y - b.GetPosition().Y
	distance := math.Sqrt(dx*dx + dy*dy)

	threshold := 0.0
	if a.GetType() == PlayerEntityType {
		threshold += PLAYER_RADIUS
	} else if a.GetType() == ProjectileEntityType {
		threshold += PROJECTILE_RADIUS
	} else if a.GetID() == string(PowerupEntityType) {
		threshold += PROJECTILE_RADIUS
	}
	if b.GetType() == PlayerEntityType {
		threshold += PLAYER_RADIUS
	} else if b.GetType() == ProjectileEntityType {
		threshold += PROJECTILE_RADIUS
	} else if b.GetID() == string(PowerupEntityType) {
		threshold += PROJECTILE_RADIUS
	}
	return distance <= threshold
}

func (g *Game) handleCollision(a Entity, b Entity) {
	// TODO: replace with something more robust
	typeA := a.GetType()
	typeB := b.GetType()

	switch {
	case typeA == PlayerEntityType && typeB == PlayerEntityType:
	case typeA == PlayerEntityType && typeB == ProjectileEntityType:
	case typeA == ProjectileEntityType && typeB == PlayerEntityType:
		g.removed = append(g.removed, a.GetID())
		g.removed = append(g.removed, b.GetID())

	case typeA == PlayerEntityType && typeB == PowerupEntityType:
		player := a.(*Player)
		powerup := b.(*Powerup)
		player.powerup = powerup
		g.removed = append(g.removed, b.GetID())

	case typeA == PowerupEntityType && typeB == PlayerEntityType:
		powerup := a.(*Powerup)
		player := b.(*Player)
		player.powerup = powerup
		g.removed = append(g.removed, a.GetID())
	}
}

func (g *Game) addPowerup() error {
	powerup, err := NewPowerup(MultishotPowerupType)
	if err != nil {
		return err
	}

	g.entities[powerup.ID] = powerup
	g.updated[powerup.ID] = powerup
	return nil
}
