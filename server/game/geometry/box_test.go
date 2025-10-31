package geometry

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

func constructBoundingBox(x float64, y float64, rotation float64, points []*Vector) *BoundingBox {
	position := NewVector(x, y)
	return NewBoundingBox(position, &rotation, &points)
}

var b1 = constructBoundingBox(0, 0, 0, square)
var b2 = constructBoundingBox(1, 2, 0, square)
var b3 = constructBoundingBox(0, 0, math.Pi/4, square)
var b4 = constructBoundingBox(1, 2, math.Pi/4, square)
var b5 = constructBoundingBox(2, 2, math.Pi/4, square)

func TestDidCollide(t *testing.T) {
	tests := map[string]struct {
		b1   *BoundingBox
		b2   *BoundingBox
		want bool
	}{
		"DidCollide":                  {b1, b4, true},
		"DidCollide with self":        {b1, b1, true},
		"DidCollide negative example": {b1, b5, false},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v collide with %v", desc, test.b1, test.b2)
		t.Run(title, func(t *testing.T) {
			got := test.b1.DidCollide(test.b2)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestHorizontalBounds(t *testing.T) {
	tests := map[string]struct {
		b       *BoundingBox
		wantMin float64
		wantMax float64
	}{
		"HorizontalBounds":                     {b1, -1, 1},
		"HorizontalBounds with rotated square": {b4, 1 - math.Sqrt2, 1 + math.Sqrt2},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v horizontal bounds", desc, test.b)
		t.Run(title, func(t *testing.T) {
			gotMin, gotMax := test.b.HorizontalBounds()
			if math.Abs(gotMin-test.wantMin) > EPSILON || math.Abs(gotMax-test.wantMax) > EPSILON {
				t.Errorf("want (%f, %f) but got (%f, %f)", test.wantMin, test.wantMax, gotMin, gotMax)
			}
		})
	}
}

func TestNormals(t *testing.T) {
	tests := map[string]struct {
		b    *BoundingBox
		want []Vector
	}{
		"normals": {b1, []Vector{{X: -1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}, {X: 1, Y: 0}}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: normals of %v", desc, test.b)
		t.Run(title, func(t *testing.T) {
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
			for i := range got {
				if !got[i].IsEqual(&test.want[i]) {
					t.Errorf("want %v but got %v", test.want[i], got[i])
				}
			}
		})
	}
}

func TestConvertToWorldSpace(t *testing.T) {
	tests := map[string]struct {
		b    *BoundingBox
		v    Vector
		want Vector
	}{
		"ConvertToWorldSpace":                               {b1, Vector{X: 1, Y: 1}, Vector{X: 1, Y: 1}},
		"ConvertToWorldSpace with translation":              {b2, Vector{X: 1, Y: 1}, Vector{X: 2, Y: 3}},
		"ConvertToWorldSpace with rotation":                 {b3, Vector{X: 1, Y: 1}, Vector{X: 0, Y: math.Sqrt2}},
		"ConvertToWorldSpace with translation and rotation": {b4, Vector{X: 1, Y: 1}, Vector{X: 1, Y: 2 + math.Sqrt2}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: convert %v to world space with respect to %v", desc, test.v, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.b.convertToWorldSpace(&test.v)
			if !got.IsEqual(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestProjectionRange(t *testing.T) {
	tests := map[string]struct {
		b       *BoundingBox
		v       Vector
		wantMin float64
		wantMax float64
	}{
		"projectionRange":                     {b3, Vector{X: 1, Y: 1}, -1, 1},
		"projectionRange with negative value": {b3, Vector{X: -1, Y: 0}, -math.Sqrt2, math.Sqrt2},
		"projectionRange with zero value":     {b3, Vector{X: 0, Y: 1}, -math.Sqrt2, math.Sqrt2},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v projectionRange on %v", desc, test.b, test.v)
		t.Run(title, func(t *testing.T) {
			gotMin, gotMax := test.b.projectionRange(&test.v)
			if math.Abs(gotMin-test.wantMin) > EPSILON ||
				math.Abs(gotMax-test.wantMax) > EPSILON {
				t.Errorf("want (%f, %f) but got (%f, %f)", test.wantMin, test.wantMax, gotMin, gotMax)
			}
		})
	}
}
