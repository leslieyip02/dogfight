package game

import (
	"encoding/json"
)

type Game struct {
	Send      chan []byte
	Broadcast chan []byte
}

func NewGame() Game {
	return Game{
		Send:      make(chan []byte),
		Broadcast: make(chan []byte),
	}
}

func (g *Game) Run() {
	for {
		select {
		case data := <-g.Send:
			var event Event
			json.Unmarshal(data, &event)

			switch event.Type {
			case EventTypeJoin:
				g.Broadcast <- data

			case EventTypeQuit:
				g.Broadcast <- data
			}
		}
	}
}
