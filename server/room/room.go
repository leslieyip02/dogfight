package room

import (
	"context"
	"encoding/json"
	"server/game"
	"server/utils"

	"sync"
)

type Room struct {
	id      string
	game    game.Game
	clients map[string]*Client
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewRoom() (*Room, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	game := game.NewGame()
	room := Room{
		id:      id,
		game:    game,
		clients: map[string]*Client{},
		mu:      sync.Mutex{},
		ctx:     ctx,
		cancel:  cancel,
	}

	game.Run(ctx)
	go room.broadcast()

	return &room, nil
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: add capacity check
	r.clients[client.id] = client
}

func (r *Room) Remove(clientId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clients, clientId)
}

func (r *Room) connect(client *Client) error {
	go client.readPump(r.game.Incoming)
	go client.writePump()
	return r.game.AddPlayer(client.id, client.username)
}

func (r *Room) broadcast() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case message := <-r.game.Outgoing:
			for _, client := range r.clients {
				client.send <- message
			}

			var event game.Event
			json.Unmarshal(message, &event)
			switch event.Type {
			case game.QuitEventType:
				var data game.QuitEventData
				json.Unmarshal(event.Data, &data)

				r.Remove(data.Id)
			}
		}
	}
}

func (r *Room) getClient(id string) (*Client, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.clients[id]
	return c, ok
}

func (r *Room) Stop() {
	r.cancel()
}
