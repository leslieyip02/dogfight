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

var basicProjectileBoundingBoxPoints = geometry.NewRectangleHull(10, 10)
var wideBeamProjectileBoundingBoxPoints = geometry.NewRectangleHull(20, 80)

type Projectile struct {
	Type     EntityType      `json:"type"`
	ID       string          `json:"id"`
	Position geometry.Vector `json:"position"`
	Velocity geometry.Vector `json:"velocity"`
	Rotation float64         `json:"rotation"`
	Flags    AbilityFlag     `json:"flags"`
	Lifetime int             `json:"lifetime"`

	shooter     *Player
	boundingBox *geometry.BoundingBox
}

func NewProjectile(position geometry.Vector, velocity geometry.Vector, shooter *Player) (*Projectile, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	rotation := velocity.Angle()
	flags := shooter.Flags
	points := chooseBoundingBoxPoints(flags)

	p := Projectile{
		Type:     ProjectileEntityType,
		ID:       id,
		Position: position,
		Velocity: velocity,
		Rotation: rotation,
		Flags:    flags,
		Lifetime: PROJECTILE_LIFETIME,
		shooter:  shooter,
	}
	p.boundingBox = geometry.NewBoundingBox(&p.Position, &p.Rotation, points)
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
	return p.Lifetime < 0
}

func (p *Projectile) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Projectile) Update() bool {
	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
	p.Lifetime--
	return true
}

func (p *Projectile) PollNewEntities() []Entity {
	return nil
}

func (p *Projectile) UpdateOnCollision(other Entity) {}

func (p *Projectile) RemoveOnCollision(other Entity) bool {
	switch other.GetType() {
	case PlayerEntityType:
		p.shooter.Score++
		return true

	case PowerupEntityType, ProjectileEntityType:
		return false

	default:
		return true
	}
}

func chooseBoundingBoxPoints(flags AbilityFlag) *[]*geometry.Vector {
	if isAbilityActive(flags, WideBeamAbilityFlag) {
		return &wideBeamProjectileBoundingBoxPoints
	}
	return &basicProjectileBoundingBoxPoints
}
