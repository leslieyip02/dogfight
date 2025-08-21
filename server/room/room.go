package room

import (
	"server/game"
	"server/utils"

	"sync"
)

type Room struct {
	id      string
	game    game.Game
	clients map[string]*Client
	mu      sync.Mutex
}

func NewRoom() (*Room, error) {
	id, err := utils.GetShortId()
	if err != nil {
		return nil, err
	}

	game := game.NewGame()
	go game.Run()

	room := Room{
		id:      id,
		game:    game,
		clients: map[string]*Client{},
		mu:      sync.Mutex{},
	}
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
	go client.readPump(r.game.Send)
	go client.writePump(r.game.Broadcast)

	joinEventMessage, err := game.NewJoinEventMessage(client.id, client.username)
	if err != nil {
		return err
	}

	r.game.Send <- joinEventMessage
	return nil
}
