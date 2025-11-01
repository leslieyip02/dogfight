package entities

import (
	"server/game/geometry"
	"server/pb"
)

// A MockEntity is a mock used for tests.
type MockEntity struct {
	Id          string
	position    geometry.Vector
	boundingBox *geometry.BoundingBox
}

func NewMockEntity(
	id string,
	x float64,
	y float64,
	rotation float64,
	points []*geometry.Vector,
) *MockEntity {
	position := geometry.NewVector(x, y)
	boundingBox := geometry.NewBoundingBox(position, &rotation, &points)
	return &MockEntity{Id: id, position: *position, boundingBox: boundingBox}
}

func (e *MockEntity) GetEntityType() pb.EntityType {
	return pb.EntityType_ENTITY_TYPE_MOCK
}

func (e *MockEntity) GetEntityData() *pb.EntityData {
	return nil
}

func (e *MockEntity) GetId() string {
	return e.Id
}

func (e *MockEntity) GetPosition() geometry.Vector {
	return e.position
}

func (e *MockEntity) GetVelocity() geometry.Vector {
	return *geometry.NewVector(0, 0)
}

func (e *MockEntity) GetIsExpired() bool {
	return false
}

func (e *MockEntity) GetBoundingBox() *geometry.BoundingBox {
	return e.boundingBox
}

func (e *MockEntity) Update() bool {
	return false
}

func (e *MockEntity) PollNewEntities() []Entity {
	return nil
}

func (e *MockEntity) UpdateOnCollision(other Entity) {}

func (e *MockEntity) RemoveOnCollision(other Entity) bool {
	return false
}

func (e *MockEntity) SyncEntityData() {}
