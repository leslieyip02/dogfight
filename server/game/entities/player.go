package entities

import (
	"math"
	"server/game/geometry"
)

const (
	PLAYER_MAX_SPEED          = 12.0
	PLAYER_ACCELERATION_DECAY = 2.0

	PLAYER_MINIMUM_TURN_RATE = 0.01
	PLAYER_MAX_TURN_RATE     = 0.4
	PLAYER_TURN_RATE_DECAY   = 8.0

	PLAYER_RADIUS = 40.0
)

var playerBoundingBoxPoints = []*geometry.Vector{
	geometry.NewVector(-40, -40),
	geometry.NewVector(40, -40),
	geometry.NewVector(40, 40),
	geometry.NewVector(-40, 40),
}

type Player struct {
	Type     EntityType      `json:"type"`
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Position geometry.Vector `json:"position"`
	Velocity geometry.Vector `json:"velocity"`
	Rotation float64         `json:"rotation"`

	Powerup     *Powerup
	boundingBox *geometry.BoundingBox

	mouseX       float64
	mouseY       float64
	mousePressed bool
}

func NewPlayer(id string, username string) *Player {
	position := *geometry.NewRandomVector(0, 0, SPAWN_AREA_WIDTH, SPAWN_AREA_HEIGHT)
	velocity := *geometry.NewVector(1, 0)
	rotation := 0.0

	p := Player{
		Type:         PlayerEntityType,
		ID:           id,
		Username:     username,
		Position:     position,
		Velocity:     velocity,
		Rotation:     rotation,
		Powerup:      nil,
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
	p.Velocity.X = p.mouseX * PLAYER_MAX_SPEED
	p.Velocity.Y = p.mouseY * PLAYER_MAX_SPEED
	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
	p.Rotation = p.Velocity.Angle()

	// TODO: reimplement
	p.move()
	p.turn()
	return true
}

func (p *Player) PollNewEntities() []Entity {
	if !p.mousePressed {
		return nil
	}
	p.mousePressed = false

	shots := 1
	if p.Powerup != nil && p.Powerup.Ability == MultishotPowerupAbility {
		shots = 3
	}

	projectiles := []Entity{}
	centrePosition := p.Position.Add(p.Velocity.Unit().Multiply(PLAYER_RADIUS + PROJECTILE_RADIUS))
	velocity := p.Velocity.Unit().Multiply(PROJECTILE_SPEED)
	for i := 0; i < shots; i++ {
		offset := (i - shots/2) * 32
		position := centrePosition.Add(geometry.NewVector(
			math.Sin(p.Rotation)*float64(offset),
			math.Cos(p.Rotation)*float64(offset),
		))

		projectile, err := NewProjectile(*position, *velocity)
		if err != nil {
			continue
		}
		projectiles = append(projectiles, projectile)
	}
	return projectiles
}

func (p *Player) RemoveOnCollision(other Entity) bool {
	return other.GetType() != PowerupEntityType
}

func (p *Player) Input(mouseX float64, mouseY float64, mousePressed bool) {
	// mouseX and mouseY are normalized (i.e. range is [0.0, 1.0])
	p.mouseX = mouseX
	p.mouseY = mouseY
	p.mousePressed = p.mousePressed || mousePressed
}

func (p *Player) move() {
	// acceleration := 1 / (1 + ACCELERATION_DECAY*p.speed)
	// throttle := min(math.Sqrt(p.mouseX*p.mouseX+p.mouseY*p.mouseY), 1.0)
	// speedDifference := throttle*MAX_PLAYER_SPEED - p.speed
	// dv := speedDifference * acceleration
	// p.speed += dv

	// p.Position.X += math.Cos(p.Position.Theta) * p.speed
	// p.Position.Y += math.Sin(p.Position.Theta) * p.speed
}

func (p *Player) turn() {
	// turnRate := 1 / (1 + TURN_RATE_DECAY*p.speed) * MAX_TURN_RATE
	// angleDifference := normalizeAngle(math.Atan2(p.mouseY, p.mouseX) - p.Position.Theta)
	// dtheta := math.Copysign(max(math.Abs(angleDifference*turnRate), MINIMUM_TURN_RATE), angleDifference)
	// p.Position.Theta = normalizeAngle(p.Position.Theta + dtheta)
}

// func normalizeAngle(angle float64) float64 {
// 	angle = math.Mod(angle, 2*math.Pi)
// 	if angle > math.Pi {
// 		angle -= 2 * math.Pi
// 	} else if angle < -math.Pi {
// 		angle += 2 * math.Pi
// 	}
// 	return angle
// }
