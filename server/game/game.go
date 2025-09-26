package game

import (
	"context"
	"encoding/json"
	"math"
	"server/utils"
	"time"
)

type Game struct {
	Incoming chan []byte
	Outgoing chan []byte
	entities map[string]Entity
}

type GameState struct {
	timestamp int64
	entities  map[string]Entity
}

func NewGame() Game {
	return Game{
		Incoming: make(chan []byte),
		Outgoing: make(chan []byte),
		entities: make(map[string]Entity),
	}
}

func (g *Game) AddPlayer(id string, username string) error {
	player := Player{
		ID:       id,
		Username: username,
		Position: randomEntityPosition(),
		speed:    MAX_PLAYER_SPEED,
		powerup:  nil,
	}

	g.entities[id] = &player
	message, err := NewJoinEventMessage(&player)
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
}

func (g *Game) GetState() GameState {
	return GameState{
		timestamp: time.Now().UnixNano(),
		entities:  g.entities,
	}
}

func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(FRAME_DURATION)

	go func() {
		defer ticker.Stop()

		var frameCounter = 0

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				frameCounter++
				if frameCounter%POWERUP_SPAWN_INTERVAL == 0 {
					g.addPowerup()
					frameCounter = 0
				}

				g.update()

			case message := <-g.Incoming:
				var event Event
				json.Unmarshal(message, &event)

				switch event.Type {
				case QuitEventType:
					var data QuitEventData
					json.Unmarshal(event.Data, &data)
					delete(g.entities, data.Id)

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
	entity, found := g.entities[data.ClientId]
	if !found {
		return
	}

	if player, ok := entity.(*Player); ok {
		player.input(data)
	}
}

func (g *Game) update() {
	expiredIDs := []string{}
	for id, entity := range g.entities {
		entity.Update(g)
		if entity.GetIsExpired() {
			expiredIDs = append(expiredIDs, id)
		}
	}
	for _, id := range expiredIDs {
		delete(g.entities, id)
	}

	g.resolveCollisions()
	g.broadcast()
}

func (g *Game) broadcast() {
	// TODO: just update all entities
	players := make(map[string]*Player)
	projectiles := make(map[string]*Projectile)
	for _, entity := range g.entities {
		if player, ok := entity.(*Player); ok {
			players[player.ID] = player
		}
		if projectile, ok := entity.(*Projectile); ok {
			projectiles[projectile.ID] = projectile
		}
	}

	message, err := NewUpdatePositionEventMessage(&players, &projectiles)
	if err != nil {
		return
	}
	g.Outgoing <- message
}

func (g *Game) resolveCollisions() {
	// TODO: use line sweep to lower time complexity to O(n log(n))
	toRemove := make(map[string]bool)
	for i, entityA := range g.entities {
		for j, entityB := range g.entities {
			if i >= j {
				continue
			}

			if g.checkCollision(entityA, entityB) {
				g.handleCollision(entityA, entityB, toRemove)
			}
		}
	}

	for id := range toRemove {
		delete(g.entities, id)
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

func (g *Game) handleCollision(a Entity, b Entity, toRemove map[string]bool) {
	// TODO: replace with something more robust
	typeA := a.GetType()
	typeB := b.GetType()

	switch {
	case typeA == PlayerEntityType && typeB == ProjectileEntityType:
		toRemove[a.GetID()] = true
		toRemove[b.GetID()] = true

	case typeA == ProjectileEntityType && typeB == PlayerEntityType:
		toRemove[a.GetID()] = true
		toRemove[b.GetID()] = true

	case typeA == PlayerEntityType && typeB == PowerupEntityType:
		player := a.(*Player)
		powerup := b.(*Powerup)
		player.powerup = powerup
		toRemove[b.GetID()] = true

	case typeA == PowerupEntityType && typeB == PlayerEntityType:
		powerup := a.(*Powerup)
		player := b.(*Player)
		player.powerup = powerup
		toRemove[a.GetID()] = true

	case typeA == PlayerEntityType && typeB == PlayerEntityType:
		toRemove[a.GetID()] = true
		toRemove[b.GetID()] = true
	}
}

func (g *Game) addPowerup() error {
	id, err := utils.NewShortId()
	if err != nil {
		return err
	}

	powerup := Powerup{
		ID:       id,
		Type:     "multishot",
		Position: randomEntityPosition(),
	}
	g.entities[id] = &powerup
	message, err := NewUpdatePowerupEventMessage(&powerup, true)
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
}
