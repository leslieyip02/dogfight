package room

import (
	"server/game"
	"server/utils"

	"sync"
)

type Room struct {
	id        string
	players   map[string]*game.Player
	broadcast []chan byte
	mu        sync.Mutex
}

func NewRoom() (*Room, error) {
	id, err := utils.GetShortId()
	if err != nil {
		return nil, err
	}

	room := Room{
		id:        id,
		players:   map[string]*game.Player{},
		broadcast: []chan byte{},
		mu:        sync.Mutex{},
	}
	return &room, nil
}

func (r *Room) AddPlayer(player *game.Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: add capacity check
	r.players[player.Id] = player
}

func (r *Room) RemovePlayer(player *game.Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.players, player.Id)
}
