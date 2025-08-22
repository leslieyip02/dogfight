package game

import (
	"encoding/json"
)

type EventType string

const (
	JoinEventType   EventType = "join"
	QuitEventType   EventType = "quit"
	UpdateEventType EventType = "update"
	InputEventType  EventType = "input"
)

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinEventData struct {
	ClientId string         `json:"clientId"`
	Username string         `json:"username"`
	Position EntityPosition `json:"position"`
}

type QuitEventData struct {
	ClientId string `json:"clientId"`
}

type UpdateEventData map[string]EntityPosition

type InputEventData struct {
	ClientId     string  `json:"clientId"`
	MouseX       float64 `json:"mouseX"`
	MouseY       float64 `json:"mouseY"`
	MousePressed bool    `json:"mousePressed"`
}

func NewJoinEventMessage(player *Player) ([]byte, error) {
	joinEventData := JoinEventData{
		ClientId: player.Id,
		Username: player.Username,
		Position: player.position,
	}
	data, err := json.Marshal(joinEventData)
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
	quitEventData := QuitEventData{
		ClientId: clientId,
	}
	data, err := json.Marshal(quitEventData)
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: QuitEventType,
		Data: data,
	}
	return json.Marshal(message)
}

func NewUpdateEventMessage(game *Game) ([]byte, error) {
	updateEventData := make(UpdateEventData)
	for clientId, player := range game.players {
		updateEventData[clientId] = player.position
	}

	data, err := json.Marshal(updateEventData)
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: UpdateEventType,
		Data: data,
	}
	return json.Marshal(message)
}
