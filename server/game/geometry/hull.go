package geometry

import (
	"math"
	"math/rand"
	"slices"
)

var square = []*Vector{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}}

func NewRectangleHull(
	width float64,
	height float64,
) []*Vector {
	points := make([]*Vector, 4)
	for i, corner := range square {
		points[i] = NewVector(corner.X*(width/2), corner.Y*(height/2))
	}
	return points
}

func NewRandomConvexHull(
	minNumPoints int,
	maxNumPoints int,
	minRadius float64,
	maxRadius float64,
) []*Vector {
	numPoints := minNumPoints + rand.Intn(maxNumPoints-minNumPoints+1)
	points := make([]*Vector, numPoints)
	for i := range numPoints {
		points[i] = &Vector{
			X: math.Copysign(
				minRadius+rand.Float64()*maxRadius,
				rand.Float64()-0.5,
			),
			Y: math.Copysign(
				minRadius+rand.Float64()*maxRadius,
				rand.Float64()-0.5,
			),
		}
	}
	return ConvexHull(points)
}

// ConvexHull uses Graham scan to choose a subset of points to form a convex
// hull.
func ConvexHull(points []*Vector) []*Vector {
	// Find the bottom-left corner.
	origin := points[0]
	for i := 1; i < len(points); i++ {
		if points[i].Y < origin.Y ||
			(points[i].Y == origin.Y && points[i].X < origin.X) {
			origin = points[i]
		}
	}

	// Sort remaining points by angle to origin.
	remainingPoints := []*Vector{}
	for _, point := range points {
		if point == origin {
			continue
		}
		remainingPoints = append(remainingPoints, point)
	}
	sortPointsAbout(origin, remainingPoints)

	// Keep adding points to hull as long as they are left turns.
	hull := []*Vector{origin}
	for _, point := range remainingPoints {
		for len(hull) > 1 &&
			!isLeftTurn(hull[len(hull)-2], hull[len(hull)-1], point) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, point)
	}
	return hull
}

// HullArea gets the area of the polygon formed by points using the shoelace
// method.
func HullArea(points []*Vector) float64 {
	area := 0.0
	for i := range len(points) {
		j := (i + 1) % len(points)
		area += points[i].X*points[j].Y - points[j].X*points[i].Y
	}
	return math.Abs(area) / 2.0
}

// sortPointsAbout orders points based on their angle to origin.
func sortPointsAbout(origin *Vector, points []*Vector) {
	slices.SortFunc(points, func(a *Vector, b *Vector) int {
		if a.Sub(origin).Angle() < b.Sub(origin).Angle() {
			return -1
		} else {
			return 1
		}
	})
}

// isLeftTurn reports whether vector bc is a left turn from vector ac.
func isLeftTurn(a *Vector, b *Vector, c *Vector) bool {
	u := b.Sub(a)
	v := c.Sub(b)
	return math.Atan2(u.X*v.Y-u.Y*v.X, u.X*v.X+u.Y*v.Y) > 0
}
