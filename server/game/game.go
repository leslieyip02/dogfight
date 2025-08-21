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

	distance := math.Sqrt(data.MouseX*data.MouseX + data.MouseY*data.MouseY)
	if distance == 0 {
		return
	}

	// TODO: move into player
	speed := 4.0
	player.x += (data.MouseX / distance) * speed
	player.y += (data.MouseY / distance) * speed
	player.theta = math.Atan2(data.MouseY, data.MouseX)
}

func (g *Game) update() {
	message, err := NewUpdateEventMessage(g)
	if err != nil {
		return
	}
	g.Outgoing <- message
}
