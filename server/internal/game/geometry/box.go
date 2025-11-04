package geometry

import (
	"math"
)

// A BoundingBox is a group of closed, convex points in 2D space.
type BoundingBox struct {
	position *Vector    // origin
	rotation *float64   // rotation of the points about position
	points   *[]*Vector // points stored in anticlockwise order
}

func NewBoundingBox(position *Vector, rotation *float64, points *[]*Vector) *BoundingBox {
	return &BoundingBox{
		position: position,
		rotation: rotation,
		points:   points,
	}
}

// DidCollide uses the Separating Axis Theorem (SAT) to determine if b1 is
// colliding with b2.
func (b1 *BoundingBox) DidCollide(b2 *BoundingBox) bool {
	// Get perpendicular vectors for each line segment.
	p := make(map[float64]bool)
	for _, normal := range b1.normals() {
		r := normal.Rotate(*b1.rotation)

		// Deduplicate parallel lines by using their gradients.
		m := r.gradient()
		p[m] = true
	}
	for _, normal := range b2.normals() {
		r := normal.Rotate(*b2.rotation)
		m := r.gradient()
		p[m] = true
	}

	// Calculate the projection range of each perpendicular vector and check
	// for overlaps.
	for m := range p {
		perpendicular := Vector{
			X: 1,
			Y: m,
		}

		min1, max1 := b1.projectionRange(&perpendicular)
		min2, max2 := b2.projectionRange(&perpendicular)
		if min1 > max2 || min2 > max1 {
			return false
		}
	}

	return true
}

// HorizontalBounds returns a pair of x-coordinates (in world space) which
// bounds b.
func (b *BoundingBox) HorizontalBounds() (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	for _, point := range *b.points {
		w := b.convertToWorldSpace(point)
		min = math.Min(w.X, min)
		max = math.Max(w.X, max)
	}
	return min, max
}

// normals returns the normal vectors of each line segment in b.
func (b *BoundingBox) normals() []*Vector {
	normals := []*Vector{}
	for i := range len(*b.points) {
		u := (*b.points)[i]
		v := (*b.points)[(i+1)%len(*b.points)]
		normals = append(normals, (v.Sub(u)).Normal())
	}
	return normals
}

// convertToWorldSpace converts a vector to world space with respect to b.
func (b *BoundingBox) convertToWorldSpace(v *Vector) *Vector {
	// transform to world space
	u := v.Rotate(*b.rotation)
	return &Vector{
		X: u.X + b.position.X,
		Y: u.Y + b.position.Y,
	}
}

// projectionRange returns a pair of coordinates which represent the bounds of
// the projection of b onto v. Each point in b is projected on to v, and the
// minimimum and maximum values along v are returned.
func (b *BoundingBox) projectionRange(v *Vector) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	for _, point := range *b.points {
		w := b.convertToWorldSpace(point)

		// scalar projection
		s := w.dot(v) / v.Length()
		min = math.Min(s, min)
		max = math.Max(s, max)
	}
	return min, max
}
