package game

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"time"
)

type Game struct {
	Incoming chan []byte
	Outgoing chan []byte
	players  map[string]*Player
}

func NewGame() Game {
	return Game{
		Incoming: make(chan []byte),
		Outgoing: make(chan []byte),
		players:  map[string]*Player{},
	}
}

func (g *Game) AddPlayer(id string, username string) error {
	player := Player{
		Id:       id,
		Username: username,
		x:        rand.Float64()*WIDTH - WIDTH/2,
		y:        rand.Float64()*HEIGHT - WIDTH/2,
		theta:    math.Pi / 2,
		speed:    MAX_VELOCITY,
	}

	g.players[id] = &player
	message, err := NewJoinEventMessage(&player)
	if err != nil {
		return err
	}
	g.Outgoing <- message
	return nil
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
					delete(g.players, data.ClientId)

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

	delta := normalizeAngle(math.Atan2(data.MouseY, data.MouseX) - player.theta)
	player.theta = normalizeAngle(player.theta + delta*0.1)

	// TODO: consider non-linear multiplier (e.g. -(x - 1)^2 + 1)
	length := math.Sqrt(data.MouseX*data.MouseX + data.MouseY*data.MouseY)
	player.x += math.Cos(player.theta) * length * player.speed
	player.y += math.Sin(player.theta) * length * player.speed
}

func (g *Game) update() {
	message, err := NewUpdateEventMessage(g)
	if err != nil {
		return
	}
	g.Outgoing <- message
}

func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle > math.Pi {
		angle -= 2 * math.Pi
	} else if angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}
