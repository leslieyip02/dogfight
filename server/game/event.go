package game

import (
	"encoding/json"
)

type EventType string

const (
	EventTypeJoin EventType = "join"
	EventTypeQuit EventType = "quit"
)

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinEvent struct {
	ClientId string `json:"clientId"`
}

type QuitEvent struct {
	ClientId string `json:"clientId"`
}

func NewJoinEventMessage(clientId string) ([]byte, error) {
	joinEvent := JoinEvent{
		ClientId: clientId,
	}
	data, err := json.Marshal(joinEvent)
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: EventTypeJoin,
		Data: data,
	}
	return json.Marshal(message)
}

func NewQuitEventMessage(clientId string) ([]byte, error) {
	quitEvent := QuitEvent{
		ClientId: clientId,
	}
	data, err := json.Marshal(quitEvent)
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: EventTypeQuit,
		Data: data,
	}
	return json.Marshal(message)
}
