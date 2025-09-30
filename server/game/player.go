package game

import (
	"math"
	"server/game/geometry"
)

const (
	ACCELERATION_DECAY = 2.0
	MAX_PLAYER_SPEED   = 12.0

	TURN_RATE_DECAY   = 8.0
	MAX_TURN_RATE     = 0.4
	MINIMUM_TURN_RATE = 0.01

	PLAYER_RADIUS = 40.0
)

var playerBoundingBox = geometry.NewBoundingBox(&[]geometry.Vector{
	geometry.NewVector(-40, -40),
	geometry.NewVector(40, -40),
	geometry.NewVector(40, 40),
	geometry.NewVector(-40, 40),
})

type Player struct {
	Type     EntityType     `json:"type"`
	ID       string         `json:"id"`
	Username string         `json:"username"`
	Position EntityPosition `json:"position"`

	boundingBox *geometry.BoundingBox
	speed       float64
	powerup     *Powerup

	mouseX       float64
	mouseY       float64
	mousePressed bool
}

func NewPlayer(id string, username string) *Player {
	p := Player{
		Type:         PlayerEntityType,
		ID:           id,
		Username:     username,
		Position:     randomEntityPosition(),
		speed:        MAX_PLAYER_SPEED,
		powerup:      nil,
		mouseX:       0,
		mouseY:       0,
		mousePressed: false,
		boundingBox:  &playerBoundingBox,
	}
	return &p
}

func (p *Player) GetType() EntityType {
	return PlayerEntityType
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) GetPosition() EntityPosition {
	return p.Position
}

func (p *Player) GetIsExpired() bool {
	return false
}

func (p *Player) GetBoundingBox() *geometry.BoundingBox {
	return playerBoundingBox.Transform(p.Position.X, p.Position.Y, p.Position.Theta)
}

func (p *Player) Update(g *Game) bool {
	p.move()
	p.turn()
	p.shootProjectiles(g)
	return true
}

func (p *Player) input(data InputEventData) {
	// mouseX and mouseY are normalized (i.e. range is [0.0, 1.0])
	p.mouseX = data.MouseX
	p.mouseY = data.MouseY
	p.mousePressed = p.mousePressed || data.MousePressed
}

func (p *Player) move() {
	acceleration := 1 / (1 + ACCELERATION_DECAY*p.speed)
	throttle := min(math.Sqrt(p.mouseX*p.mouseX+p.mouseY*p.mouseY), 1.0)
	speedDifference := throttle*MAX_PLAYER_SPEED - p.speed
	dv := speedDifference * acceleration
	p.speed += dv

	p.Position.X += math.Cos(p.Position.Theta) * p.speed
	p.Position.Y += math.Sin(p.Position.Theta) * p.speed
}

func (p *Player) turn() {
	turnRate := 1 / (1 + TURN_RATE_DECAY*p.speed) * MAX_TURN_RATE
	angleDifference := normalizeAngle(math.Atan2(p.mouseY, p.mouseX) - p.Position.Theta)
	dtheta := math.Copysign(max(math.Abs(angleDifference*turnRate), MINIMUM_TURN_RATE), angleDifference)
	p.Position.Theta = normalizeAngle(p.Position.Theta + dtheta)
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

func (p *Player) shootProjectiles(g *Game) {
	if !p.mousePressed {
		return
	}
	p.mousePressed = false

	var shots int
	if p.powerup == nil {
		shots = 1
	} else {
		shots = 3
	}

	centerX := p.Position.X + math.Cos(p.Position.Theta)*(PLAYER_RADIUS+PROJECTILE_RADIUS+p.speed)
	centerY := p.Position.Y + math.Sin(p.Position.Theta)*(PLAYER_RADIUS+PROJECTILE_RADIUS+p.speed)
	for i := 0; i < shots; i++ {
		offset := (i - shots/2) * 32
		position := EntityPosition{
			X:     centerX + math.Sin(p.Position.Theta)*float64(offset),
			Y:     centerY - math.Cos(p.Position.Theta)*float64(offset),
			Theta: p.Position.Theta,
		}
		projectile, err := NewProjectile(position)
		if err != nil {
			continue
		}
		g.entities[projectile.ID] = projectile
	}
}
