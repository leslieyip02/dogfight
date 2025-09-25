package room

import (
	"encoding/json"
	"net/http"
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

	client, err := NewClient(request.Username)
	if err != nil {
		http.Error(w, "unable to create client", http.StatusInternalServerError)
		return
	}

	room, err := m.getRoom(request.RoomId)
	if err != nil {
		http.Error(w, "unable to get room", http.StatusInternalServerError)
		return
	}
	room.Add(client)

	token, err := m.session.createToken(room.id, client.id)
	if err != nil {
		http.Error(w, "unable to issue token", http.StatusInternalServerError)
		return
	}

	body := map[string]string{
		"roomId":   room.id,
		"clientId": client.id,
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
	room.connect(client)
}

func (m *Manager) HandleFetchState(w http.ResponseWriter, r *http.Request) {
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
	room, err := m.getRoom(&roomId)
	if err != nil {
		http.Error(w, "unable to get room", http.StatusInternalServerError)
		return
	}

	body := room.game.GetState()
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to get room state", http.StatusInternalServerError)
	}
}
