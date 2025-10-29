package room

import (
	"sync"
)

type Manager struct {
	session *Session
	rooms   map[string]*Room
	roomIds []string
	mu      sync.Mutex
}

func NewManager(session *Session) *Manager {
	return &Manager{
		session: session,
		rooms:   map[string]*Room{},
		roomIds: []string{},
		mu:      sync.Mutex{},
	}
}

func (m *Manager) getRoom(roomId string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.rooms[roomId]
}

func (m *Manager) getVacantRoom() (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
