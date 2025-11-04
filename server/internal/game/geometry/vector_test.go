package geometry

import (
	"fmt"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		v    Vector
		want Vector
	}{
		"Add":                      {Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, Vector{X: 4, Y: 6}},
		"Add with negative values": {Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, Vector{X: -2, Y: -6}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v + %v", desc, test.u, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.u.Add(&test.v)
			if !got.IsEqual(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		v    Vector
		want Vector
	}{
		"Sub":                      {Vector{X: 3, Y: 4}, Vector{X: 1, Y: 2}, Vector{X: 2, Y: 2}},
		"Sub with negative values": {Vector{X: 1, Y: -2}, Vector{X: -3, Y: 4}, Vector{X: 4, Y: -6}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v - %v", desc, test.u, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.u.Sub(&test.v)
			if !got.IsEqual(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		s    float64
		want Vector
	}{
		"Multiply":                     {Vector{X: 1, Y: 2}, 2, Vector{X: 2, Y: 4}},
		"Multiply with negative value": {Vector{X: 1, Y: 2}, -2, Vector{X: -2, Y: -4}},
		"Multiply with zero":           {Vector{X: 1, Y: 2}, 0, Vector{X: 0, Y: 0}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %f * %v", desc, test.s, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.Multiply(test.s)
			if !got.IsEqual(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestLength(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		want float64
	}{
		"Length":                {Vector{X: 1, Y: 2}, math.Sqrt(5)},
		"Length of zero vector": {Vector{X: 0, Y: 0}, 0},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: ||%v||", desc, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.Length()
			if math.Abs(got-test.want) > EPSILON {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestUnit(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		want Vector
	}{
		"Unit":                {Vector{X: 1, Y: 1}, Vector{X: math.Sqrt2 / 2, Y: math.Sqrt2 / 2}},
		"Unit of zero vector": {Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: unit vector of %v", desc, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.Unit()
			if !got.IsEqual(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestAngle(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		want float64
	}{
		"Angle":                     {Vector{X: 1, Y: 2}, 1.10714},
		"Angle with negative value": {Vector{X: 1, Y: -2}, -1.10714},
		"Angle of zero vector":      {Vector{X: 0, Y: 0}, 0},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: angle of %v", desc, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.Angle()
			if math.Abs(got-test.want) > EPSILON {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	tests := map[string]struct {
		u     Vector
		theta float64
		want  Vector
	}{
		"Rotate":                     {Vector{X: 1, Y: 0}, math.Pi / 2, Vector{X: 0, Y: 1}},
		"Rotate with negative value": {Vector{X: 1, Y: 0}, -math.Pi / 2, Vector{X: 0, Y: -1}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: rotate %v", desc, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.Rotate(test.theta)
			if !got.IsEqual(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestNormal(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		want Vector
	}{
		"Normal":                 {Vector{X: 1, Y: 1}, Vector{X: -math.Sqrt2 / 2, Y: math.Sqrt2 / 2}},
		"Normal with zero value": {Vector{X: 0, Y: 1}, Vector{X: 1, Y: 0}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: normal of %v", desc, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.Normal()
			if math.Abs(got.Length()-1) > EPSILON {
				t.Errorf("expected vector of length 1 but got %v", got)
			}
			if !got.isParallel(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		v    Vector
		want bool
	}{
		"IsEqual":                            {Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, true},
		"IsEqual with acceptable difference": {Vector{X: 1, Y: 2}, Vector{X: 1 + 1e-8, Y: 2}, true},
		"IsEqual negative example":           {Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, false},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v = %v", desc, test.u, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.u.IsEqual(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestDot(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		v    Vector
		want float64
	}{
		"dot":                      {Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, 11},
		"dot with negative values": {Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, 5},
		"dot with zero vector":     {Vector{X: 1, Y: 2}, Vector{X: 0, Y: 0}, 0},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v . %v", desc, test.u, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.u.dot(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestCross(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		v    Vector
		want float64
	}{
		"cross":                      {Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, -2},
		"cross with negative values": {Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, -10},
		"cross with zero vector":     {Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, 0},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v x %v", desc, test.u, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.u.cross(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestGradient(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		want float64
	}{
		"gradient":                      {Vector{X: 1, Y: 2}, 2},
		"gradient with negative values": {Vector{X: 1, Y: -2}, -2},
		"gradient with zero vector":     {Vector{X: 0, Y: 0}, math.NaN()},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: gradient of %v", desc, test.u)
		t.Run(title, func(t *testing.T) {
			got := test.u.gradient()

			if math.IsNaN(test.want) {
				if !math.IsNaN(got) {
					t.Errorf("want %v but got %v", test.want, got)
				}
				return
			}

			if math.Abs(got-test.want) > EPSILON {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestIsParallel(t *testing.T) {
	tests := map[string]struct {
		u    Vector
		v    Vector
		want bool
	}{
		"isParallel":                         {Vector{X: 1, Y: 2}, Vector{X: 2, Y: 4}, true},
		"isParallel with self":               {Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, true},
		"isParallel with zero vector":        {Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, true},
		"isParallel with opposite direction": {Vector{X: 1, Y: 2}, Vector{X: -1, Y: -2}, true},
		"isParallel negative example":        {Vector{X: 1, Y: 0}, Vector{X: 0, Y: 1}, false},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: %v isParallel to %v", desc, test.u, test.v)
		t.Run(title, func(t *testing.T) {
			got := test.u.isParallel(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}
