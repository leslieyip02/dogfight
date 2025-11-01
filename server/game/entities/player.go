package entities

import (
	"math"
	"server/game/geometry"
	"server/id"
	"server/pb"
)

const (
	PLAYER_MAX_SPEED              = 20.0
	PLAYER_ACCELERATION_DECAY     = 8.0
	PLAYER_MAX_TURN_RATE          = 0.8
	PLAYER_TURN_RATE_DECAY_FACTOR = 4.0
	PLAYER_RADIUS                 = 40.0
)

var playerBoundingBoxPoints = geometry.NewRectangleHull(
	PLAYER_RADIUS*2,
	PLAYER_RADIUS*2,
)

// A Player is the spaceship which players control.
//
// Behaviors:
//   - destroyed on impact with an asteroid, projectile or another player
//   - can shoot projectiles
//   - can collect powerups
//   - takes 1 shot to destroy
type Player struct {
	entityData *pb.EntityData

	// Internal duplicates of EntityData state.
	position geometry.Vector
	velocity geometry.Vector
	rotation float64

	boundingBox  *geometry.BoundingBox
	mouseX       float64
	mouseY       float64
	mousePressed bool
}

func newPlayer(
	id string,
	position geometry.Vector,
	velocity geometry.Vector,
	rotation float64,
	username string,
) *Player {
	mouseX := 0.0
	mouseY := 0.0
	mousePressed := false

	entityData := &pb.EntityData{
		Type:     pb.EntityType_ENTITY_TYPE_PLAYER,
		Id:       id,
		Position: position.ToPb(),
		Velocity: velocity.ToPb(),
		Rotation: rotation,
		Data: &pb.EntityData_PlayerData_{
			PlayerData: &pb.EntityData_PlayerData{
				Username: username,
				Score:    0,
				Flags:    0,
			},
		},
	}

	p := Player{
		entityData:   entityData,
		position:     position,
		velocity:     velocity,
		rotation:     rotation,
		mouseX:       mouseX,
		mouseY:       mouseY,
		mousePressed: mousePressed,
	}
	p.boundingBox = geometry.NewBoundingBox(
		&p.position,
		&p.rotation,
		&playerBoundingBoxPoints,
	)
	return &p
}

func (p *Player) GetEntityType() pb.EntityType {
	return pb.EntityType_ENTITY_TYPE_PLAYER
}

func (p *Player) GetEntityData() *pb.EntityData {
	return p.entityData
}

func (p *Player) GetId() string {
	return p.entityData.Id
}

func (p *Player) GetPosition() geometry.Vector {
	return p.position
}

func (p *Player) GetVelocity() geometry.Vector {
	return p.velocity
}

func (p *Player) GetIsExpired() bool {
	return false
}

func (p *Player) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Player) Update() bool {
	// TODO: continue iterating on this
	targetVelocity := geometry.NewVector(p.mouseX, p.mouseY)

	speed := accelerate(p.velocity, targetVelocity)
	p.velocity.X = math.Cos(p.rotation) * speed
	p.velocity.Y = math.Sin(p.rotation) * speed

	p.position.X += p.velocity.X
	p.position.Y += p.velocity.Y
	p.rotation = rotate(p.velocity, targetVelocity)

	p.SyncEntityData()
	return true
}

func (p *Player) PollNewEntities() []Entity {
	if !p.mousePressed {
		return nil
	}
	p.mousePressed = false

	return p.spawnProjectiles()
}

func (p *Player) UpdateOnCollision(other Entity) {
	if other.GetEntityType() == pb.EntityType_ENTITY_TYPE_POWERUP {
		powerup := other.(*Powerup)
		ability := powerup.entityData.GetPowerupData().Ability
		p.entityData.GetPlayerData().Flags |= ability
	}
}

