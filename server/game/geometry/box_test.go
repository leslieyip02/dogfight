package geometry

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

var b1 = BoundingBox{
	origin: Vector{X: 0, Y: 0},
	theta:  0,
	points: &[]Vector{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	},
}

var b2 = BoundingBox{
	origin: Vector{X: 1, Y: 2},
	theta:  math.Pi / 4,
	points: &[]Vector{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	},
}

var b3 = BoundingBox{
	origin: Vector{X: 0, Y: 0},
	theta:  math.Pi / 4,
	points: &[]Vector{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	},
}

var b4 = BoundingBox{
	origin: Vector{X: 2, Y: 2},
	theta:  math.Pi / 4,
	points: &[]Vector{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	},
}

func TestNormals(t *testing.T) {
	tests := []struct {
		b    *BoundingBox
		want []Vector
	}{
		{&b1, []Vector{{X: -1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}, {X: 1, Y: 0}}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v normals", test.b), func(t *testing.T) {
			got := test.b.normals()
			if len(got) != len(test.want) {
				t.Errorf("want %d normals but got %d", len(test.want), len(got))
			}

			slices.SortFunc(got, func(a *Vector, b *Vector) int {
				if a.X == b.X {
					return int(math.Copysign(1, a.Y-b.Y))
				} else {
					return int(math.Copysign(1, a.X-b.X))
				}
			})
			for i, want := range test.want {
				if math.Abs(got[i].X-want.X) > epsilon ||
					math.Abs(got[i].Y-want.Y) > epsilon {
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
		{&b2, Vector{X: -1, Y: -1}, Vector{X: 1, Y: 2 - math.Sqrt2}},
		{&b2, Vector{X: 1, Y: -1}, Vector{X: 1 + math.Sqrt2, Y: 2}},
		{&b2, Vector{X: 1, Y: 1}, Vector{X: 1, Y: 2 + math.Sqrt2}},
		{&b2, Vector{X: -1, Y: 1}, Vector{X: 1 - math.Sqrt2, Y: 2}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v wordlify %v", test.b, test.v), func(t *testing.T) {
			got := test.b.convertToWorldSpace(&test.v)
			if math.Abs(got.X-test.want.X) > epsilon ||
				math.Abs(got.Y-test.want.Y) > epsilon {
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
		{&b3, Vector{X: 0, Y: 1}, -math.Sqrt2, math.Sqrt2},
		{&b3, Vector{X: 1, Y: 0}, -math.Sqrt2, math.Sqrt2},
		{&b3, Vector{X: 1, Y: 1}, -1, 1},
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
