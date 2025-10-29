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

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := m.session.parseToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	room, found := m.rooms[claims.roomId]
	if !found {
		http.Error(w, "no such room", http.StatusNotFound)
		return
	}

	conn, err := m.session.createConn(&w, r, claims.clientId, room)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to create connection %v", err), http.StatusInternalServerError)
		return
	}
	client := NewClient(claims.clientId, claims.username, conn)

	room.Add(client)
	room.connect(client)
}

func (m *Manager) HandleSnapshot(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := m.session.parseToken(token)
	if err != nil {
		http.Error(w, "unable to validate JWT", http.StatusUnauthorized)
		return
	}

	room := m.getRoom(claims.roomId)
	body := room.game.GetSnapshot()
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to get room state", http.StatusInternalServerError)
	}
}
