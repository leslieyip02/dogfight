package entities

import (
	"server/game/geometry"
	"server/pb"
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
	entity *pb.Entity

	// state
	shooter     *Player
	position    geometry.Vector
	velocity    geometry.Vector
	rotation    float64
	boundingBox *geometry.BoundingBox
}

func NewProjectile(position geometry.Vector, velocity geometry.Vector, shooter *Player) (*Projectile, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}
	rotation := velocity.Angle()
	flags := shooter.entity.GetPlayerData().GetFlags()

	entity := &pb.Entity{
		Type:     pb.EntityType_ENTITY_TYPE_PROJECTILE,
		Id:       id,
		Position: &pb.Vector{X: position.X, Y: position.Y},
		Velocity: &pb.Vector{X: velocity.X, Y: velocity.Y},
		Rotation: rotation,
		Data: &pb.Entity_ProjectileData_{
			ProjectileData: &pb.Entity_ProjectileData{
				Flags:    shooter.entity.GetPlayerData().GetFlags(),
				Lifetime: PROJECTILE_LIFETIME,
			},
		},
	}

	p := Projectile{
		entity:   entity,
		position: position,
		velocity: velocity,
		rotation: rotation,
		shooter:  shooter,
	}
	p.boundingBox = geometry.NewBoundingBox(
		&p.position,
		&p.rotation,
		chooseBoundingBoxPoints(AbilityFlag(flags)),
	)
	return &p, nil
}

func (p *Projectile) GetType() pb.EntityType {
	return pb.EntityType_ENTITY_TYPE_PROJECTILE
}

func (p *Projectile) GetEntity() *pb.Entity {
	return p.entity
}

func (p *Projectile) GetID() string {
	return p.entity.Id
}

func (p *Projectile) GetPosition() geometry.Vector {
	return p.position
}

func (p *Projectile) GetVelocity() geometry.Vector {
	return p.velocity
}

func (p *Projectile) GetIsExpired() bool {
	return p.entity.GetProjectileData().Lifetime < 0
}

func (p *Projectile) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Projectile) Update() bool {
	p.position.X += p.velocity.X
	p.position.Y += p.velocity.Y
	p.entity.GetProjectileData().Lifetime--

	// copy to entity
	p.entity.Position.X = p.position.X
	p.entity.Position.Y = p.position.Y

	return true
}

func (p *Projectile) PollNewEntities() []Entity {
	return nil
}

func (p *Projectile) UpdateOnCollision(other Entity) {}

func (p *Projectile) RemoveOnCollision(other Entity) bool {
	switch other.GetType() {
	case pb.EntityType_ENTITY_TYPE_PLAYER:
		p.shooter.entity.GetPlayerData().Score++
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
