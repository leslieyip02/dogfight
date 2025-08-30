package game

import "math"

type Projectile struct {
	Id       string         `json:"id"`
	Position EntityPosition `json:"position"`
	speed    float64
	lifetime int
}

func (p *Projectile) update() {
	p.Position.X += math.Cos(p.Position.Theta) * p.speed
	p.Position.Y += math.Sin(p.Position.Theta) * p.speed
	p.lifetime--
}
