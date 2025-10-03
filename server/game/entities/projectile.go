package entities

import (
	"server/game/geometry"
	"server/utils"
)

const (
	PROJECTILE_RADIUS   = 10.0
	PROJECTILE_SPEED    = 24.0
	PROJECTILE_LIFETIME = 2.4 * FPS
)

var projectileBoundingBoxPoints = []*geometry.Vector{
	geometry.NewVector(-10, -10),
	geometry.NewVector(10, -10),
	geometry.NewVector(10, 10),
	geometry.NewVector(-10, 10),
}

type Projectile struct {
	Type     EntityType      `json:"type"`
	ID       string          `json:"id"`
	Position geometry.Vector `json:"position"`
	Velocity geometry.Vector `json:"velocity"`
	Rotation float64         `json:"rotation"`

	lifetime    int
	boundingBox *geometry.BoundingBox
}

func NewProjectile(position geometry.Vector, velocity geometry.Vector) (*Projectile, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	rotation := velocity.Angle()

	p := Projectile{
		Type:     ProjectileEntityType,
		ID:       id,
		Position: position,
		Velocity: velocity,
		Rotation: rotation,
		lifetime: PROJECTILE_LIFETIME,
	}
	p.boundingBox = geometry.NewBoundingBox(
		&p.Position,
		&p.Rotation,
		&projectileBoundingBoxPoints,
	)
	return &p, nil
}

func (p *Projectile) GetType() EntityType {
	return ProjectileEntityType
}

func (p *Projectile) GetID() string {
	return p.ID
}

func (p *Projectile) GetPosition() geometry.Vector {
	return p.Position
}

func (p *Projectile) GetVelocity() geometry.Vector {
	return p.Velocity
}

func (p *Projectile) GetIsExpired() bool {
	return p.lifetime <= 0
}

func (p *Projectile) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Projectile) Update() bool {
	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
	p.lifetime--
	return true
}

func (p *Projectile) PollNewEntities() []Entity {
	return nil
}

func (p *Projectile) RemoveOnCollision(other Entity) bool {
	return other.GetType() != ProjectileEntityType && other.GetType() != PowerupEntityType
}
