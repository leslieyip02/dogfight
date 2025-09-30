package geometry

import (
	"math"
)

type BoundingBox struct {
	origin Vector
	theta  float64

	// points take position as (0, 0)
	// points are stored in anticlockwise order
	points *[]Vector
}

func NewBoundingBox(points *[]Vector) BoundingBox {
	return BoundingBox{
		origin: Vector{x: 0, y: 0},
		theta:  0,
		points: points,
	}
}

func (b *BoundingBox) Transform(x float64, y float64, theta float64) *BoundingBox {
	return &BoundingBox{
		origin: Vector{x: x, y: y},
		theta:  theta,
		points: b.points,
	}
}

func (b *BoundingBox) normals() []*Vector {
	normals := []*Vector{}
	for i := range len(*b.points) {
		u := (*b.points)[i]
		v := (*b.points)[(i+1)%len((*b.points))]
		normals = append(normals, (v.sub(&u)).normal())
	}
	return normals
}

func (b *BoundingBox) convertToWorldSpace(v *Vector) *Vector {
	// transform to world space
	u := v.rotate(b.theta)
	return &Vector{
		x: u.x + b.origin.x,
		y: u.y + b.origin.y,
	}
}

func (b *BoundingBox) projectionRange(v *Vector) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	for _, point := range *b.points {
		w := b.convertToWorldSpace(&point)

		// scalar projection
		s := w.dot(v) / v.length()
		min = math.Min(s, min)
		max = math.Max(s, max)
	}
	return min, max
}

func (b1 *BoundingBox) DidCollide(b2 *BoundingBox) bool {
	// separatable axis theorem

	// for each line segment,
	// get a perpendicular vector
	p := make(map[float64]bool)
	for _, normal := range b1.normals() {
		r := normal.rotate(b1.theta)
		// take gradients to deduplicate parallel lines
		g := r.gradient()
		p[g] = true
	}
	for _, normal := range b2.normals() {
		r := normal.rotate(b2.theta)
		g := r.gradient()
		p[g] = true
	}

	// for each perpendicular vector,
	// calculate the projection
	for g := range p {
		perpendicular := Vector{
			x: 1,
			y: g,
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
