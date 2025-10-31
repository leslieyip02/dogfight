package entities

import (
	"server/game/geometry"
	"server/pb"
)

const (
	PROJECTILE_RADIUS   = 10.0
	PROJECTILE_SPEED    = 24.0
	PROJECTILE_LIFETIME = 2.4 * FPS
)

var basicProjectileBoundingBoxPoints = geometry.NewRectangleHull(10, 10)
var wideBeamProjectileBoundingBoxPoints = geometry.NewRectangleHull(20, 80)

type Projectile struct {
	entityData *pb.EntityData

	// internal duplicates of EntityData state
	position geometry.Vector
	velocity geometry.Vector
	rotation float64

	boundingBox *geometry.BoundingBox
	onRemove    ProjectileOnRemoveCallback
}

type ProjectileOnRemoveCallback func(other *Entity)

func NewProjectile(
	id string,
	position geometry.Vector,
	velocity geometry.Vector,
	flags AbilityFlag,
	onRemove ProjectileOnRemoveCallback,
) *Projectile {
	rotation := velocity.Angle()
	points := chooseBoundingBoxPoints(AbilityFlag(flags))

	entityData := &pb.EntityData{
		Type:     pb.EntityType_ENTITY_TYPE_PROJECTILE,
		Id:       id,
		Position: position.ToPb(),
		Velocity: velocity.ToPb(),
		Rotation: rotation,
		Data: &pb.EntityData_ProjectileData_{
			ProjectileData: &pb.EntityData_ProjectileData{
				Flags:    uint32(flags),
				Lifetime: PROJECTILE_LIFETIME,
			},
		},
	}

	p := Projectile{
		entityData: entityData,
		position:   position,
		velocity:   velocity,
		rotation:   rotation,
		onRemove:   onRemove,
	}
	p.boundingBox = geometry.NewBoundingBox(&p.position, &p.rotation, points)
	return &p
}

func (p *Projectile) GetEntityType() pb.EntityType {
	return pb.EntityType_ENTITY_TYPE_PROJECTILE
}

func (p *Projectile) GetEntityData() *pb.EntityData {
	return p.entityData
}

func (p *Projectile) GetID() string {
	return p.entityData.Id
}

func (p *Projectile) GetPosition() geometry.Vector {
	return p.position
}

func (p *Projectile) GetVelocity() geometry.Vector {
	return p.velocity
}

func (p *Projectile) GetIsExpired() bool {
	if p.entityData.GetProjectileData().Lifetime < 0 {
		p.onRemove(nil)
		return true
	}
	return false
}

func (p *Projectile) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Projectile) Update() bool {
	p.position.X += p.velocity.X
	p.position.Y += p.velocity.Y
	p.entityData.GetProjectileData().Lifetime--

	p.syncEntityData()
	return true
}

func (p *Projectile) PollNewEntities() []Entity {
	return nil
}

func (p *Projectile) UpdateOnCollision(other Entity) {}

func (p *Projectile) RemoveOnCollision(other Entity) bool {
	p.onRemove(&other)

	switch other.GetEntityType() {
	case pb.EntityType_ENTITY_TYPE_PLAYER:
		return true

	case pb.EntityType_ENTITY_TYPE_POWERUP, pb.EntityType_ENTITY_TYPE_PROJECTILE:
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

func (p *Projectile) syncEntityData() {
	p.entityData.Position.X = p.position.X
	p.entityData.Position.Y = p.position.Y
}
