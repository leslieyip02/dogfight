package room

import (
	"fmt"
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

func (m *Manager) getRoom(roomId *string) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if roomId != nil {
		room, found := m.rooms[*roomId]
		if !found {
			return nil, fmt.Errorf("room %v not found", roomId)
		}
		return room, nil
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
