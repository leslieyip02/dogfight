package geometry

import (
	"fmt"
	"testing"
)

func TestSortPoints(t *testing.T) {
	tests := []struct {
		origin Vector
		points []Vector
		want   []Vector
	}{
		{
			origin: Vector{X: -2, Y: -2},
			points: []Vector{
				{X: 1, Y: 1},
				{X: 4, Y: 2},
				{X: -3, Y: 1},
				{X: 0, Y: 4},
			},
			want: []Vector{
				{X: 4, Y: 2},
				{X: 1, Y: 1},
				{X: 0, Y: 4},
				{X: -3, Y: 1},
			},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("sort %v w.r.t. %v", test.points, test.origin), func(t *testing.T) {
			sortPoints(&test.origin, test.points)
			got := test.points
			if len(got) != len(test.want) {
				t.Errorf("want sorted points of length %d but got %v", len(test.want), got)
			}
			for i, want := range test.want {
				if got[i].X != want.X || got[i].Y != want.Y {
					t.Errorf("want %v but got %v", test.want, got)
				}
			}
		})
	}
}

func TestConvexHull(t *testing.T) {
	tests := []struct {
		points []Vector
		want   []Vector
	}{
		{
			points: []Vector{
				{X: -1, Y: -1},
				{X: 1, Y: -1},
				{X: 1, Y: 1},
				{X: -1, Y: 1},
			},
			want: []Vector{
				{X: -1, Y: -1},
				{X: 1, Y: -1},
				{X: 1, Y: 1},
				{X: -1, Y: 1},
			},
		},
		{
			points: []Vector{
				{X: -1, Y: -1},
				{X: 1, Y: -1},
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: -1, Y: 1},
			},
			want: []Vector{
				{X: -1, Y: -1},
				{X: 1, Y: -1},
				{X: 1, Y: 1},
				{X: -1, Y: 1},
			},
		},
		{
			points: []Vector{
				{X: 1, Y: 1},
				{X: 4, Y: 2},
				{X: -3, Y: 1},
				{X: -2, Y: -2},
				{X: 0, Y: 4},
			},
			want: []Vector{
				{X: -2, Y: -2},
				{X: 4, Y: 2},
				{X: 0, Y: 4},
				{X: -3, Y: 1},
			},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("convex hull of %v", test.points), func(t *testing.T) {
			got := ConvexHull(test.points)
			if len(got) != len(test.want) {
				t.Errorf("want hull of length %d but got %v", len(test.want), got)
			}
			for i, want := range test.want {
				if got[i].X != want.X || got[i].Y != want.Y {
					t.Errorf("want %v but got %v", test.want, got)
				}
			}
		})
	}
}
