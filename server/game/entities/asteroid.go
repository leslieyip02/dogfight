package entities

import (
	"math"
	"math/rand"
	"server/game/geometry"
	"server/utils"
)

const (
	ASTEROID_MAX_SPEED    = 1.0
	ASTEROID_MAX_ROTATION = 0.001
	ASTEROID_MAX_HEALTH   = 3
)

type Asteroid struct {
	Type     EntityType         `json:"type"`
	ID       string             `json:"id"`
	Position EntityPosition     `json:"position"`
	Points   *[]geometry.Vector `json:"points"`

	speed     float64
	rotation  float64
	health    int
	destroyed bool

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
		rotation:    rand.Float64()*ASTEROID_MAX_ROTATION*2 - ASTEROID_MAX_ROTATION,
		health:      ASTEROID_MAX_HEALTH,
		destroyed:   false,
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
	a.Position.Theta += a.rotation
	return true
}

func (a *Asteroid) PollNewEntities() []Entity {
	if !a.destroyed {
		return nil
	}
	return a.split()
}

func (a *Asteroid) split() []Entity {
	// TODO: finish this/reconsider
	if len(*a.Points) < 4 {
		return nil
	}

	fragments := []Entity{}
	// mid := 2 + rand.Intn(len(*a.Points)-3)
	// starts := []int{0, mid}
	// for i, start := range starts {
	// 	next := starts[(i+1)%len(starts)]

	// 	points := []geometry.Vector{}
	// 	for j := start; j != next; j = (j + 1) % len(*a.Points) {
	// 		points = append(points, geometry.NewVector((*a.Points)[j].X*0.8, (*a.Points)[j].Y*0.8))
	// 	}
	// 	points = append(points, geometry.NewVector((*a.Points)[next].X*0.8, (*a.Points)[next].Y*0.8))

	// 	if len(points) < 3 {
	// 		continue
	// 	}

	// 	id, err := utils.NewShortId()
	// 	if err != nil {
	// 		continue
	// 	}

	// 	theta := points[len(points)-1].Sub(&points[0]).Normal().Multiply(-1).Angle()
	// 	position := EntityPosition{
	// 		X:     a.Position.X + math.Cos(theta)*20,
	// 		Y:     a.Position.Y + math.Sin(theta)*20,
	// 		Theta: theta,
	// 	}
	// 	boundingBox := geometry.NewBoundingBox(&points)

	// 	fragment := &Asteroid{
	// 		Type:        AsteroidEntityType,
	// 		ID:          id,
	// 		Position:    position,
	// 		Points:      &points,
	// 		speed:       rand.Float64() * ASTEROID_MAX_SPEED,
	// 		rotation:    rand.Float64()*ASTEROID_MAX_ROTATION*2 - ASTEROID_MAX_ROTATION,
	// 		destroyed:   false,
	// 		boundingBox: &boundingBox,
	// 	}
	// 	fragments = append(fragments, fragment)
	// }
	return fragments
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
