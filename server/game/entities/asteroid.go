package entities

import (
	"math"
	"math/rand"
	"server/game/geometry"
	"server/utils"
)

const (
	ASTEROID_MAX_SPEED = 2.0
)

type Asteroid struct {
	Type     EntityType         `json:"type"`
	ID       string             `json:"id"`
	Position EntityPosition     `json:"position"`
	Points   *[]geometry.Vector `json:"points"`

	speed float64

	boundingBox *geometry.BoundingBox
}

func NewAsteroid() (*Asteroid, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	position := randomEntityPosition()
	position.Theta = rand.Float64() * 2 * math.Pi

	points := geometry.NewRandomConvexHull()
	boundingBox := geometry.NewBoundingBox(&points)

	return &Asteroid{
		Type:        AsteroidEntityType,
		ID:          id,
		Position:    position,
		Points:      &points,
		speed:       rand.Float64() * ASTEROID_MAX_SPEED,
		boundingBox: &boundingBox,
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
	a.Position.X += math.Cos(a.Position.Theta) * a.speed
	a.Position.Y += math.Sin(a.Position.Theta) * a.speed
	return true
}

func (a *Asteroid) PollNewEntities() []Entity {
	// TODO: make it split?
	return nil
}
