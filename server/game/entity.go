package game

import (
	"math"
	"math/rand"
)

type EntityPosition struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
}

func randomEntityPosition() EntityPosition {
	return EntityPosition{
		X:     rand.Float64()*WIDTH - WIDTH/2,
		Y:     rand.Float64()*HEIGHT - WIDTH/2,
		Theta: -math.Pi / 2,
	}
}
