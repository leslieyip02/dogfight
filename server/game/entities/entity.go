package entities

import (
	"server/game/geometry"
	"server/pb"
	"time"
)

const (
	FPS            = 60
	FRAME_DURATION = time.Second / FPS
)

type Entity interface {
	// wrap over protobuf data
	GetEntityType() pb.EntityType
	GetEntityData() *pb.EntityData

	GetID() string
	GetIsExpired() bool

	// internal representations
	GetPosition() geometry.Vector
	GetVelocity() geometry.Vector
	GetBoundingBox() *geometry.BoundingBox

	Update() bool
	PollNewEntities() []Entity
	UpdateOnCollision(other Entity)
	RemoveOnCollision(other Entity) bool
}
