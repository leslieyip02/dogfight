package entities

import (
	"math"
	"server/game/geometry"
)

const (
	PLAYER_MAX_SPEED          = 20.0
	PLAYER_ACCELERATION_DECAY = 8.0
	PLAYER_MAX_TURN_RATE      = 0.8
	PLAYER_TURN_RATE_DECAY    = 4.0

	PLAYER_RADIUS = 40.0
)

var playerBoundingBoxPoints = geometry.NewRectangleHull(80, 80)

type Player struct {
	Type     EntityType      `json:"type"`
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Position geometry.Vector `json:"position"`
	Velocity geometry.Vector `json:"velocity"`
	Rotation float64         `json:"rotation"`
	Score    int             `json:"score"`
	Flags    AbilityFlag     `json:"flags"`

	boundingBox *geometry.BoundingBox

	mouseX       float64
	mouseY       float64
	mousePressed bool
}

func NewPlayer(id string, username string) *Player {
	position := *geometry.NewRandomVector(0, 0, SPAWN_AREA_WIDTH, SPAWN_AREA_HEIGHT)
	velocity := *geometry.NewVector(0, 0)
	rotation := 0.0

	p := Player{
		Type:         PlayerEntityType,
		ID:           id,
		Username:     username,
		Position:     position,
		Velocity:     velocity,
		Rotation:     rotation,
		Score:        0,
		Flags:        0,
		mouseX:       0,
		mouseY:       0,
		mousePressed: false,
	}
	p.boundingBox = geometry.NewBoundingBox(
		&p.Position,
		&p.Rotation,
		&playerBoundingBoxPoints,
	)
	return &p
}

func (p *Player) GetType() EntityType {
	return PlayerEntityType
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) GetPosition() geometry.Vector {
	return p.Position
}

func (p *Player) GetVelocity() geometry.Vector {
	return p.Velocity
}

func (p *Player) GetIsExpired() bool {
	return false
}

func (p *Player) GetBoundingBox() *geometry.BoundingBox {
	return p.boundingBox
}

func (p *Player) Update() bool {
	// TODO: continue iterating on this
	currentSpeed := p.Velocity.Length()
	targetVelocity := geometry.NewVector(p.mouseX, p.mouseY)

	currentAngle := p.Velocity.Angle()
	targetAngle := targetVelocity.Angle()
	turnRate := normalizeAngle(targetAngle - currentAngle)
	maxTurnRate := PLAYER_MAX_TURN_RATE / (1 + PLAYER_TURN_RATE_DECAY*currentSpeed)
	if math.Abs(turnRate) > maxTurnRate {
		turnRate = math.Copysign(maxTurnRate, turnRate)
	}
	p.Rotation = normalizeAngle(currentAngle + turnRate)

	throttle := targetVelocity.Length()
	targetSpeed := throttle * PLAYER_MAX_SPEED
	acceleration := 1 / (1 + PLAYER_ACCELERATION_DECAY*currentSpeed)
	currentSpeed += (targetSpeed - currentSpeed) * acceleration

	p.Velocity.X = math.Cos(p.Rotation) * currentSpeed
	p.Velocity.Y = math.Sin(p.Rotation) * currentSpeed

	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
	return true
}

func (p *Player) PollNewEntities() []Entity {
	if !p.mousePressed {
		return nil
	}
	p.mousePressed = false

	shots := 1
	if isAbilityActive(p.Flags, MultishotAbilityFlag) {
		shots = 3
	}

	projectiles := []Entity{}
	velocity := p.Velocity.Unit().Multiply(PROJECTILE_SPEED)
	for i := range shots {
		offset := float64(i-shots/2) * 32.0
		translated := p.Position.Add(p.Velocity.Normal().Multiply(offset))
		position := translated.Add(p.Velocity.Unit().Multiply(PLAYER_RADIUS*1.1 + PROJECTILE_RADIUS))

		projectile, err := NewProjectile(*position, *velocity, p)
		if err != nil {
			continue
		}
		projectiles = append(projectiles, projectile)
	}
	return projectiles
}

func (p *Player) UpdateOnCollision(other Entity) {
	if other.GetType() == PowerupEntityType {
		powerup := other.(*Powerup)
		p.Flags |= powerup.Ability
	}
}

func (p *Player) RemoveOnCollision(other Entity) bool {
	if other.GetType() == PowerupEntityType {
		return false
	}

	if isAbilityActive(p.Flags, ShieldAbilityFlag) {
		p.Flags ^= ShieldAbilityFlag
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
