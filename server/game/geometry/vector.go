package geometry

import (
	"math"
)

const epsilon = 1e-5

type Vector struct {
	x float64
	y float64
}

func (u *Vector) sub(v *Vector) *Vector {
	return &Vector{
		x: u.x - v.x,
		y: u.y - v.y,
	}
}

func (u *Vector) multiply(s float64) *Vector {
	return &Vector{
		x: s * u.x,
		y: s * u.y,
	}
}

func (u *Vector) dot(v *Vector) float64 {
	return u.x*v.x + u.y*v.y
}

func (u *Vector) length() float64 {
	return math.Sqrt(u.dot(u))
}

func (u *Vector) gradient() float64 {
	return u.y / u.x
}

func (u *Vector) angle() float64 {
	return math.Atan2(u.y, u.x)
}

func (u *Vector) rotate(theta float64) *Vector {
	h := math.Sqrt(u.dot(u))
	theta += u.angle()
	return &Vector{
		x: math.Cos(theta) * h,
		y: math.Sin(theta) * h,
	}
}

func (u *Vector) normal() *Vector {
	v := u.rotate(math.Pi / 2)
	return v.multiply(1 / v.length())
}

func (u *Vector) isParallel(v *Vector) bool {
	a := u.angle()
	b := v.angle()
	return math.Abs(math.Mod(a-b, math.Pi)) < epsilon
}
