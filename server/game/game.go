package game

import (
	"context"
	"encoding/json"
	"math"
	"server/utils"
	"time"
)

type Game struct {
	Incoming    chan []byte
	Outgoing    chan []byte
	players     map[string]*Player
	projectiles map[string]*Projectile
	powerups    map[string]*Powerup
}

type GameState struct {
	Players  []*Player  `json:"players"`
	Powerups []*Powerup `json:"powerups"`
}

func NewGame() Game {
	return Game{
		Incoming:    make(chan []byte),
		Outgoing:    make(chan []byte),
		players:     map[string]*Player{},
		projectiles: map[string]*Projectile{},
		powerups:    map[string]*Powerup{},
	}
}

func (g *Game) AddPlayer(id string, username string) error {
	player := Player{
		Id:       id,
		Username: username,
		Position: randomEntityPosition(),
		speed:    MAX_PLAYER_SPEED,
		powerup:  nil,
	}

	g.players[id] = &player
	message, err := NewJoinEventMessage(&player)
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
}

func (g *Game) GetState() GameState {
	players := []*Player{}
	for _, player := range g.players {
		players = append(players, player)
	}
	powerups := []*Powerup{}
	for _, powerup := range g.powerups {
		powerups = append(powerups, powerup)
	}
	return GameState{
		Players:  players,
		Powerups: powerups,
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
					delete(g.players, data.Id)

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
	player, found := g.players[data.ClientId]
	if !found {
		return
	}
	player.input(data, g)
}

func (g *Game) update() {
	g.updateProjectiles()
	g.resolveCollisions()

	message, err := NewUpdatePositionEventMessage(&g.players, &g.projectiles)
	if err != nil {
		return
	}
	g.Outgoing <- message
}

func (g *Game) updateProjectiles() {
	expiredIds := []string{}
	for _, projectile := range g.projectiles {
		projectile.update()
		if projectile.lifetime < 0 {
			expiredIds = append(expiredIds, projectile.Id)
		}
	}

	for _, id := range expiredIds {
		delete(g.projectiles, id)
	}
}

func (g *Game) resolveCollisions() {
	// TODO: use line sweep to lower time complexity to O(n log(n))
	collidedPlayerIds := []string{}
	for i, player := range g.players {
		for j, other := range g.players {
			if i == j {
				continue
			}

			// players are modelled as circles
			dx := player.Position.X - other.Position.X
			dy := player.Position.Y - other.Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance <= 2*PLAYER_RADIUS {
				collidedPlayerIds = append(collidedPlayerIds, i, j)
			}
		}
	}

	collidedProjectileIds := []string{}
	consumedPowerupIds := []string{}
	for i, player := range g.players {
		for j, projectile := range g.projectiles {
			// projectiles are modelled as circles
			dx := player.Position.X - projectile.Position.X
			dy := player.Position.Y - projectile.Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance <= PLAYER_RADIUS+PROJECTILE_RADIUS {
				collidedPlayerIds = append(collidedPlayerIds, i)
				collidedProjectileIds = append(collidedProjectileIds, j)
			}
		}

		for j, powerup := range g.powerups {
			dx := player.Position.X - powerup.Position.X
			dy := player.Position.Y - powerup.Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance <= PLAYER_RADIUS+PROJECTILE_RADIUS {
				player.powerup = powerup
				consumedPowerupIds = append(consumedPowerupIds, j)
			}
		}
	}

	for _, id := range collidedPlayerIds {
		delete(g.players, id)
	}
	for _, id := range collidedProjectileIds {
		delete(g.projectiles, id)
	}
	for _, id := range consumedPowerupIds {
		message, err := NewUpdatePowerupEventMessage(g.powerups[id], false)
		if err == nil {
			g.Outgoing <- message
		}
		delete(g.powerups, id)
	}
}

func (g *Game) addPowerup() error {
	id, err := utils.NewShortId()
	if err != nil {
		return err
	}

	powerup := Powerup{
		Id:       id,
		Type:     "multishot",
		Position: randomEntityPosition(),
	}
	g.powerups[id] = &powerup
	message, err := NewUpdatePowerupEventMessage(&powerup, true)
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
}
