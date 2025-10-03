package entities

import (
	"fmt"
	"math"
	"math/rand"
	"server/game/geometry"
	"server/utils"
)

const (
	ASTEROID_MAX_SPEED      = 0.5
	ASTEROID_MAX_ROTATION   = 0.001
	ASTEROID_MAX_HEALTH     = 3
	ASTEROID_MIN_RADIUS     = 20
	ASTEROID_MAX_RADIUS     = 100
	ASTEROID_MIN_NUM_POINTS = 8
	ASTEROID_MAX_NUM_POINTS = 16
	ASTEROID_MIN_AREA       = 200
)

type Asteroid struct {
	Type     EntityType          `json:"type"`
	ID       string              `json:"id"`
	Position EntityPosition      `json:"position"`
	Points   *[]*geometry.Vector `json:"points"`

	velocity *geometry.Vector
	rotation float64
	health   int

	boundingBox *geometry.BoundingBox
}

func NewAsteroid() (*Asteroid, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	position := randomEntityPosition()
	position.Theta = rand.Float64() * 2 * math.Pi

	points := geometry.NewRandomConvexHull(ASTEROID_MIN_NUM_POINTS, ASTEROID_MAX_NUM_POINTS, ASTEROID_MIN_RADIUS, ASTEROID_MAX_RADIUS)
	if geometry.HullArea(points) < ASTEROID_MIN_AREA {
		return nil, fmt.Errorf("too small")
	}
	boundingBox := geometry.NewBoundingBox(&points)

	return &Asteroid{
		Type:        AsteroidEntityType,
		ID:          id,
		Position:    position,
		Points:      &points,
		velocity:    geometry.NewVector(rand.Float64()*ASTEROID_MAX_SPEED, rand.Float64()*ASTEROID_MAX_SPEED),
		rotation:    rand.Float64()*ASTEROID_MAX_ROTATION*2 - ASTEROID_MAX_ROTATION,
		health:      ASTEROID_MAX_HEALTH,
		boundingBox: boundingBox,
	}, nil
}

func (a *Asteroid) GetType() EntityType {
	return AsteroidEntityType
}

func (a *Asteroid) GetID() string {
	return a.ID
}

func (a *Asteroid) GetPosition() EntityPosition {
	return a.Position
}

func (a *Asteroid) GetIsExpired() bool {
	return false
}

func (a *Asteroid) GetBoundingBox() *geometry.BoundingBox {
	return a.boundingBox.Transform(a.Position.X, a.Position.Y, a.Position.Theta)
}

func (a *Asteroid) Update() bool {
	a.Position.X += a.velocity.X
	a.Position.Y += a.velocity.Y
	a.Position.Theta += a.rotation
	return true
}

func (a *Asteroid) PollNewEntities() []Entity {
	return nil
}

func (a *Asteroid) RemoveOnCollision(other Entity) bool {
	switch other.GetType() {
	case ProjectileEntityType:
		a.health--
		return a.health <= 0

	case PowerupEntityType:
		return false

	default:
		return true
	}
}
