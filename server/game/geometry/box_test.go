package geometry

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

var b1 = BoundingBox{
	origin: Vector{x: 0, y: 0},
	theta:  0,
	points: &[]Vector{
		{x: -1, y: -1},
		{x: 1, y: -1},
		{x: 1, y: 1},
		{x: -1, y: 1},
	},
}

var b2 = BoundingBox{
	origin: Vector{x: 1, y: 2},
	theta:  math.Pi / 4,
	points: &[]Vector{
		{x: -1, y: -1},
		{x: 1, y: -1},
		{x: 1, y: 1},
		{x: -1, y: 1},
	},
}

var b3 = BoundingBox{
	origin: Vector{x: 0, y: 0},
	theta:  math.Pi / 4,
	points: &[]Vector{
		{x: -1, y: -1},
		{x: 1, y: -1},
		{x: 1, y: 1},
		{x: -1, y: 1},
	},
}

var b4 = BoundingBox{
	origin: Vector{x: 2, y: 2},
	theta:  math.Pi / 4,
	points: &[]Vector{
		{x: -1, y: -1},
		{x: 1, y: -1},
		{x: 1, y: 1},
		{x: -1, y: 1},
	},
}

func TestNormals(t *testing.T) {
	tests := []struct {
		b    *BoundingBox
		want []Vector
	}{
		{&b1, []Vector{{x: -1, y: 0}, {x: 0, y: -1}, {x: 0, y: 1}, {x: 1, y: 0}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v normals", test.b), func(t *testing.T) {
			got := test.b.normals()
			if len(got) != len(test.want) {
				t.Errorf("want %d normals but got %d", len(test.want), len(got))
			}

			slices.SortFunc(got, func(a *Vector, b *Vector) int {
				if a.x == b.x {
					return int(math.Copysign(1, a.y-b.y))
				} else {
					return int(math.Copysign(1, a.x-b.x))
				}
			})
			for i, want := range test.want {
				if math.Abs(got[i].x-want.x) > epsilon ||
					math.Abs(got[i].y-want.y) > epsilon {
					t.Errorf("want %v but got %v", test.want, got)
				}
			}
		})
	}
}

func TestConvertToWorldSpace(t *testing.T) {
	tests := []struct {
		b    *BoundingBox
		v    Vector
		want Vector
	}{
		{&b2, Vector{x: -1, y: -1}, Vector{x: 1, y: 2 - math.Sqrt2}},
		{&b2, Vector{x: 1, y: -1}, Vector{x: 1 + math.Sqrt2, y: 2}},
		{&b2, Vector{x: 1, y: 1}, Vector{x: 1, y: 2 + math.Sqrt2}},
		{&b2, Vector{x: -1, y: 1}, Vector{x: 1 - math.Sqrt2, y: 2}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v wordlify %v", test.b, test.v), func(t *testing.T) {
			got := test.b.convertToWorldSpace(&test.v)
			if math.Abs(got.x-test.want.x) > epsilon ||
				math.Abs(got.y-test.want.y) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestProjectionRange(t *testing.T) {
	tests := []struct {
		b       *BoundingBox
		v       Vector
		wantMin float64
		wantMax float64
	}{
		{&b3, Vector{x: 0, y: 1}, -math.Sqrt2, math.Sqrt2},
		{&b3, Vector{x: 1, y: 0}, -math.Sqrt2, math.Sqrt2},
		{&b3, Vector{x: 1, y: 1}, -1, 1},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v projectionRange on %v", test.b, test.v), func(t *testing.T) {
			gotMin, gotMax := test.b.projectionRange(&test.v)
			if math.Abs(gotMin-test.wantMin) > epsilon ||
				math.Abs(gotMax-test.wantMax) > epsilon {
				t.Errorf("want (%f, %f) but got (%f, %f)", test.wantMin, test.wantMax, gotMin, gotMax)
			}
		})
	}
}

func TestDidCollide(t *testing.T) {
	tests := []struct {
		b1   *BoundingBox
		b2   *BoundingBox
		want bool
	}{
		{&b1, &b1, true},
		{&b1, &b2, true},
		{&b1, &b3, true},
		{&b1, &b4, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v collide with %v", test.b1, test.b2), func(t *testing.T) {
			got := test.b1.DidCollide(test.b2)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}
