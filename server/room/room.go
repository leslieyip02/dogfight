package room

import (
	"context"
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

	return &room, nil
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: add capacity check
	r.clients[client.id] = client
}

func (r *Room) Remove(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clients, client.id)
}

func (r *Room) Connect(client *Client) error {
	go client.readPump(r.game.Incoming)
	go client.writePump(r.game.Outgoing)
	return r.game.AddPlayer(client.id, client.username)
}

func (r *Room) Stop() {
	r.cancel()
}
