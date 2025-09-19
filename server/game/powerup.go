package game

type PowerupType string

const (
	MultishotPowerupType PowerupType = "multishot"
)

type Powerup struct {
	Id       string         `json:"id"`
	Type     PowerupType    `json:"type"`
	Position EntityPosition `json:"position"`
}

const MAX_POWERUP_COUNT = 16
const POWERUP_SPAWN_INTERVAL = FPS * 30
