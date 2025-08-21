package game

type Player struct {
	Id       string
	Username string
	x        float64
	y        float64
	theta    float64
	speed    float64
}

type PlayerState struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
}

func (p *Player) getState() PlayerState {
	return PlayerState{
		X:     p.x,
		Y:     p.y,
		Theta: p.theta,
	}
}
