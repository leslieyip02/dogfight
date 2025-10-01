package geometry

import (
	"math"
)

const epsilon = 1e-5

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func NewVector(x float64, y float64) Vector {
	return Vector{X: x, Y: y}
}

func (u *Vector) sub(v *Vector) *Vector {
	return &Vector{
		X: u.X - v.X,
		Y: u.Y - v.Y,
	}
}

func (u *Vector) multiply(s float64) *Vector {
	return &Vector{
		X: s * u.X,
		Y: s * u.Y,
	}
}

func (u *Vector) dot(v *Vector) float64 {
	return u.X*v.X + u.Y*v.Y
}

func (u *Vector) length() float64 {
	return math.Sqrt(u.dot(u))
}

func (u *Vector) gradient() float64 {
	return u.Y / u.X
}

func (u *Vector) angle() float64 {
	return math.Atan2(u.Y, u.X)
}

func (u *Vector) rotate(theta float64) *Vector {
	h := math.Sqrt(u.dot(u))
	theta += u.angle()
	return &Vector{
		X: math.Cos(theta) * h,
		Y: math.Sin(theta) * h,
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
