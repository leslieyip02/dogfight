package entities

import (
	"server/game/geometry"
	"server/pb"
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

	// internal duplicates of EntityData state
	position geometry.Vector
	velocity geometry.Vector
	rotation float64

	boundingBox *geometry.BoundingBox
	spin        float64
	health      int
}

func NewAsteroid(
	id string,
	position geometry.Vector,
	velocity geometry.Vector,
	rotation float64,
	points *[]*geometry.Vector,
	spin float64,
) *Asteroid {
	entityPoints := make([]*pb.Vector, len(*points))
	for i, point := range *points {
		entityPoints[i] = point.ToPb()
	}
	entityData := &pb.EntityData{
		Type:     pb.EntityType_ENTITY_TYPE_ASTEROID,
		Id:       id,
		Position: position.ToPb(),
		Velocity: velocity.ToPb(),
		Rotation: rotation,
		Data: &pb.EntityData_AsteroidData_{
			AsteroidData: &pb.EntityData_AsteroidData{
				Points: entityPoints,
			},
		},
	}

	a := &Asteroid{
		entityData: entityData,
		position:   position,
		velocity:   velocity,
		rotation:   rotation,
		spin:       spin,
		health:     ASTEROID_MAX_HEALTH,
	}
	a.boundingBox = geometry.NewBoundingBox(&a.position, &a.rotation, points)
	return a
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

	a.syncEntityData()
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

func (a *Asteroid) syncEntityData() {
	a.entityData.Position.X = a.position.X
	a.entityData.Position.Y = a.position.Y
	a.entityData.Rotation = a.rotation
}
