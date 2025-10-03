package geometry

import (
	"fmt"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		u    Vector
		v    Vector
		want Vector
	}{
		{Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}},
		{Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, Vector{X: 2, Y: 4}},
		{Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, Vector{X: 4, Y: 6}},
		{Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, Vector{X: -2, Y: -6}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v + %v", test.u, test.v), func(t *testing.T) {
			got := test.u.Add(&test.v)
			if math.Abs(got.X-test.want.X) > epsilon ||
				math.Abs(got.Y-test.want.Y) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		u    Vector
		v    Vector
		want Vector
	}{
		{Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}},
		{Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, Vector{X: 0, Y: 0}},
		{Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, Vector{X: -2, Y: -2}},
		{Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, Vector{X: 4, Y: 2}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v - %v", test.u, test.v), func(t *testing.T) {
			got := test.u.Sub(&test.v)
			if math.Abs(got.X-test.want.X) > epsilon ||
				math.Abs(got.Y-test.want.Y) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		u    Vector
		s    float64
		want Vector
	}{
		{Vector{X: 0, Y: 0}, 1, Vector{X: 0, Y: 0}},
		{Vector{X: 1, Y: 2}, 1, Vector{X: 1, Y: 2}},
		{Vector{X: 1, Y: 2}, 2, Vector{X: 2, Y: 4}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%f * %v", test.s, test.u), func(t *testing.T) {
			got := test.u.Multiply(test.s)
			if math.Abs(got.X-test.want.X) > epsilon ||
				math.Abs(got.Y-test.want.Y) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestDot(t *testing.T) {
	tests := []struct {
		u    Vector
		v    Vector
		want float64
	}{
		{Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, 0},
		{Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, 5},
		{Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, 11},
		{Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, 5},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v . %v", test.u, test.v), func(t *testing.T) {
			got := test.u.dot(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestCross(t *testing.T) {
	tests := []struct {
		u    Vector
		v    Vector
		want float64
	}{
		{Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, 0},
		{Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, 0},
		{Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, -2},
		{Vector{X: 1, Y: -2}, Vector{X: -3, Y: -4}, -10},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v x %v", test.u, test.v), func(t *testing.T) {
			got := test.u.cross(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestLength(t *testing.T) {
	tests := []struct {
		u    Vector
		want float64
	}{
		{Vector{X: 0, Y: 0}, 0},
		{Vector{X: 1, Y: 2}, math.Sqrt(5)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("length of %v", test.u), func(t *testing.T) {
			got := test.u.Length()
			if math.Abs(got-test.want) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestGradient(t *testing.T) {
	tests := []struct {
		u    Vector
		want float64
	}{
		{Vector{X: 1, Y: 2}, 2},
		{Vector{X: 2, Y: 1}, 0.5},
		{Vector{X: 1, Y: -2}, -2},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("gradient of %v", test.u), func(t *testing.T) {
			got := test.u.gradient()
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestGradientNaN(t *testing.T) {
	u := Vector{X: 0, Y: 0}
	got := u.gradient()
	if !math.IsNaN(got) {
		t.Errorf("want NaN but got %v", got)
	}
}

func TestAngle(t *testing.T) {
	tests := []struct {
		u    Vector
		want float64
	}{
		{Vector{X: 1, Y: 1}, math.Pi / 4},
		{Vector{X: 1, Y: 2}, 1.10714},
		{Vector{X: 1, Y: -2}, -1.10714},
		{Vector{X: 0, Y: 0}, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("angle of %v", test.u), func(t *testing.T) {
			got := test.u.Angle()
			if math.Abs(got-test.want) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		u     Vector
		theta float64
		want  Vector
	}{
		{Vector{X: 1, Y: 0}, math.Pi / 2, Vector{X: 0, Y: 1}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("rotate %v", test.u), func(t *testing.T) {
			got := test.u.rotate(test.theta)
			if math.Abs(got.X-test.want.X) > epsilon ||
				math.Abs(got.Y-test.want.Y) > epsilon {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestNormal(t *testing.T) {
	tests := []struct {
		u    Vector
		want Vector
	}{
		{Vector{X: 0, Y: 1}, Vector{X: 1, Y: 0}},
		{Vector{X: 1, Y: 1}, Vector{X: -math.Sqrt2 / 2, Y: math.Sqrt2 / 2}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("normalize %v", test.u), func(t *testing.T) {
			got := test.u.normal()
			if math.Abs(got.Length()-1) > epsilon {
				t.Errorf("want vector of length 1 but got %v", got)
			}
			if !got.isParallel(&test.want) {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestIsParallel(t *testing.T) {
	tests := []struct {
		u    Vector
		v    Vector
		want bool
	}{
		{Vector{X: 1, Y: 2}, Vector{X: 1, Y: 2}, true},
		{Vector{X: 1, Y: 2}, Vector{X: 2, Y: 4}, true},
		{Vector{X: 1, Y: 2}, Vector{X: -1, Y: -2}, true},
		{Vector{X: 0, Y: 0}, Vector{X: 0, Y: 0}, true},
		{Vector{X: 1, Y: 0}, Vector{X: 0, Y: 1}, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v isParallel to %v", test.u, test.v), func(t *testing.T) {
			got := test.u.isParallel(&test.v)
			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}
