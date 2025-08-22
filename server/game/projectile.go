package game

import "math"

type Projectile struct {
	Id       string
	position EntityPosition
	speed    float64
	lifetime int
}

func (p *Projectile) update() {
	p.position.X += math.Cos(p.position.Theta) * p.speed
	p.position.Y += math.Sin(p.position.Theta) * p.speed
	p.lifetime--
}
