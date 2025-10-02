package geometry

import (
	"math"
	"math/rand"
	"slices"
)

const (
	MIN_RADIUS     = 20
	MAX_RADIUS     = 100
	MIN_NUM_POINTS = 4
	MAX_NUM_POINTS = 16
)

func NewRandomConvexHull() []Vector {
	numPoints := MIN_NUM_POINTS + rand.Intn(MAX_NUM_POINTS-MIN_NUM_POINTS+1)
	points := []Vector{}
	for range numPoints {
		points = append(points, Vector{
			X: math.Copysign(MIN_RADIUS+rand.Float64()*MAX_RADIUS, rand.Float64()-0.5),
			Y: math.Copysign(MIN_RADIUS+rand.Float64()*MAX_RADIUS, rand.Float64()-0.5),
		})
	}
	return ConvexHull(points)
}

func ConvexHull(points []Vector) []Vector {
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
		if a.Sub(origin).Angle() < b.Sub(origin).Angle() {
			return -1
		} else {
			return 1
		}
	})
}

func isLeftTurn(a *Vector, b *Vector, c *Vector) bool {
	u := b.Sub(a)
	v := c.Sub(b)
	return math.Atan2(u.X*v.Y-u.Y*v.X, u.X*v.X+u.Y*v.Y) > 0
}
