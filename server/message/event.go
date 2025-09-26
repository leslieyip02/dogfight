package message

import (
	"encoding/json"
	"server/game"
	"time"
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
	ID       string              `json:"id"`
	Username string              `json:"username"`
	Position game.EntityPosition `json:"position"`
}

type QuitEventData struct {
	ID string `json:"id"`
}

type UpdatePositionEventData struct {
	Timestamp   int64                          `json:"timestamp"`
	Players     map[string]game.EntityPosition `json:"players"`
	Projectiles map[string]game.EntityPosition `json:"projectiles"`
}

type UpdatePowerupEventData struct {
	ID       string               `json:"id"`
	Type     game.PowerupType     `json:"type"`
	Position *game.EntityPosition `json:"position,omitempty"`
}

type InputEventData struct {
	ClientID     string  `json:"clientId"`
	MouseX       float64 `json:"mouseX"`
	MouseY       float64 `json:"mouseY"`
	MousePressed bool    `json:"mousePressed"`
}

// Event Factory Functions
func CreateJoinEvent(player *game.Player) ([]byte, error) {
	return NewEvent(JoinEventType).
		WithData(JoinEventData{
			ID:       player.Id,
			Username: player.Username,
			Position: player.Position,
		}).
		Build()
}

func CreateQuitEvent(clientID string) ([]byte, error) {
	return NewEvent(QuitEventType).
		WithData(QuitEventData{
			ID: clientID,
		}).
		Build()
}

func CreateUpdatePositionEvent(players map[string]*game.Player, projectiles map[string]*game.Projectile) ([]byte, error) {
	playerPositions := make(map[string]game.EntityPosition, len(players))
	for id, player := range players {
		playerPositions[id] = player.Position
	}

	projectilePositions := make(map[string]game.EntityPosition, len(projectiles))
	for id, projectile := range projectiles {
		projectilePositions[id] = projectile.Position
	}

	return NewEvent(UpdatePositionEventType).
		WithData(UpdatePositionEventData{
			Timestamp:   time.Now().UnixNano(),
			Players:     playerPositions,
			Projectiles: projectilePositions,
		}).
		Build()
}

func CreateUpdatePowerupEvent(powerup *game.Powerup, active bool) ([]byte, error) {
	data := UpdatePowerupEventData{
		ID:   powerup.Id,
		Type: powerup.Type,
	}

	if active {
		data.Position = &powerup.Position
	}

	return NewEvent(UpdatePowerupEventType).
		WithData(data).
		Build()
}
