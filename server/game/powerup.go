package game

type PowerupType string

const (
	MAX_POWERUP_COUNT      = 16
	POWERUP_SPAWN_INTERVAL = 30 * FPS
)

// TODO: add more powerups (e.g. invincibilty)
const (
	MultishotPowerupType PowerupType = "multishot"
)

type Powerup struct {
	Id       string         `json:"id"`
	Type     PowerupType    `json:"type"`
	Position EntityPosition `json:"position"`
}
