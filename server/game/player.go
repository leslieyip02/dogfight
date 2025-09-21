package game

import (
	"math"
	"server/utils"
)

const (
	ACCELERATION_DECAY = 2.0
	MAX_PLAYER_SPEED   = 12.0

	TURN_RATE_DECAY   = 8.0
	MAX_TURN_RATE     = 0.4
	MINIMUM_TURN_RATE = 0.01

	PLAYER_RADIUS = 40.0
)

type Player struct {
	Id       string         `json:"id"`
	Username string         `json:"username"`
	Position EntityPosition `json:"position"`
	speed    float64
	powerup  *Powerup
}

func (p *Player) input(data InputEventData, game *Game) {
	p.updatePosition(data)
	p.shootProjectiles(data, game)
}

func (p *Player) updatePosition(data InputEventData) {
	// mouseX and mouseY are normalized (i.e. range is [0.0, 1.0])
	p.move(data.MouseX, data.MouseY)
	p.turn(data.MouseX, data.MouseY)
}

func (p *Player) move(mouseX float64, mouseY float64) {
	acceleration := 1 / (1 + ACCELERATION_DECAY*p.speed)
	throttle := min(math.Sqrt(mouseX*mouseX+mouseY*mouseY), 1.0)
	speedDifference := throttle*MAX_PLAYER_SPEED - p.speed
	dv := speedDifference * acceleration
	p.speed += dv

	p.Position.X += math.Cos(p.Position.Theta) * p.speed
	p.Position.Y += math.Sin(p.Position.Theta) * p.speed
}

func (p *Player) turn(mouseX float64, mouseY float64) {
	turnRate := 1 / (1 + TURN_RATE_DECAY*p.speed) * MAX_TURN_RATE
	angleDifference := normalizeAngle(math.Atan2(mouseY, mouseX) - p.Position.Theta)
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

func (p *Player) shootProjectiles(data InputEventData, game *Game) {
	if !data.MousePressed {
		return
	}

	var shots int
	if p.powerup == nil {
		shots = 1
	} else {
		shots = 3
	}

	centerX := p.Position.X + math.Cos(p.Position.Theta)*(PLAYER_RADIUS+PROJECTILE_RADIUS)
	centerY := p.Position.Y + math.Sin(p.Position.Theta)*(PLAYER_RADIUS+PROJECTILE_RADIUS)
	for i := 0; i < shots; i++ {
		id, err := utils.NewShortId()
		if err != nil {
			continue
		}

		offset := (i - shots/2) * 32
		position := EntityPosition{
			X:     centerX + math.Sin(p.Position.Theta)*float64(offset),
			Y:     centerY - math.Cos(p.Position.Theta)*float64(offset),
			Theta: p.Position.Theta,
		}
		projectile := Projectile{
			Id:       id,
			Position: position,
			speed:    PROJECTILE_SPEED,
			lifetime: PROJECTILE_LIFETIME,
		}
		game.projectiles[id] = &projectile
	}
}
