package geometry

import (
	"fmt"
	"math"
	"testing"
)

func TestNewRectangleHull(t *testing.T) {
	tests := map[string]struct {
		width  float64
		height float64
		want   []*Vector
	}{
		"RectangleHull":             {10, 20, []*Vector{{X: -5, Y: -10}, {X: 5, Y: -10}, {X: 5, Y: 10}, {X: -5, Y: 10}}},
		"RectangleHull with square": {10, 10, []*Vector{{X: -5, Y: -5}, {X: 5, Y: -5}, {X: 5, Y: 5}, {X: -5, Y: 5}}},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: rectangle hull with w = %f and h = %f", desc, test.width, test.height)
		t.Run(title, func(t *testing.T) {
			got := NewRectangleHull(test.width, test.height)
			if len(got) != len(test.want) {
				t.Errorf("want hull of length %d but got %v", len(test.want), got)
			}
			for i := range got {
				if !got[i].IsEqual(test.want[i]) {
					t.Errorf("want %v but got %v", test.want[i], got[i])
				}
			}
		})
	}
}

func TestConvexHull(t *testing.T) {
	tests := map[string]struct {
		points []*Vector
		want   []*Vector
	}{
		"ConvexHull": {
			[]*Vector{{X: 1, Y: 1}, {X: 4, Y: 2}, {X: -3, Y: 1}, {X: -2, Y: -2}, {X: 0, Y: 4}},
			[]*Vector{{X: -2, Y: -2}, {X: 4, Y: 2}, {X: 0, Y: 4}, {X: -3, Y: 1}},
		},
		"ConvexHull with extra point": {
			[]*Vector{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 0, Y: 0}, {X: 1, Y: 1}, {X: -1, Y: 1}},
			[]*Vector{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}},
		},
		"ConvexHull with square": {
			[]*Vector{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}},
			[]*Vector{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}},
		},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: convex hull of %v", desc, test.points)
		t.Run(title, func(t *testing.T) {
			got := ConvexHull(test.points)
			if len(got) != len(test.want) {
				t.Errorf("want hull of length %d but got %v", len(test.want), got)
			}
			for i := range got {
				if !got[i].IsEqual(test.want[i]) {
					t.Errorf("want %v but got %v", test.want, got)
				}
			}
		})
	}
}

func TestHullArea(t *testing.T) {
	tests := map[string]struct {
		points []*Vector
		want   float64
	}{
		"HullArea":        {[]*Vector{{X: 2, Y: -3}, {X: 3, Y: -1}, {X: 2, Y: 1}, {X: 0, Y: 1}, {X: -1, Y: -1}}, 10},
		"HullArea square": {[]*Vector{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}, 4},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: convex hull of %v", desc, test.points)
		t.Run(title, func(t *testing.T) {
			got := HullArea(test.points)
			if math.Abs(got-test.want) > EPSILON {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}

func TestSortPointsAbout(t *testing.T) {
	tests := map[string]struct {
		origin Vector
		points []*Vector
		want   []*Vector
	}{
		"sortPointsAbout": {
			Vector{X: -2, Y: -2},
			[]*Vector{{1, 1}, {4, 2}, {-3, 1}, {0, 4}},
			[]*Vector{{4, 2}, {1, 1}, {0, 4}, {-3, 1}},
		},
		"sortPointsAbout origin": {
			Vector{X: 0, Y: 0},
			[]*Vector{{-1, 1}, {-1, -1}, {1, -1}, {1, 1}},
			[]*Vector{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}},
		},
	}

	for desc, test := range tests {
		title := fmt.Sprintf("%s: sort %v about %v", desc, test.points, test.origin)
		t.Run(title, func(t *testing.T) {
			sortPointsAbout(&test.origin, test.points)

			got := test.points
			if len(got) != len(test.want) {
				t.Errorf("want sorted points of length %d but got %v", len(test.want), got)
			}
			for i := range got {
				if !got[i].IsEqual(test.want[i]) {
					t.Errorf("want %v but got %v", test.want[i], got[i])
				}
			}
		})
	}
}
