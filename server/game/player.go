package game

import (
	"math"
)

type Player struct {
	Id       string
	Username string
	speed    float64
	position EntityPosition
}

type PlayerState struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
}

func (p *Player) update(mouseX float64, mouseY float64) {
	delta := normalizeAngle(math.Atan2(mouseY, mouseX) - p.position.Theta)
	p.position.Theta = normalizeAngle(p.position.Theta + delta*0.1)

	// TODO: consider non-linear multiplier (e.g. -(x - 1)^2 + 1)
	length := math.Sqrt(mouseX*mouseX + mouseY*mouseY)
	p.position.X += math.Cos(p.position.Theta) * length * p.speed
	p.position.Y += math.Sin(p.position.Theta) * length * p.speed
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
