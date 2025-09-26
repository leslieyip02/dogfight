package room

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/utils"
)

type JoinRequest struct {
	Username string  `json:"username"`
	RoomId   *string `json:"roomId,omitempty"`
}

func (m *Manager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	var request JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "unable to parse room join request", http.StatusBadRequest)
		return
	}

	var room *Room
	if request.RoomId != nil {
		room = m.getRoom(*request.RoomId)
	} else {
		// TODO: this is kinda dumb
		room2, err := m.getVacantRoom()
		if err != nil {
			http.Error(w, "unable to assign room", http.StatusInternalServerError)
			return
		}
		room = room2
	}

	clientId, err := utils.NewShortId()
	if err != nil {
		http.Error(w, "unable to assign client", http.StatusInternalServerError)
		return
	}

	token, err := m.session.createToken(room.id, clientId, request.Username)
	if err != nil {
		http.Error(w, "unable to issue token", http.StatusInternalServerError)
		return
	}

	body := map[string]string{
		"clientId": clientId,
		"token":    token,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to write body", http.StatusInternalServerError)
		return
	}
}

func (m *Manager) HandleConnect(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := m.session.validateToken(token)
	if err != nil {
		http.Error(w, "unable to validate JWT", http.StatusUnauthorized)
		return
	}

	roomId, found := claims["roomId"].(string)
	if !found {
		http.Error(w, "missing room ID", http.StatusBadRequest)
		return
	}

	room, found := m.rooms[roomId]
	if !found {
		http.Error(w, "no such room", http.StatusNotFound)
		return
	}

	clientId, found := claims["clientId"].(string)
	if !found {
		http.Error(w, "missing client ID", http.StatusBadRequest)
		return
	}

	username, found := claims["username"].(string)
	if !found {
		username = "testificate"
	}

	client := NewClient(clientId, username)
	room.Add(client)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to create connection %v", err), http.StatusInternalServerError)
		return
	}

	client.conn = conn
	room.connect(client)
}

func (m *Manager) HandleFetchSnapshot(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := m.session.validateToken(token)
	if err != nil {
		http.Error(w, "unable to validate JWT", http.StatusUnauthorized)
		return
	}

	roomId, found := claims["roomId"].(string)
	if !found {
		http.Error(w, "missing room ID", http.StatusBadRequest)
		return
	}
	room := m.getRoom(roomId)
	body := room.game.GetSnapshot()
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to get room state", http.StatusInternalServerError)
	}
}
