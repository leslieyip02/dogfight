package geometry

import (
	"fmt"
	"math"
	"testing"
)

func TestSub(t *testing.T) {
	tests := []struct {
		u    Vector
		v    Vector
		want Vector
	}{
		{Vector{x: 0, y: 0}, Vector{x: 0, y: 0}, Vector{x: 0, y: 0}},
		{Vector{x: 1, y: 2}, Vector{x: 1, y: 2}, Vector{x: 0, y: 0}},
		{Vector{x: 1, y: 2}, Vector{x: 3, y: 4}, Vector{x: -2, y: -2}},
		{Vector{x: 1, y: -2}, Vector{x: -3, y: -4}, Vector{x: 4, y: 2}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v - %v", test.u, test.v), func(t *testing.T) {
			got := test.u.sub(&test.v)
			if math.Abs(got.x-test.want.x) > epsilon ||
				math.Abs(got.y-test.want.y) > epsilon {
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
		{Vector{x: 0, y: 0}, 1, Vector{x: 0, y: 0}},
		{Vector{x: 1, y: 2}, 1, Vector{x: 1, y: 2}},
		{Vector{x: 1, y: 2}, 2, Vector{x: 2, y: 4}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%f * %v", test.s, test.u), func(t *testing.T) {
			got := test.u.multiply(test.s)
			if math.Abs(got.x-test.want.x) > epsilon ||
				math.Abs(got.y-test.want.y) > epsilon {
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
		{Vector{x: 0, y: 0}, Vector{x: 0, y: 0}, 0},
		{Vector{x: 1, y: 2}, Vector{x: 1, y: 2}, 5},
		{Vector{x: 1, y: 2}, Vector{x: 3, y: 4}, 11},
		{Vector{x: 1, y: -2}, Vector{x: -3, y: -4}, 5},
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

func TestLength(t *testing.T) {
	tests := []struct {
		u    Vector
		want float64
	}{
		{Vector{x: 0, y: 0}, 0},
		{Vector{x: 1, y: 2}, math.Sqrt(5)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("length of %v", test.u), func(t *testing.T) {
			got := test.u.length()
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
		{Vector{x: 1, y: 2}, 2},
		{Vector{x: 2, y: 1}, 0.5},
		{Vector{x: 1, y: -2}, -2},
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
	u := Vector{x: 0, y: 0}
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
		{Vector{x: 1, y: 1}, math.Pi / 4},
		{Vector{x: 1, y: 2}, 1.10714},
		{Vector{x: 1, y: -2}, -1.10714},
		{Vector{x: 0, y: 0}, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("angle of %v", test.u), func(t *testing.T) {
			got := test.u.angle()
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
		{Vector{x: 1, y: 0}, math.Pi / 2, Vector{x: 0, y: 1}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("rotate %v", test.u), func(t *testing.T) {
			got := test.u.rotate(test.theta)
			if math.Abs(got.x-test.want.x) > epsilon ||
				math.Abs(got.y-test.want.y) > epsilon {
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
		{Vector{x: 0, y: 1}, Vector{x: 1, y: 0}},
		{Vector{x: 1, y: 1}, Vector{x: -math.Sqrt2 / 2, y: math.Sqrt2 / 2}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("normalize %v", test.u), func(t *testing.T) {
			got := test.u.normal()
			if math.Abs(got.length()-1) > epsilon {
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
		{Vector{x: 1, y: 2}, Vector{x: 1, y: 2}, true},
		{Vector{x: 1, y: 2}, Vector{x: 2, y: 4}, true},
		{Vector{x: 1, y: 2}, Vector{x: -1, y: -2}, true},
		{Vector{x: 0, y: 0}, Vector{x: 0, y: 0}, true},
		{Vector{x: 1, y: 0}, Vector{x: 0, y: 1}, false},
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
