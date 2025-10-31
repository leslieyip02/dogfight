package room

import (
	"context"
	"server/game"
	"server/id"
	"server/pb"

	"sync"

	"google.golang.org/protobuf/proto"
)

const ROOM_CAPACITY = 32

type Room struct {
	id      string
	game    *game.Game
	clients map[string]*Client
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewRoom() (*Room, error) {
	id, err := id.NewShortId()
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

	r.clients[client.id] = client
}

func (r *Room) Remove(clientId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.game.RemovePlayer(clientId)
	delete(r.clients, clientId)
}

func (r *Room) Stop() {
	r.cancel()
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

			var event pb.Event
			proto.Unmarshal(message, &event)

			switch event.Type {
			case pb.EventType_EVENT_TYPE_QUIT:
				data := event.GetQuitEventData()
				r.Remove(data.GetId())
			}
		}
	}
}

func (r *Room) hasCapacity() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.clients) < ROOM_CAPACITY
}
