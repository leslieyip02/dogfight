package game

import (
	"encoding/json"
)

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type EventType string

const (
	JoinEventType     EventType = "join"
	QuitEventType     EventType = "quit"
	RespawnEventType  EventType = "respawn"
	InputEventType    EventType = "input"
	SnapshotEventType EventType = "snapshot"
	DeltaEventType    EventType = "delta"
)

type JoinEventData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type QuitEventData struct {
	ID string `json:"id"`
}

type RespawnEventData struct {
	ID string `json:"id"`
}

type InputEventData struct {
	ID           string  `json:"id"`
	MouseX       float64 `json:"mouseX"`
	MouseY       float64 `json:"mouseY"`
	MousePressed bool    `json:"mousePressed"`
}

type SnapshotEventData struct {
	Timestamp int64             `json:"timestamp"`
	Entities  map[string]Entity `json:"entities"`
}

type DeltaEventData struct {
	Timestamp int64             `json:"timestamp"`
	Updated   map[string]Entity `json:"updated"`
	Removed   []string          `json:"removed"`
}

func CreateMessage(eventType EventType, data any) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	event := Event{
		Type: eventType,
		Data: bytes,
	}
	return json.Marshal(event)
}
