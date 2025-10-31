package entities

import (
	"server/game/geometry"
	"server/pb"
	"server/utils"
)

const (
	MAX_POWERUP_COUNT = 16
)

var powerupBoundingBoxPoints = geometry.NewRectangleHull(20, 20)

type Powerup struct {
	entityData *pb.EntityData

	// state
	position    geometry.Vector
	velocity    geometry.Vector
	rotation    float64
	boundingBox *geometry.BoundingBox
}

func newRandomPowerup() (*Powerup, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}
	position := *geometry.NewRandomVector(0, 0, SPAWN_AREA_WIDTH, SPAWN_AREA_HEIGHT)
	velocity := *geometry.NewVector(0, 0)
	rotation := 0.0
	ability := newRandomAbility()

	entity := &pb.EntityData{
		Type:     pb.EntityType_ENTITY_TYPE_POWERUP,
		Id:       id,
		Position: &pb.Vector{X: position.X, Y: position.Y},
		Velocity: &pb.Vector{X: velocity.X, Y: velocity.Y},
		Rotation: rotation,
		Data: &pb.EntityData_PowerupData_{
			PowerupData: &pb.EntityData_PowerupData{
				Ability: uint32(ability),
			},
		},
	}

	p := Powerup{
		entityData: entity,
		position:   position,
		velocity:   velocity,
		rotation:   rotation,
	}
	p.boundingBox = geometry.NewBoundingBox(
		&p.position,
		&p.rotation,
		&powerupBoundingBoxPoints,
	)
	return &p, nil
}

func (p *Powerup) GetEntityType() pb.EntityType {
	return pb.EntityType_ENTITY_TYPE_POWERUP
}

func (p *Powerup) GetEntityData() *pb.EntityData {
	return p.entityData
}

func (p *Powerup) GetID() string {
	return p.entityData.Id
}

func (p *Powerup) GetPosition() geometry.Vector {
	return p.position
}

func (p *Powerup) GetVelocity() geometry.Vector {
	return p.velocity
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

func (p *Powerup) UpdateOnCollision(other Entity) {}

func (p *Powerup) RemoveOnCollision(other Entity) bool {
	return other.GetEntityType() == pb.EntityType_ENTITY_TYPE_PLAYER
}
