package room

import (
	"context"
	"server/internal/game"
	"server/pb"

	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const ROOM_CAPACITY = 32

// A Room allows multiple clients to connect, and runs a single game instance.
type Room struct {
	id      string
	game    *game.Game
	clients map[string]*Client

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func newRoom(id string) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	game := game.NewGame()

	return &Room{
		id:      id,
		game:    game,
		clients: map[string]*Client{},
		mu:      sync.Mutex{},
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (r *Room) InitClient(clientId string, username string, conn *websocket.Conn) {
	conn.SetCloseHandler(func(code int, text string) error {
		r.remove(clientId)
		return nil
	})

	client := newClient(clientId, username, conn)
	r.add(client)
	r.connect(client)
}

func (r *Room) init() {
	go r.game.Run(r.ctx)
	go r.broadcast()
}

func (r *Room) add(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[client.id] = client
}

func (r *Room) remove(clientId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.game.RemovePlayer(clientId)
	delete(r.clients, clientId)
}

func (r *Room) stop() {
	r.cancel()
}

// connect allows clients to connect to the room, and sends a join event
// message to other clients.
func (r *Room) connect(client *Client) error {
	go client.readPump(r.game.Incoming)
	go client.writePump()

	err := r.game.AddPlayer(client.id, client.username)
	if err != nil {
		return err
	}
	return r.sendJoinEvent(client.id, client.username)
}

// broadcast relays messages from the game to all clients.
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
				r.remove(data.GetId())
			}
		}
	}
}

func (r *Room) sendJoinEvent(id string, username string) error {
	data := &pb.Event{
		Type: pb.EventType_EVENT_TYPE_JOIN,
		Data: &pb.Event_JoinEventData_{
			JoinEventData: &pb.Event_JoinEventData{
				Id:       id,
				Username: username,
			},
		},
	}

	message, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	r.game.Outgoing <- message
	return nil
}

// hasCapacity reports if more players could be added to the room.
func (r *Room) hasCapacity() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.clients) < ROOM_CAPACITY
}

// getOccupancy returns the current number of connected players
func (r *Room) getOccupancy() uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()

	return uint32(len(r.clients))
}
