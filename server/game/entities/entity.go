package entities

import (
	"server/game/geometry"
	"server/pb"
)

// An Entity is a struct that represents a game entity. It wraps over
// pb.EntityData, and provides an interface to interact with other entities.
type Entity interface {
	// Entity wraps over pb.EntityData, which makes serialization convenient.
	GetEntityType() pb.EntityType
	GetEntityData() *pb.EntityData

	GetId() string
	GetIsExpired() bool

	// Getters for the internal representation of the entity's state. Even
	// though Vector is also defined in Protobuf, the internal representation
	// from the geometry package provides more utilities.
	GetPosition() geometry.Vector
	GetVelocity() geometry.Vector
	GetBoundingBox() *geometry.BoundingBox

	Update() bool

	// PollNewEntities returns any entities created by the entity.
	// (e.g. player shooting a projectile).
	PollNewEntities() []Entity

	// UpdateOnCollision processes updates after a collision with other.
	UpdateOnCollision(other Entity)

	// RemoveOnCollision reports if the entity should be removed after a
	// collision with other.
	RemoveOnCollision(other Entity) bool

	// SyncEntityData syncs the entity's internal game state with the state
	// stored in pb.EntityData by writing updated values into pb.EntityData.
	SyncEntityData()
}