func (p *Player) RemoveOnCollision(other Entity) bool {
	if other.GetEntityType() == pb.EntityType_ENTITY_TYPE_POWERUP {
		return false
	}

	flags := AbilityFlag(p.entityData.GetPlayerData().Flags)
	if isAbilityActive(flags, ShieldAbilityFlag) {
		p.entityData.GetPlayerData().Flags ^= uint32(ShieldAbilityFlag)
		return false
	}

	return true
}

func (p *Player) Input(mouseX float64, mouseY float64, mousePressed bool) {
	// mouseX and mouseY are normalized (i.e. range is [0.0, 1.0])
	p.mouseX = mouseX
	p.mouseY = mouseY
	p.mousePressed = p.mousePressed || mousePressed
}

// spawnProjectiles creates a volley of projectiles based on the player's
// active abilities.
func (p *Player) spawnProjectiles() []Entity {
	flags := AbilityFlag(p.entityData.GetPlayerData().Flags)
	shots := 1
	if isAbilityActive(flags, MultishotAbilityFlag) {
		shots = 3
	}

	projectiles := []Entity{}
	for i := range shots {
		offset := float64(i-shots/2) * 32.0

		projectile, err := p.spawnProjectile(offset)
		if err != nil {
			continue
		}

		projectiles = append(projectiles, projectile)
	}
	return projectiles
}

// spawnProjectile creates a single projectile with a given offset. Offset is
// the perpendicular distance between the player's velocity and position of the
// projectile.
func (p *Player) spawnProjectile(offset float64) (*Projectile, error) {
	id, err := id.NewShortId()
	if err != nil {
		return nil, err
	}

	translated := p.position.Add(p.velocity.Normal().Multiply(offset))
	position := translated.Add(
		p.velocity.Unit().Multiply(PLAYER_RADIUS*1.1 + PROJECTILE_RADIUS),
	)
	velocity := p.velocity.Unit().Multiply(PROJECTILE_SPEED)

	return newProjectile(
		id,
		*position,
		*velocity,
		AbilityFlag(p.entityData.GetPlayerData().Flags),
		p.projectileOnRemove,
	), nil
}

// projectileOnRemove is a callback when the projectile is destroyed. It
// increments the player's score if the shot hit another player.
func (p *Player) projectileOnRemove(other *Entity) {
	if other == nil {
		return
	}

	if (*other).GetEntityType() == pb.EntityType_ENTITY_TYPE_PLAYER {
		p.entityData.GetPlayerData().Score++
	}
}

// rotate computes a new rotation angle based on the current and target
// velocities. The angle is between velocity and the positive x-axis.
func rotate(
	currentVelocity geometry.Vector,
	targetVelocity *geometry.Vector,
) float64 {
	currentSpeed := currentVelocity.Length()
	currentAngle := currentVelocity.Angle()
	targetAngle := targetVelocity.Angle()

	rotation := normalizeAngle(targetAngle - currentAngle)
	decay := 1 / (1 + PLAYER_TURN_RATE_DECAY_FACTOR*currentSpeed)
	return normalizeAngle(currentAngle + rotation*decay)
}

// accelerate computes a speed based on the current and target velocities.
func accelerate(
	currentVelocity geometry.Vector,
	targetVelocity *geometry.Vector,
) float64 {
	currentSpeed := currentVelocity.Length()
	targetSpeed := targetVelocity.Length() * PLAYER_MAX_SPEED

	acceleration := targetSpeed - currentSpeed
	decay := 1 / (1 + PLAYER_ACCELERATION_DECAY*currentSpeed)
	return currentSpeed + acceleration*decay
}

// normalizeAngle clamps angle to the range [-π, π].
func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle > math.Pi {
		angle -= 2 * math.Pi
	} else if angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

func (p *Player) SyncEntityData() {
	p.entityData.Position.X = p.position.X
	p.entityData.Position.Y = p.position.Y
	p.entityData.Velocity.X = p.velocity.X
	p.entityData.Velocity.Y = p.velocity.Y
	p.entityData.Rotation = p.rotation
}
