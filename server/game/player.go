package game

import (
	"math"
	"server/utils"
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
	p.move(data.MouseX, data.MouseY)
	p.turn(data.MouseX, data.MouseY)
}

func (p *Player) move(mouseX float64, mouseY float64) {
	acceleration := 1 / (1 + ACCELERATION_DECAY*p.speed)
	throttle := math.Sqrt(mouseX*mouseX+mouseY*mouseY) / math.Sqrt(2)
	speedDifference := throttle*MAX_PLAYER_SPEED - p.speed
	p.speed += speedDifference * acceleration

	p.Position.X += math.Cos(p.Position.Theta) * p.speed
	p.Position.Y += math.Sin(p.Position.Theta) * p.speed
}

func (p *Player) turn(mouseX float64, mouseY float64) {
	turnRate := 1 / (1 + TURN_RATE_DECAY*p.speed) * MAX_TURN_RATE
	angleDifference := normalizeAngle(math.Atan2(mouseY, mouseX) - p.Position.Theta)
	p.Position.Theta = normalizeAngle(p.Position.Theta + angleDifference*turnRate)
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
			lifetime: 1 * FPS,
		}
		game.projectiles[id] = &projectile
	}
}
