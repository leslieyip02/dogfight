package game

import (
	"math"
	"math/rand"
)

type EntityType string

type Entity interface {
	GetType() EntityType
	GetID() string
	GetPosition() EntityPosition

	Update(g *Game)
	GetIsExpired() bool
}

const (
	PlayerEntityType     EntityType = "player"
	ProjectileEntityType EntityType = "projectile"
	PowerupEntityType    EntityType = "powerup"
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
