package collision

import (
	"fmt"
	"math"
	"math/rand/v2"
	"server/internal/game/entities"
	"server/internal/game/geometry"
	"slices"
	"testing"
)

func randomEntity(id string) *entities.MockEntity {
	position := *geometry.NewRandomVector(-100, 100, -100, 100)
	rotation := rand.Float64() * math.Pi * 2
	points := geometry.NewRandomConvexHull(8, 16, 8, 32)
	return entities.NewMockEntity(id, position.X, position.Y, rotation, points)
}

var square = geometry.NewRectangleHull(2, 2)

var e1 = entities.NewMockEntity("1", 0, 0, 0, square)
var e2 = entities.NewMockEntity("2", 1, 2, 0, square)
var e3 = entities.NewMockEntity("3", 0, 0, math.Pi/4, square)
var e4 = entities.NewMockEntity("4", 1, 2, math.Pi/4, square)
var e5 = entities.NewMockEntity("5", -1, 10, math.Pi/4, square)
var e6 = entities.NewMockEntity("6", 2, 2, 0, square)

var benchmarkEntities map[string]entities.Entity

func mockCollisionHandler(id1 *string, id2 *string) {}

func init() {
	benchmarkEntities = make(map[string]entities.Entity)
	for i := range 512 {
		id := fmt.Sprintf("%d", i)
		benchmarkEntities[id] = randomEntity(id)
	}
}

func TestResolveCollisionsLineSweep(t *testing.T) {
	tests := map[string]struct {
		entities map[string]entities.Entity
		want     map[string][]string
	}{
		"ResolveCollisionsLineSweep": {
			map[string]entities.Entity{
				"1": e1,
				"4": e4,
			},
			map[string][]string{
				"1": {"4"},
				"4": {"1"},
			},
		},
		"ResolveCollisionsLineSweep with overlapping edge": {
			map[string]entities.Entity{
				"1": e1,
				"6": e6,
			},
			map[string][]string{
				"1": {"6"},
				"6": {"1"},
			},
		},
		"ResolveCollisionsLineSweep with empty collision list": {
			map[string]entities.Entity{
				"1": e1,
				"2": e2,
				"5": e5,
			},
			map[string][]string{
				"1": {"2"},
				"2": {"1"},
				"5": {},
			},
		},
		"ResolveCollisionsLineSweep with more entities": {
			map[string]entities.Entity{
				"1": e1,
				"2": e2,
				"3": e3,
				"4": e4,
			},
			map[string][]string{
				"1": {"2", "3", "4"},
				"2": {"1", "3", "4"},
				"3": {"1", "2"},
				"4": {"1", "2"},
			},
		},
	}

	for desc, test := range tests {
		got := make(map[string][]string)
		for id := range test.want {
			got[id] = []string{}
		}
		handleCollision := func(id1 *string, id2 *string) {
			got[*id1] = append(got[*id1], *id2)
			got[*id2] = append(got[*id2], *id1)
		}

		title := desc
		t.Run(title, func(t *testing.T) {
			ResolveCollisionsLineSweep(&test.entities, handleCollision)
			if len(got) != len(test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
			for id := range test.want {
				if len(test.want[id]) != len(got[id]) {
					t.Errorf("want %v but got %v", test.want[id], got[id])
				}
				for _, otherId := range test.want[id] {
					if !slices.Contains(got[id], otherId) {
						t.Errorf("want collision between %v and %v", id, otherId)
					}
				}
			}
		})
	}
}

func TestGetSortedEdges(t *testing.T) {
	tests := map[string]struct {
		entities map[string]entities.Entity
		want     []Edge
	}{
		"getSortedEdges": {
			map[string]entities.Entity{
				"1": e1,
				"4": e4,
			},
			[]Edge{
				{id: &e1.Id, x: -1, isLeft: true},
				{id: &e4.Id, x: 1 - math.Sqrt2, isLeft: true},
				{id: &e1.Id, x: 1, isLeft: false},
				{id: &e4.Id, x: 1 + math.Sqrt2, isLeft: false},
			},
		},
		"getSortedEdges with overlapping edge": {
			map[string]entities.Entity{
				"1": e1,
				"6": e6,
			},
			[]Edge{
				{id: &e1.Id, x: -1, isLeft: true},
				{id: &e6.Id, x: 1, isLeft: true},
				{id: &e1.Id, x: 1, isLeft: false},
				{id: &e6.Id, x: 3, isLeft: false},
			},
		},
		"getSortedEdges with more entities": {
			map[string]entities.Entity{
				"1": e1,
				"2": e2,
				"3": e3,
				"4": e4,
				"5": e5,
				"6": e6,
			},
			[]Edge{
				{id: &e5.Id, x: -1 - math.Sqrt2, isLeft: true},
				{id: &e3.Id, x: -math.Sqrt2, isLeft: true},
				{id: &e1.Id, x: -1, isLeft: true},
				{id: &e4.Id, x: 1 - math.Sqrt2, isLeft: true},
				{id: &e2.Id, x: 0, isLeft: true},
				{id: &e5.Id, x: -1 + math.Sqrt2, isLeft: false},
				{id: &e6.Id, x: 1, isLeft: true},
				{id: &e1.Id, x: 1, isLeft: false},
				{id: &e3.Id, x: math.Sqrt2, isLeft: false},
				{id: &e2.Id, x: 2, isLeft: false},
				{id: &e4.Id, x: 1 + math.Sqrt2, isLeft: false},
				{id: &e6.Id, x: 3, isLeft: false},
			},
		},
	}

	for desc, test := range tests {
		title := desc
		t.Run(title, func(t *testing.T) {
			got := getSortedEdges(&test.entities)
			if len(got) != len(test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
			for i := range got {
				if *got[i].id != *test.want[i].id || got[i].isLeft != test.want[i].isLeft {
					t.Errorf("want %v but got %v", test.want[i], got[i])
				}
				if math.Abs(got[i].x-test.want[i].x) > geometry.EPSILON {
					t.Errorf("want %v but got %v", test.want[i], got[i])
				}
			}
		})
	}
}

func BenchmarkResolveCollisionsLineSweep(b *testing.B) {
	for b.Loop() {
		ResolveCollisionsLineSweep(&benchmarkEntities, mockCollisionHandler)
	}
}

func BenchmarkResolveCollisionsNaive(b *testing.B) {
	for b.Loop() {
		resolveCollisionsNaive(&benchmarkEntities, mockCollisionHandler)
	}
}
