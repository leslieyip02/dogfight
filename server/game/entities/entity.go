package entities

import (
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
	GetPosition() geometry.Vector
	GetVelocity() geometry.Vector
	GetIsExpired() bool
	GetBoundingBox() *geometry.BoundingBox

	Update() bool
	PollNewEntities() []Entity
	UpdateOnCollision(other Entity)
	RemoveOnCollision(other Entity) bool
}

const (
	AsteroidEntityType   EntityType = "asteroid"
	PlayerEntityType     EntityType = "player"
	ProjectileEntityType EntityType = "projectile"
	PowerupEntityType    EntityType = "powerup"
	MockEntityType       EntityType = "mock"
)
