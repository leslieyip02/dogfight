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
