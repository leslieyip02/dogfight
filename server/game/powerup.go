package game

import (
	"server/game/geometry"
	"server/utils"
)

type PowerupAbility string

const (
	MAX_POWERUP_COUNT      = 16
	POWERUP_SPAWN_INTERVAL = 30 * FPS
)

// TODO: add more powerups (e.g. invincibilty)
const (
	MultishotPowerupType PowerupAbility = "multishot"
)

var powerupBoundingBox = geometry.NewBoundingBox(&[]geometry.Vector{
	geometry.NewVector(-10, -10),
	geometry.NewVector(10, -10),
	geometry.NewVector(10, 10),
	geometry.NewVector(-10, 10),
})

type Powerup struct {
	Type     EntityType     `json:"type"`
	ID       string         `json:"id"`
	Position EntityPosition `json:"position"`
	Ability  PowerupAbility `json:"ability"`

	boundingBox *geometry.BoundingBox
}

func NewPowerup(ability PowerupAbility) (*Powerup, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	return &Powerup{
		Type:        PowerupEntityType,
		ID:          id,
		Position:    randomEntityPosition(),
		Ability:     ability,
		boundingBox: &powerupBoundingBox,
	}, nil
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

func (p *Powerup) GetBoundingBox() *geometry.BoundingBox {
	return powerupBoundingBox.Transform(p.Position.X, p.Position.Y, p.Position.Theta)
}

func (p *Powerup) Update(g *Game) bool {
	return false
}
