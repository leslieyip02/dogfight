package game

type PowerupType string

const (
	MultishotPowerupType PowerupType = "multishot"
)

type Powerup struct {
	Id       string
	Type     PowerupType
	position EntityPosition
}
