package game

import (
	"context"
	"encoding/json"
	"math"
	"time"
)

type Game struct {
	Incoming    chan []byte
	Outgoing    chan []byte
	players     map[string]*Player
	projectiles map[string]*Projectile
}

func NewGame() Game {
	return Game{
		Incoming:    make(chan []byte),
		Outgoing:    make(chan []byte),
		players:     map[string]*Player{},
		projectiles: map[string]*Projectile{},
	}
}

func (g *Game) AddPlayer(id string, username string) error {
	position := EntityPosition{
		// X:     rand.Float64()*WIDTH - WIDTH/2,
		// Y:     rand.Float64()*HEIGHT - WIDTH/2,
		X:     0,
		Y:     0,
		Theta: math.Pi / 2,
	}
	player := Player{
		Id:       id,
		Username: username,
		Position: position,
		speed:    MAX_SPEED,
	}

	g.players[id] = &player
	message, err := NewJoinEventMessage(&player)
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
}

func (g *Game) GetPlayers() []*Player {
	players := []*Player{}
	for _, player := range g.players {
		players = append(players, player)
	}
	return players
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

	projectile := player.input(data)
	if projectile != nil {
		g.projectiles[projectile.Id] = projectile
	}
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
	for i, player := range g.players {
		for j, projectile := range g.projectiles {
			// projectiles are modelled as circles
			dx := player.Position.X - projectile.position.X
			dy := player.Position.Y - projectile.position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance <= PLAYER_RADIUS+PROJECTILE_RADIUS {
				collidedPlayerIds = append(collidedPlayerIds, i)
				collidedProjectileIds = append(collidedProjectileIds, j)
			}
		}
	}

	for _, id := range collidedPlayerIds {
		delete(g.players, id)
	}
	for _, id := range collidedProjectileIds {
		delete(g.projectiles, id)
	}
}
