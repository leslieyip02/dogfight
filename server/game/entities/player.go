package entities

import (
	"math"
	"server/game/geometry"
	"server/pb"
)

const (
	PLAYER_MAX_SPEED          = 20.0
	PLAYER_ACCELERATION_DECAY = 8.0
	PLAYER_MAX_TURN_RATE      = 0.8
	PLAYER_TURN_RATE_DECAY    = 4.0
	PLAYER_RADIUS             = 40.0
)

var playerBoundingBoxPoints = geometry.NewRectangleHull(PLAYER_RADIUS*2, PLAYER_RADIUS*2)

type Player struct {
	entityData *pb.EntityData

	// state
	position    geometry.Vector
	velocity    geometry.Vector
	rotation    float64
	boundingBox *geometry.BoundingBox

	// input
	mouseX       float64
	mouseY       float64
	mousePressed bool
}

func NewPlayer(id string, username string) *Player {
	position := *geometry.NewRandomVector(0, 0, SPAWN_AREA_WIDTH, SPAWN_AREA_HEIGHT)
	velocity := *geometry.NewVector(0, 0)
	rotation := 0.0

	entity := &pb.EntityData{
		Type:     pb.EntityType_ENTITY_TYPE_PLAYER,
		Id:       id,
		Position: &pb.Vector{X: position.X, Y: position.Y},
		Velocity: &pb.Vector{X: velocity.X, Y: velocity.Y},
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
		entityData:   entity,
		position:     position,
		velocity:     velocity,
		rotation:     rotation,
		mouseX:       0,
		mouseY:       0,
		mousePressed: false,
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

func (p *Player) GetID() string {
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
	currentSpeed := p.velocity.Length()
	targetVelocity := geometry.NewVector(p.mouseX, p.mouseY)

	currentAngle := p.velocity.Angle()
	targetAngle := targetVelocity.Angle()
	turnRate := normalizeAngle(targetAngle - currentAngle)
	maxTurnRate := PLAYER_MAX_TURN_RATE / (1 + PLAYER_TURN_RATE_DECAY*currentSpeed)
	if math.Abs(turnRate) > maxTurnRate {
		turnRate = math.Copysign(maxTurnRate, turnRate)
	}
	p.rotation = normalizeAngle(currentAngle + turnRate)

	throttle := targetVelocity.Length()
	targetSpeed := throttle * PLAYER_MAX_SPEED
	acceleration := 1 / (1 + PLAYER_ACCELERATION_DECAY*currentSpeed)
	currentSpeed += (targetSpeed - currentSpeed) * acceleration

	p.velocity.X = math.Cos(p.rotation) * currentSpeed
	p.velocity.Y = math.Sin(p.rotation) * currentSpeed

	p.position.X += p.velocity.X
	p.position.Y += p.velocity.Y

	// copy to entity
	p.entityData.Position.X = p.position.X
	p.entityData.Position.Y = p.position.Y
	p.entityData.Velocity.X = p.velocity.X
	p.entityData.Velocity.Y = p.velocity.Y
	p.entityData.Rotation = p.rotation

	return true
}

func (p *Player) PollNewEntities() []Entity {
	if !p.mousePressed {
		return nil
	}
	p.mousePressed = false

	shots := 1
	flags := AbilityFlag(p.entityData.GetPlayerData().Flags)
	if isAbilityActive(flags, MultishotAbilityFlag) {
		shots = 3
	}

	projectiles := []Entity{}
	velocity := p.velocity.Unit().Multiply(PROJECTILE_SPEED)
	for i := range shots {
		offset := float64(i-shots/2) * 32.0
		translated := p.position.Add(p.velocity.Normal().Multiply(offset))
		position := translated.Add(p.velocity.Unit().Multiply(PLAYER_RADIUS*1.1 + PROJECTILE_RADIUS))

		projectile, err := NewProjectile(*position, *velocity, p)
		if err != nil {
			continue
		}
		projectiles = append(projectiles, projectile)
	}
	return projectiles
}

func (p *Player) UpdateOnCollision(other Entity) {
	if other.GetEntityType() == pb.EntityType_ENTITY_TYPE_POWERUP {
		powerup := other.(*Powerup)
		p.entityData.GetPlayerData().Flags |= powerup.entityData.GetPowerupData().Ability
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

func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle > math.Pi {
		angle -= 2 * math.Pi
	} else if angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}
