package geometry

import (
	"math"
)

type BoundingBox struct {
	position *Vector
	rotation *float64

	// points take position as (0, 0)
	// points are stored in anticlockwise order
	points *[]*Vector
}

func NewBoundingBox(position *Vector, rotation *float64, points *[]*Vector) *BoundingBox {
	return &BoundingBox{
		position: position,
		rotation: rotation,
		points:   points,
	}
}

func (b1 *BoundingBox) DidCollide(b2 *BoundingBox) bool {
	// separatable axis theorem

	// for each line segment,
	// get a perpendicular vector
	p := make(map[float64]bool)
	for _, normal := range b1.normals() {
		r := normal.Rotate(*b1.rotation)
		// take gradients to deduplicate parallel lines
		m := r.gradient()
		p[m] = true
	}
	for _, normal := range b2.normals() {
		r := normal.Rotate(*b2.rotation)
		m := r.gradient()
		p[m] = true
	}

	// for each perpendicular vector,
	// calculate the projection
	for m := range p {
		perpendicular := Vector{
			X: 1,
			Y: m,
		}

		// check overlaps
		min1, max1 := b1.projectionRange(&perpendicular)
		min2, max2 := b2.projectionRange(&perpendicular)
		if min1 > max2 || min2 > max1 {
			return false
		}
	}

	return true
}

func (b *BoundingBox) normals() []*Vector {
	normals := []*Vector{}
	for i := range len(*b.points) {
		u := (*b.points)[i]
		v := (*b.points)[(i+1)%len(*b.points)]
		normals = append(normals, (v.Sub(u)).Normal())
	}
	return normals
}

func (b *BoundingBox) convertToWorldSpace(v *Vector) *Vector {
	// transform to world space
	u := v.Rotate(*b.rotation)
	return &Vector{
		X: u.X + b.position.X,
		Y: u.Y + b.position.Y,
	}
}

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
