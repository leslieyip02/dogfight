package room

import (
	"server/pb"
	"sync"
)

// A Lobby manages rooms.
type Lobby struct {
	rooms   map[string]*Room
	roomIds []string
	mu      sync.Mutex
}

func NewLobby() *Lobby {
	return &Lobby{
		rooms:   map[string]*Room{},
		roomIds: []string{},
		mu:      sync.Mutex{},
	}
}

// GetRoom returns the room with roomId.
func (l *Lobby) GetRoom(roomId string) *Room {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.rooms[roomId]
}

// CreateRoom creates a new room with roomId.
func (l *Lobby) CreateRoom(roomId string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	room := newRoom(roomId)
	room.init()

	l.rooms[roomId] = room
	l.roomIds = append(l.roomIds, roomId)
}

// GetSnapshot gets the game state for the requested room.
func (l *Lobby) GetSnapshot(roomId string) *pb.Event {
	l.mu.Lock()
	defer l.mu.Unlock()

	room, found := l.rooms[roomId]
	if !found {
		return nil
	}
	return room.game.GetSnapshot()
}

// GetSnapshot gets the total player count over all rooms.
func (l *Lobby) GetPlayerCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	count := 0
	for _, room := range l.rooms {
		count += room.getPlayerCount()
	}
	return count
}
