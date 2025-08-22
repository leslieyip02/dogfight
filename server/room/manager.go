package room

import (
	"encoding/json"
	"math/rand"
	"net/http"

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
}

type JoinRequest struct {
	Username string  `json:"username"`
	RoomId   *string `json:"roomId,omitempty"`
}

func NewManager() (*Manager, error) {
	roomManager := Manager{
		rooms:   map[string]*Room{},
		roomIds: []string{},
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

func (m Manager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	var request JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "unable to parse room join request", http.StatusBadRequest)
		return
	}

	client, err := NewClient(request.Username)
	if err != nil {
		http.Error(w, "unable to create client", http.StatusInternalServerError)
		return
	}

	room := m.getRoom(request.RoomId)
	room.Add(client)

	// TODO: use JWT
	body := map[string]string{
		"clientId": client.id,
		"roomId":   room.id,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to create room manager", http.StatusInternalServerError)
	}
}

func (m Manager) HandleConnect(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")
	if clientId == "" {
		http.Error(w, "missing client ID", http.StatusBadRequest)
		return
	}

	roomId := r.URL.Query().Get("roomId")
	if roomId == "" {
		http.Error(w, "missing room ID", http.StatusBadRequest)
		return
	}

	room, found := m.rooms[roomId]
	if !found {
		http.Error(w, "no such room", http.StatusNotFound)
		return
	}

	client, found := room.clients[clientId]
	if !found {
		http.Error(w, "no such client", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "unable to create connection", http.StatusInternalServerError)
		return
	}

	client.conn = conn
	room.Connect(client)
}

func (m Manager) getRoom(roomId *string) *Room {
	if roomId != nil {
		return m.rooms[*roomId]
	}

	randomId := m.roomIds[rand.Intn(len(m.roomIds))]
	return m.rooms[randomId]
}
