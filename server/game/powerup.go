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
	ID       string         `json:"id"`
	Type     PowerupType    `json:"type"`
	Position EntityPosition `json:"position"`
}

func (p *Powerup) GetType() EntityType {
	return PowerupEntityType
}

func (p *Powerup) GetID() string {
	return p.ID
}

func (p *Powerup) GetPosition() EntityPosition {
	return p.Position
}

func (p *Powerup) GetIsExpired() bool {
	return false
}

func (p *Powerup) Update(g *Game) {
	// No update needed
}
