package room

import (
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Manager struct {
	rooms   map[string]*Room
	roomIds []string
	mu      sync.Mutex
}

type JoinRequest struct {
	Username string  `json:"username"`
	RoomId   *string `json:"roomId,omitempty"`
}

func NewManager() (*Manager, error) {
	roomManager := Manager{
		rooms:   map[string]*Room{},
		roomIds: []string{},
		mu:      sync.Mutex{},
	}

	// TODO: handle adding more rooms
	room, err := NewRoom()
	if err != nil {
		return nil, err
	}

	roomManager.rooms[room.id] = room
	roomManager.roomIds = append(roomManager.roomIds, room.id)
	return &roomManager, nil
}

func (m *Manager) getRoom(roomId *string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	if roomId != nil {
		return m.rooms[*roomId]
	}

	randomId := m.roomIds[rand.Intn(len(m.roomIds))]
	return m.rooms[randomId]
}
