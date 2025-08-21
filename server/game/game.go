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
	for message := range g.Send {
		var event Event
		json.Unmarshal(message, &event)

		switch event.Type {
		case JoinEventType:
			g.Broadcast <- message

		case QuitEventType:
			g.Broadcast <- message

		case InputEventType:
			var data InputEventData
			json.Unmarshal(event.Data, &data)

			// TODO: handle input
		}
	}
}
