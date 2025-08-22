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
}

func (p *Player) input(data InputEventData) *Projectile {
	p.updatePosition(data.MouseX, data.MouseY)
	if data.MousePressed {
		return p.shootProjectile()
	} else {
		return nil
	}
}

func (p *Player) updatePosition(mouseX float64, mouseY float64) {
	delta := normalizeAngle(math.Atan2(mouseY, mouseX) - p.Position.Theta)
	p.Position.Theta = normalizeAngle(p.Position.Theta + delta*0.1)

	// TODO: consider non-linear multiplier (e.g. -(x - 1)^2 + 1)
	length := math.Sqrt(mouseX*mouseX + mouseY*mouseY)
	p.Position.X += math.Cos(p.Position.Theta) * length * p.speed
	p.Position.Y += math.Sin(p.Position.Theta) * length * p.speed
}

func (p *Player) shootProjectile() *Projectile {
	id, err := utils.NewShortId()
	if err != nil {
		return nil
	}

	projectile := Projectile{
		Id:       id,
		position: p.Position,
		speed:    8.0,
		lifetime: 1 * FPS,
	}
	return &projectile
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
