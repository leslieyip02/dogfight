package geometry

import (
	"math"
	"math/rand"
	"slices"
)

const (
	POINT_RANGE_WIDTH  = 300
	POINT_RANGE_HEIGHT = 300
	MIN_NUM_POINTS     = 3
	MAX_NUM_POINTS     = 16
)

func NewRandomConvexHull() []Vector {
	numPoints := MIN_NUM_POINTS + rand.Intn(MAX_NUM_POINTS-MIN_NUM_POINTS+1)
	points := []Vector{}
	for range numPoints {
		points = append(points, Vector{
			X: rand.Float64()*POINT_RANGE_WIDTH - POINT_RANGE_WIDTH/2,
			Y: rand.Float64()*POINT_RANGE_HEIGHT - POINT_RANGE_HEIGHT/2,
		})
	}
	return convexHull(points)
}

func convexHull(points []Vector) []Vector {
	// graham scan

	// start at the point with the smallest y (ties broken by x)
	origin := points[0]
	for i := 1; i < len(points); i++ {
		if points[i].Y < origin.Y || (points[i].Y == origin.Y && points[i].X < origin.X) {
			origin = points[i]
		}
	}

	// sort remaining points by angle to origin
	remainingPoints := []Vector{}
	for _, point := range points {
		if point == origin {
			continue
		}
		remainingPoints = append(remainingPoints, point)
	}
	sortPoints(&origin, remainingPoints)

	hull := []Vector{origin}
	for _, point := range remainingPoints {
		for len(hull) > 1 && !isLeftTurn(&hull[len(hull)-2], &hull[len(hull)-1], &point) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, point)
	}
	return hull
}

func sortPoints(origin *Vector, points []Vector) {
	slices.SortFunc(points, func(a Vector, b Vector) int {
		// return int(math.Copysign(a.sub(origin).angle()-b.sub(origin).angle(), 1))
		if a.sub(origin).angle() < b.sub(origin).angle() {
			return -1
		} else {
			return 1
		}
	})
}

func isLeftTurn(a *Vector, b *Vector, c *Vector) bool {
	u := b.sub(a)
	v := c.sub(b)
	return math.Atan2(u.X*v.Y-u.Y*v.X, u.X*v.X+u.Y*v.Y) > 0
}
