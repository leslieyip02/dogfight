package game

import (
	"math"
	"server/utils"
)

type Player struct {
	Id       string
	Username string
	speed    float64
	position EntityPosition
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
	delta := normalizeAngle(math.Atan2(mouseY, mouseX) - p.position.Theta)
	p.position.Theta = normalizeAngle(p.position.Theta + delta*0.1)

	// TODO: consider non-linear multiplier (e.g. -(x - 1)^2 + 1)
	length := math.Sqrt(mouseX*mouseX + mouseY*mouseY)
	p.position.X += math.Cos(p.position.Theta) * length * p.speed
	p.position.Y += math.Sin(p.position.Theta) * length * p.speed
}

func (p *Player) shootProjectile() *Projectile {
	id, err := utils.NewShortId()
	if err != nil {
		return nil
	}

	projectile := Projectile{
		Id:       id,
		position: p.position,
		speed:    MAX_SPEED,
	}
	return &projectile
}

func (p *Player) shoot() Projectile {
	position := EntityPosition{
		X:     p.position.X,
		Y:     p.position.Y,
		Theta: p.position.Theta,
	}
	return Projectile{
		Id:       "",
		speed:    0,
		position: position,
		lifetime: 3 * FPS,
	}
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
