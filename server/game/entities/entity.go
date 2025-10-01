package entities

import (
	"math"
	"math/rand"
	"server/game/geometry"
	"time"
)

const (
	SPAWN_AREA_WIDTH  = 10000.0
	SPAWN_AREA_HEIGHT = 10000.0

	FPS            = 60
	FRAME_DURATION = time.Second / FPS
)

type EntityType string

type Entity interface {
	GetType() EntityType
	GetID() string
	GetPosition() EntityPosition
	GetIsExpired() bool
	GetBoundingBox() *geometry.BoundingBox

	Update() bool
	PollNewEntities() []Entity
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
		X:     rand.Float64()*SPAWN_AREA_WIDTH - SPAWN_AREA_WIDTH/2,
		Y:     rand.Float64()*SPAWN_AREA_HEIGHT - SPAWN_AREA_WIDTH/2,
		Theta: -math.Pi / 2,
	}
}
