package game

import (
	"encoding/json"
)

type EventType string

const (
	JoinEventType           EventType = "join"
	QuitEventType           EventType = "quit"
	UpdatePositionEventType EventType = "position"
	UpdatePowerupEventType  EventType = "powerup"
	InputEventType          EventType = "input"
)

type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinEventData struct {
	Id       string         `json:"id"`
	Username string         `json:"username"`
	Position EntityPosition `json:"position"`
}

type QuitEventData struct {
	Id string `json:"id"`
}

type UpdatePositionEventData struct {
	Players     map[string]EntityPosition `json:"players"`
	Projectiles map[string]EntityPosition `json:"projectiles"`
}

type UpdatePowerupEventData struct {
	Id       string          `json:"id"`
	Type     PowerupType     `json:"type"`
	Position *EntityPosition `json:"position,omitempty"`
}

type InputEventData struct {
	ClientId     string  `json:"clientId"`
	MouseX       float64 `json:"mouseX"`
	MouseY       float64 `json:"mouseY"`
	MousePressed bool    `json:"mousePressed"`
}

func NewJoinEventMessage(player *Player) ([]byte, error) {
	joinEventData := JoinEventData{
		Id:       player.Id,
		Username: player.Username,
		Position: player.Position,
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
		Id: clientId,
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

func NewUpdatePositionEventMessage(players *map[string]*Player, projectiles *map[string]*Projectile) ([]byte, error) {
	playerPositions := map[string]EntityPosition{}
	for id, player := range *players {
		playerPositions[id] = player.Position
	}
	projectilePositions := map[string]EntityPosition{}
	for id, projectile := range *projectiles {
		projectilePositions[id] = projectile.Position
	}

	data, err := json.Marshal(UpdatePositionEventData{
		Players:     playerPositions,
		Projectiles: projectilePositions,
	})
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: UpdatePositionEventType,
		Data: data,
	}
	return json.Marshal(message)
}

func NewUpdatePowerupEventMessage(powerup *Powerup, active bool) ([]byte, error) {
	var position *EntityPosition
	if active {
		position = &powerup.Position
	} else {
		position = nil
	}

	data, err := json.Marshal(UpdatePowerupEventData{
		Id:       powerup.Id,
		Type:     powerup.Type,
		Position: position,
	})
	if err != nil {
		return nil, err
	}

	message := Event{
		Type: UpdatePowerupEventType,
		Data: data,
	}
	return json.Marshal(message)
}
