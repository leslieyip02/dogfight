package room

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/id"

	"google.golang.org/protobuf/proto"
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

	room, err := m.getRoom(request.RoomId)
	if err != nil {
		http.Error(w, "unable to assign room", http.StatusInternalServerError)
		return
	}

	clientId, err := id.NewShortId()
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

	room, err := m.getRoom(&claims.roomId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	message, err := proto.Marshal(room.game.GetSnapshot())
	if err != nil {
		http.Error(w, "unable to serialize room state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err = w.Write(message); err != nil {
		http.Error(w, "unable to get room state", http.StatusInternalServerError)
		return
	}
}
