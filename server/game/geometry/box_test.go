package geometry

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

var square = []*Vector{
	NewVector(-1, -1),
	NewVector(1, -1),
	NewVector(1, 1),
	NewVector(-1, 1),
}

var p1 = NewVector(0, 0)
var r1 = 0.0
var b1 = NewBoundingBox(p1, &r1, &square)

var p2 = NewVector(1, 2)
var r2 = math.Pi / 4
var b2 = NewBoundingBox(p2, &r2, &square)

var p3 = NewVector(0, 0)
var r3 = math.Pi / 4
var b3 = NewBoundingBox(p3, &r3, &square)

var p4 = NewVector(2, 2)
var r4 = math.Pi / 4
var b4 = NewBoundingBox(p4, &r4, &square)

func TestDidCollide(t *testing.T) {
	tests := []struct {
		b1   *BoundingBox
		b2   *BoundingBox
		want bool
	}{
		{b1, b1, true},
		{b1, b2, true},
		{b1, b3, true},
		{b1, b4, false},
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

func TestDidCollideWithMovement(t *testing.T) {
	p5 := NewVector(5, 0)
	r5 := 0.0
	b5 := NewBoundingBox(p5, &r5, &square)

	got := b5.DidCollide(b1)
	want := false
	if got != want {
		t.Error("should not collide before movement")
	}

	p5.X = 1
	got = b5.DidCollide(b1)
	want = true
	if got != want {
		t.Error("should collide after movement")
	}
}

func TestHorizontalBounds(t *testing.T) {
	tests := []struct {
		b    *BoundingBox
		want struct {
			min float64
			max float64
		}
	}{
		{
			b1, struct {
				min float64
				max float64
			}{-1, 1},
		},
		{
			b2, struct {
				min float64
				max float64
			}{1 - math.Sqrt2, 1 + math.Sqrt2},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v horizontal bounds", test.b), func(t *testing.T) {
			gotMin, gotMax := test.b.HorizontalBounds()
			if math.Abs(gotMin-test.want.min) > EPSILON ||
				math.Abs(gotMax-test.want.max) > EPSILON {
				t.Errorf("want %v but got (%v, %v)", test.want, gotMin, gotMax)
			}
		})
	}
}

func TestNormals(t *testing.T) {
	tests := []struct {
		b    *BoundingBox
		want []Vector
	}{
		{
			b1, []Vector{
				{X: -1, Y: 0},
				{X: 0, Y: -1},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
			},
		},
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
				if math.Abs(got[i].X-want.X) > EPSILON ||
					math.Abs(got[i].Y-want.Y) > EPSILON {
					t.Errorf("want %v but got %v", want, got[i])
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
		{b2, Vector{X: -1, Y: -1}, Vector{X: 1, Y: 2 - math.Sqrt2}},
		{b2, Vector{X: 1, Y: -1}, Vector{X: 1 + math.Sqrt2, Y: 2}},
		{b2, Vector{X: 1, Y: 1}, Vector{X: 1, Y: 2 + math.Sqrt2}},
		{b2, Vector{X: -1, Y: 1}, Vector{X: 1 - math.Sqrt2, Y: 2}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v wordlify %v", test.b, test.v), func(t *testing.T) {
			got := test.b.convertToWorldSpace(&test.v)
			if math.Abs(got.X-test.want.X) > EPSILON ||
				math.Abs(got.Y-test.want.Y) > EPSILON {
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
		{b3, Vector{X: 0, Y: 1}, -math.Sqrt2, math.Sqrt2},
		{b3, Vector{X: 1, Y: 0}, -math.Sqrt2, math.Sqrt2},
		{b3, Vector{X: 1, Y: 1}, -1, 1},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v projectionRange on %v", test.b, test.v), func(t *testing.T) {
			gotMin, gotMax := test.b.projectionRange(&test.v)
			if math.Abs(gotMin-test.wantMin) > EPSILON ||
				math.Abs(gotMax-test.wantMax) > EPSILON {
				t.Errorf("want (%f, %f) but got (%f, %f)", test.wantMin, test.wantMax, gotMin, gotMax)
			}
		})
	}
}
