package room

import (
	"encoding/json"
	"net/http"
)

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

func (m *Manager) HandleConnect(w http.ResponseWriter, r *http.Request) {
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
	room.connect(client)
}

func (m *Manager) HandleFetchState(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	if roomId == "" {
		http.Error(w, "missing room ID", http.StatusBadRequest)
		return
	}

	room := m.getRoom(&roomId)
	body := room.game.GetState()
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "unable to get room state", http.StatusInternalServerError)
	}
}
