package entities

import (
	"server/game/geometry"
	"server/pb"
	"time"
)

const (
	SPAWN_AREA_WIDTH  = 10000.0
	SPAWN_AREA_HEIGHT = 10000.0

	FPS            = 60
	FRAME_DURATION = time.Second / FPS
)

type Entity interface {
	GetType() pb.EntityType
	GetEntity() *pb.Entity

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
