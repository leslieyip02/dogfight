package collision

import (
	"math"
	"math/rand/v2"
	"server/game/entities"
	"server/game/geometry"
)

type MockEntity struct {
	id          string
	position    geometry.Vector
	boundingBox *geometry.BoundingBox
}

func (e *MockEntity) GetType() entities.EntityType {
	return entities.MockEntityType
}

func (e *MockEntity) GetID() string {
	return e.id
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

func (e *MockEntity) PollNewEntities() []entities.Entity {
	return nil
}

func (e *MockEntity) UpdateOnCollision(other entities.Entity) {}

func (e *MockEntity) RemoveOnCollision(other entities.Entity) bool {
	return false
}

func newRandomMockEntity(id string) *MockEntity {
	position := *geometry.NewRandomVector(-100, 100, -100, 100)
	rotation := rand.Float64() * math.Pi * 2
	points := geometry.NewRandomConvexHull(8, 16, 8, 32)
	boundingBox := geometry.NewBoundingBox(&position, &rotation, &points)

	return &MockEntity{
		id:          id,
		position:    position,
		boundingBox: boundingBox,
	}
}

func mockCollisionHandler(id1 *string, id2 *string) {}
