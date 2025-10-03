package geometry

import (
	"math"
	"math/rand"
	"slices"
)

func NewRandomConvexHull(
	minNumPoints int,
	maxNumPoints int,
	minRadius float64,
	maxRadius float64,
) []*Vector {
	numPoints := minNumPoints + rand.Intn(maxNumPoints-minNumPoints+1)
	points := []*Vector{}
	for range numPoints {
		points = append(points, &Vector{
			X: math.Copysign(minRadius+rand.Float64()*maxRadius, rand.Float64()-0.5),
			Y: math.Copysign(minRadius+rand.Float64()*maxRadius, rand.Float64()-0.5),
		})
	}
	return ConvexHull(points)
}

func ConvexHull(points []*Vector) []*Vector {
	// graham scan

	// start at the point with the smallest y (ties broken by x)
	origin := points[0]
	for i := 1; i < len(points); i++ {
		if points[i].Y < origin.Y || (points[i].Y == origin.Y && points[i].X < origin.X) {
			origin = points[i]
		}
	}

	// sort remaining points by angle to origin
	remainingPoints := []*Vector{}
	for _, point := range points {
		if point == origin {
			continue
		}
		remainingPoints = append(remainingPoints, point)
	}
	sortPointsAround(origin, remainingPoints)

	hull := []*Vector{origin}
	for _, point := range remainingPoints {
		for len(hull) > 1 && !isLeftTurn(hull[len(hull)-2], hull[len(hull)-1], point) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, point)
	}
	return hull
}

func HullArea(points []*Vector) float64 {
	area := 0.0
	for i := range len(points) {
		j := (i + 1) % len(points)
		area += points[i].X*points[j].Y - points[j].X*points[i].Y
	}
	return math.Abs(area) / 2.0
}

func sortPointsAround(origin *Vector, points []*Vector) {
	slices.SortFunc(points, func(a *Vector, b *Vector) int {
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
