package room

import (
	"fmt"
	"server/id"
	"sync"
)

// A Manager manages REST API requests and room assignments.
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

// getRoom returns a room. If a roomId is not specified, the first vacant room
// is picked.
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

// getVacantRoom returns the first room that has capacity for more players. If
// all existing rooms are full, a new room is created.
func (m *Manager) getVacantRoom() (*Room, error) {
	for _, room := range m.rooms {
		if room.hasCapacity() {
			return room, nil
		}
	}
	return m.initNewRoom()
}

// initNewRoom initializes a new room.
func (m *Manager) initNewRoom() (*Room, error) {
	id, err := id.NewShortId()
	if err != nil {
		return nil, err
	}

	room := newRoom(id)
	room.init()

	m.rooms[room.id] = room
	m.roomIds = append(m.roomIds, room.id)
	return room, nil
}
