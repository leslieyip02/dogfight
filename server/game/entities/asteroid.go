package entities

import (
	"fmt"
	"math"
	"math/rand"
	"server/game/geometry"
	"server/pb"
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
	entityData *pb.EntityData

	// state
	position    geometry.Vector
	velocity    geometry.Vector
	rotation    float64
	boundingBox *geometry.BoundingBox

	spin   float64
	health int
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

	entityPoints := make([]*pb.Vector, len(points))
	for i, point := range points {
		entityPoints[i] = &pb.Vector{
			X: point.X,
			Y: point.Y,
		}
	}

	entity := &pb.EntityData{
		Type:     pb.EntityType_ENTITY_TYPE_ASTEROID,
		Id:       id,
		Position: &pb.Vector{X: position.X, Y: position.Y},
		Velocity: &pb.Vector{X: velocity.X, Y: velocity.Y},
		Rotation: rotation,
		Data: &pb.EntityData_AsteroidData_{
			AsteroidData: &pb.EntityData_AsteroidData{
				Points: entityPoints,
			},
		},
	}
	a := Asteroid{
		entityData: entity,
		position:   position,
		velocity:   velocity,
		rotation:   rotation,
		spin:       spin,
		health:     ASTEROID_MAX_HEALTH,
	}
	a.boundingBox = geometry.NewBoundingBox(
		&a.position,
		&a.rotation,
		&points,
	)
	return &a, nil
}

func (a *Asteroid) GetEntityType() pb.EntityType {
	return a.entityData.GetType()
}

func (a *Asteroid) GetEntityData() *pb.EntityData {
	return a.entityData
}

func (a *Asteroid) GetID() string {
	return a.entityData.Id
}

func (a *Asteroid) GetPosition() geometry.Vector {
	return a.position
}

func (a *Asteroid) GetVelocity() geometry.Vector {
	return a.velocity
}

func (a *Asteroid) GetIsExpired() bool {
	return false
}

func (a *Asteroid) GetBoundingBox() *geometry.BoundingBox {
	return a.boundingBox
}

func (a *Asteroid) Update() bool {
	a.position.X += a.velocity.X
	a.position.Y += a.velocity.Y
	a.rotation += a.spin

	// copy to entity
	a.entityData.Position.X = a.position.X
	a.entityData.Position.Y = a.position.Y
	a.entityData.Rotation = a.rotation

	return true
}

func (a *Asteroid) PollNewEntities() []Entity {
	return nil
}

func (a *Asteroid) UpdateOnCollision(other Entity) {}

func (a *Asteroid) RemoveOnCollision(other Entity) bool {
	switch other.GetEntityType() {
	case pb.EntityType_ENTITY_TYPE_PROJECTILE:
		a.health--
		return a.health <= 0

	case pb.EntityType_ENTITY_TYPE_POWERUP:
		return false

	default:
		return true
	}
}
