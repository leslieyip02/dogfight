package room

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"server/game"
)

type RoomManager struct {
	rooms   map[string]*Room
	roomIds []string
}

type JoinRoomRequest struct {
	Username string  `json:"username"`
	RoomId   *string `json:"roomId"`
}

func NewRoomManager() (*RoomManager, error) {
	roomManager := RoomManager{
		rooms: map[string]*Room{},
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

func (m RoomManager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	var request JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "unable to parse join room request", http.StatusBadRequest)
		return
	}

	player, err := game.NewPlayer(request.Username)
	if err != nil {
		http.Error(w, "unable to create player", http.StatusInternalServerError)
		return
	}

	room := m.pickRoom(request.RoomId)
	room.AddPlayer(player)

	body := map[string]string{
		"playerId": player.Id,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to create room manager", http.StatusInternalServerError)
	}
}

func (m RoomManager) pickRoom(roomId *string) *Room {
	if roomId != nil {
		return m.rooms[*roomId]
	}

	randomId := m.roomIds[rand.Intn(len(m.roomIds))]
	return m.rooms[randomId]
}
