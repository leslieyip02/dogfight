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

// HandleJoin assigns a client to a room and issues a session token for
// subsequent requests.
func (m *Manager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	var request JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		message := "unable to parse room join request"
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	room, err := m.getRoom(request.RoomId)
	if err != nil {
		message := "unable to assign room"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	clientId, err := id.NewShortId()
	if err != nil {
		message := "unable to assign client"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	token, err := m.session.createToken(room.id, clientId, request.Username)
	if err != nil {
		message := "unable to issue token"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	body := map[string]string{
		"clientId": clientId,
		"token":    token,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		message := "unable to write body"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
}

// HandleWS creates a WebSocket connection with the client.
func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := m.session.parseToken(token)
	if err != nil {
		message := err.Error()
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	room, found := m.rooms[claims.roomId]
	if !found {
		message := "no such room"
		http.Error(w, message, http.StatusNotFound)
		return
	}

	conn, err := m.session.createConn(&w, r, claims.clientId, room)
	if err != nil {
		message := fmt.Sprintf("unable to create connection %v", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	client := newClient(claims.clientId, claims.username, conn)

	room.add(client)
	room.connect(client)
}

// HandleSnapshot sends with a snapshot of the game's entire state to the
// client.
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

	snapshot, err := proto.Marshal(room.game.GetSnapshot())
	if err != nil {
		message := "unable to serialize room state"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err = w.Write(snapshot); err != nil {
		message := "unable to get room state"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
}
