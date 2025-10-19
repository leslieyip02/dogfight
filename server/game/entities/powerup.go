package entities

import (
	"server/game/geometry"
	"server/utils"
)

const (
	MAX_POWERUP_COUNT      = 16
	POWERUP_SPAWN_INTERVAL = 30 * FPS
)

var powerupBoundingBoxPoints = []*geometry.Vector{
	geometry.NewVector(-10, -10),
	geometry.NewVector(10, -10),
	geometry.NewVector(10, 10),
	geometry.NewVector(-10, 10),
}

type Powerup struct {
	Type     EntityType      `json:"type"`
	ID       string          `json:"id"`
	Position geometry.Vector `json:"position"`
	Velocity geometry.Vector `json:"velocity"`
	Rotation float64         `json:"rotation"`
	Ability  AbilityFlag     `json:"ability"`

	boundingBox *geometry.BoundingBox
}

func NewPowerup(ability AbilityFlag) (*Powerup, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	position := *geometry.NewRandomVector(0, 0, SPAWN_AREA_WIDTH, SPAWN_AREA_HEIGHT)
	velocity := *geometry.NewVector(0, 0)
	rotation := 0.0

	p := Powerup{
		Type:     PowerupEntityType,
		ID:       id,
		Position: position,
		Velocity: velocity,
		Rotation: rotation,
		Ability:  ability,
	}
	p.boundingBox = geometry.NewBoundingBox(
		&p.Position,
		&p.Rotation,
		&powerupBoundingBoxPoints,
	)
	return &p, nil
}

func (p *Powerup) GetType() EntityType {
	return PowerupEntityType
}

func (p *Powerup) GetID() string {
	return p.ID
}

func (p *Powerup) GetPosition() geometry.Vector {
	return p.Position
}

func (p *Powerup) GetVelocity() geometry.Vector {
	return p.Velocity
}

func (p *Powerup) GetIsExpired() bool {
	return false
}

func (p *Powerup) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Powerup) Update() bool {
	return false
}

func (p *Powerup) PollNewEntities() []Entity {
	return nil
}

func (p *Powerup) RemoveOnCollision(other Entity) bool {
	return other.GetType() == PlayerEntityType
}
