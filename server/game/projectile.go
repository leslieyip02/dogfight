package game

import (
	"math"
	"server/utils"
)

const (
	PROJECTILE_RADIUS   = 10.0
	PROJECTILE_SPEED    = 24.0
	PROJECTILE_LIFETIME = 2.4 * FPS
)

type Projectile struct {
	Type     EntityType     `json:"type"`
	ID       string         `json:"id"`
	Position EntityPosition `json:"position"`
	speed    float64
	lifetime int
}

func NewProjectile(position EntityPosition) (*Projectile, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	return &Projectile{
		Type:     ProjectileEntityType,
		ID:       id,
		Position: position,
		speed:    PROJECTILE_SPEED,
		lifetime: PROJECTILE_LIFETIME,
	}, err
}

func (p *Projectile) GetType() EntityType {
	return ProjectileEntityType
}

func (p *Projectile) GetID() string {
	return p.ID
}

func (p *Projectile) GetPosition() EntityPosition {
	return p.Position
}

func (p *Projectile) GetIsExpired() bool {
	return p.lifetime <= 0
}

func (p *Projectile) Update(g *Game) bool {
	p.Position.X += math.Cos(p.Position.Theta) * p.speed
	p.Position.Y += math.Sin(p.Position.Theta) * p.speed
	p.lifetime--
	return true
}
