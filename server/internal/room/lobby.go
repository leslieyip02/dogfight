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

func (l *Lobby) GetStatus() *pb.StatusResponse {
	l.mu.Lock()
	defer l.mu.Unlock()

	roomStatuses := make([]*pb.StatusResponse_RoomStatus, len(l.rooms))
	i := 0
	for roomId, room := range l.rooms {
		roomStatuses[i] = &pb.StatusResponse_RoomStatus{
			RoomId:    roomId,
			Occupancy: room.getOccupancy(),
		}
		i++
	}

	return &pb.StatusResponse{
		RoomStatuses: roomStatuses,
	}
}
