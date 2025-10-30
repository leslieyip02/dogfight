package collision

import (
	"fmt"
	"math"
	"server/game/entities"
	"server/game/geometry"
	"slices"
	"testing"
)

var square = []*geometry.Vector{
	geometry.NewVector(-1, -1),
	geometry.NewVector(1, -1),
	geometry.NewVector(1, 1),
	geometry.NewVector(-1, 1),
}

var p1 = geometry.NewVector(0, 0)
var r1 = 0.0
var b1 = geometry.NewBoundingBox(p1, &r1, &square)
var id1 = "1"

var p2 = geometry.NewVector(1, 2)
var r2 = math.Pi / 4
var b2 = geometry.NewBoundingBox(p2, &r2, &square)
var id2 = "2"

var p3 = geometry.NewVector(0, 0)
var r3 = math.Pi / 4
var b3 = geometry.NewBoundingBox(p3, &r3, &square)
var id3 = "3"

var p4 = geometry.NewVector(2, 2)
var r4 = math.Pi / 4
var b4 = geometry.NewBoundingBox(p4, &r4, &square)
var id4 = "4"

var p5 = geometry.NewVector(-1, 10)
var r5 = 0.0
var b5 = geometry.NewBoundingBox(p5, &r5, &square)
var id5 = "5"

var benchmarkEntities map[string]entities.Entity

func init() {
	benchmarkEntities = make(map[string]entities.Entity)
	for i := range 512 {
		id := fmt.Sprintf("%d", i)
		benchmarkEntities[id] = newRandomMockEntity(id)
	}
}

func TestResolveCollisionsLineSweep(t *testing.T) {
	tests := []struct {
		entities map[string]entities.Entity
		want     map[string][]string
	}{
		{
			entities: map[string]entities.Entity{
				id1: &MockEntity{id: id1, position: *p1, boundingBox: b1},
				id2: &MockEntity{id: id2, position: *p2, boundingBox: b2},
			},
			want: map[string][]string{
				id1: {id2},
				id2: {id1},
			},
		},
		{
			entities: map[string]entities.Entity{
				id1: &MockEntity{id: id1, position: *p1, boundingBox: b1},
				id2: &MockEntity{id: id2, position: *p2, boundingBox: b2},
				id3: &MockEntity{id: id3, position: *p3, boundingBox: b3},
				id4: &MockEntity{id: id4, position: *p4, boundingBox: b4},
				id5: &MockEntity{id: id5, position: *p5, boundingBox: b5},
			},
			want: map[string][]string{
				id1: {id2, id3},
				id2: {id1, id4},
				id3: {id1},
				id4: {id2},
				id5: {},
			},
		},
	}
	for _, test := range tests {
		got := make(map[string][]string)
		for id := range test.want {
			got[id] = []string{}
		}
		handleCollision := func(id1 *string, id2 *string) {
			got[*id1] = append(got[*id1], *id2)
			got[*id2] = append(got[*id2], *id1)
		}

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
					t.Errorf("expected collision between %v and %v", id, otherId)
				}
			}
		}
	}
}

func TestGetSortedEdges(t *testing.T) {
	tests := []struct {
		entities map[string]entities.Entity
		want     []Edge
	}{
		{
			entities: map[string]entities.Entity{
				id1: &MockEntity{id: id1, position: *p1, boundingBox: b1},
				id2: &MockEntity{id: id2, position: *p2, boundingBox: b2},
			},
			want: []Edge{
				{id: &id1, x: -1, isLeft: true},
				{id: &id2, x: 1 - math.Sqrt2, isLeft: true},
				{id: &id1, x: 1, isLeft: false},
				{id: &id2, x: 1 + math.Sqrt2, isLeft: false},
			},
		},
		{
			entities: map[string]entities.Entity{
				id1: &MockEntity{id: id1, position: *p1, boundingBox: b1},
				id2: &MockEntity{id: id2, position: *p2, boundingBox: b2},
				id3: &MockEntity{id: id3, position: *p3, boundingBox: b3},
				id4: &MockEntity{id: id4, position: *p4, boundingBox: b4},
				id5: &MockEntity{id: id5, position: *p5, boundingBox: b5},
			},
			want: []Edge{
				{id: &id5, x: -2, isLeft: true},
				{id: &id3, x: -math.Sqrt2, isLeft: true},
				{id: &id1, x: -1, isLeft: true},
				{id: &id2, x: 1 - math.Sqrt2, isLeft: true},
				{id: &id5, x: 0, isLeft: false},
				{id: &id4, x: 2 - math.Sqrt2, isLeft: true},
				{id: &id1, x: 1, isLeft: false},
				{id: &id3, x: math.Sqrt2, isLeft: false},
				{id: &id2, x: 1 + math.Sqrt2, isLeft: false},
				{id: &id4, x: 2 + math.Sqrt2, isLeft: false},
			},
		},
		{
			entities: map[string]entities.Entity{
				id1: &MockEntity{id: id1, position: *p1, boundingBox: b1},
				id2: &MockEntity{id: id2, position: *p2, boundingBox: b1},
			},
			want: []Edge{
				{id: &id1, x: -1, isLeft: true},
				{id: &id2, x: -1, isLeft: true},
				{id: &id1, x: 1, isLeft: false},
				{id: &id2, x: 1, isLeft: false},
			},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("sort edges of %v", test.entities), func(t *testing.T) {
			got := getSortedEdges(&test.entities)
			if len(got) != len(test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
			for i := 0; i < len(got); i++ {
				// as long as x and isLeft is correct, whether id matches doesn't matter
				// because the edge will be included in the comparison window anyways
				if got[i].isLeft != test.want[i].isLeft {
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
