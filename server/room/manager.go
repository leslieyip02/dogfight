package room

import (
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

func NewManager() *Manager {
	return &Manager{
		rooms:   map[string]*Room{},
		roomIds: []string{},
		mu:      sync.Mutex{},
	}
}

func (m *Manager) getRoom(roomId *string) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if roomId != nil {
		return m.rooms[*roomId], nil
	}

	return m.getVacantRoom()
}

func (m *Manager) getVacantRoom() (*Room, error) {
	for _, room := range m.rooms {
		if room.hasCapacity() {
			return room, nil
		}
	}
	return m.makeNewRoom()
}

func (m *Manager) makeNewRoom() (*Room, error) {
	room, err := NewRoom()
	if err != nil {
		return nil, err
	}

	m.rooms[room.id] = room
	m.roomIds = append(m.roomIds, room.id)
	return room, nil
}
