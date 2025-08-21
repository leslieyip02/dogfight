package game

import (
	"encoding/json"
	"math/rand"
)

type EventType string

const (
	JoinEventType   EventType = "join"
	QuitEventType   EventType = "quit"
	UpdateEventType EventType = "update"
)

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinEventData struct {
	ClientId string  `json:"clientId"`
	Username string  `json:"username"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
}

type QuitEventData struct {
	ClientId string `json:"clientId"`
}

func NewJoinEventMessage(clientId string, username string) ([]byte, error) {
	joinEvent := JoinEventData{
		ClientId: clientId,
		Username: username,
		X:        rand.Float32() * 100,
		Y:        rand.Float32() * 100,
	}
	data, err := json.Marshal(joinEvent)
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: JoinEventType,
		Data: data,
	}
	return json.Marshal(message)
}

func NewQuitEventMessage(clientId string) ([]byte, error) {
	quitEvent := QuitEventData{
		ClientId: clientId,
	}
	data, err := json.Marshal(quitEvent)
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: QuitEventType,
		Data: data,
	}
	return json.Marshal(message)
}
