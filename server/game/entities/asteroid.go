package entities

import (
	"fmt"
	"math"
	"math/rand"
	"server/game/geometry"
	"server/utils"
)

const (
	ASTEROID_MAX_SPEED  = 0.5
	ASTEROID_MAX_SPIN   = 0.001
	ASTEROID_MAX_HEALTH = 3

	ASTEROID_MIN_NUM_POINTS = 8
	ASTEROID_MAX_NUM_POINTS = 16
	ASTEROID_MIN_RADIUS     = 20
	ASTEROID_MAX_RADIUS     = 100
	ASTEROID_MIN_AREA       = 200
)

type Asteroid struct {
	Type     EntityType          `json:"type"`
	ID       string              `json:"id"`
	Position geometry.Vector     `json:"position"`
	Velocity geometry.Vector     `json:"velocity"`
	Rotation float64             `json:"rotation"`
	Points   *[]*geometry.Vector `json:"points"`

	spin        float64
	health      int
	boundingBox *geometry.BoundingBox
}

func newRandomAsteroid() (*Asteroid, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	position := *geometry.NewRandomVector(0, 0, SPAWN_AREA_WIDTH, SPAWN_AREA_HEIGHT)
	velocity := *geometry.NewRandomVector(0, 0, ASTEROID_MAX_SPEED, ASTEROID_MAX_SPEED)
	rotation := rand.Float64() * math.Pi * 2
	spin := rand.Float64()*ASTEROID_MAX_SPIN*2 - ASTEROID_MAX_SPIN

	points := geometry.NewRandomConvexHull(
		ASTEROID_MIN_NUM_POINTS,
		ASTEROID_MAX_NUM_POINTS,
		ASTEROID_MIN_RADIUS,
		ASTEROID_MAX_RADIUS,
	)
	if geometry.HullArea(points) < ASTEROID_MIN_AREA {
		return nil, fmt.Errorf("too small")
	}

	a := Asteroid{
		Type:     AsteroidEntityType,
		ID:       id,
		Position: position,
		Velocity: velocity,
		Rotation: rotation,
		Points:   &points,
		spin:     spin,
		health:   ASTEROID_MAX_HEALTH,
	}
	a.boundingBox = geometry.NewBoundingBox(
		&a.Position,
		&a.Rotation,
		&points,
	)
	return &a, nil
}

func (a *Asteroid) GetType() EntityType {
	return AsteroidEntityType
}

func (a *Asteroid) GetID() string {
	return a.ID
}

func (a *Asteroid) GetPosition() geometry.Vector {
	return a.Position
}

func (a *Asteroid) GetVelocity() geometry.Vector {
	return a.Velocity
}

func (a *Asteroid) GetIsExpired() bool {
	return false
}

func (a *Asteroid) GetBoundingBox() *geometry.BoundingBox {
	return a.boundingBox
}

func (a *Asteroid) Update() bool {
	a.Position.X += a.Velocity.X
	a.Position.Y += a.Velocity.Y
	a.Rotation += a.spin
	return true
}

func (a *Asteroid) PollNewEntities() []Entity {
	return nil
}

func (a *Asteroid) UpdateOnCollision(other Entity) {}

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
