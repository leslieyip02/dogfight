package geometry

import (
	"math"
	"math/rand"
)

const EPSILON = 1e-5

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func NewVector(x float64, y float64) *Vector {
	return &Vector{X: x, Y: y}
}

func NewRandomVector(minX float64, minY float64, maxX float64, maxY float64) *Vector {
	return &Vector{
		X: minX + rand.Float64()*(maxX-minX),
		Y: minY + rand.Float64()*(maxY-minY),
	}
}

func (u *Vector) Add(v *Vector) *Vector {
	return &Vector{
		X: u.X + v.X,
		Y: u.Y + v.Y,
	}
}

func (u *Vector) Sub(v *Vector) *Vector {
	return &Vector{
		X: u.X - v.X,
		Y: u.Y - v.Y,
	}
}

func (u *Vector) Multiply(s float64) *Vector {
	return &Vector{
		X: s * u.X,
		Y: s * u.Y,
	}
}

func (u *Vector) Length() float64 {
	return math.Sqrt(u.dot(u))
}

func (u *Vector) Unit() *Vector {
	if u.Length() == 0 {
		return &Vector{0, 0}
	}
	return u.Multiply(1 / u.Length())
}

func (u *Vector) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

func (u *Vector) Rotate(theta float64) *Vector {
	h := math.Sqrt(u.dot(u))
	theta += u.Angle()
	return &Vector{
		X: math.Cos(theta) * h,
		Y: math.Sin(theta) * h,
	}
}

func (u *Vector) Normal() *Vector {
	return u.Rotate(math.Pi / 2).Unit()
}

func (u *Vector) dot(v *Vector) float64 {
	return u.X*v.X + u.Y*v.Y
}

func (u *Vector) cross(v *Vector) float64 {
	return u.X*v.Y - u.Y*v.X
}

func (u *Vector) gradient() float64 {
	return u.Y / u.X
}

func (u *Vector) isParallel(v *Vector) bool {
	a := u.Angle()
	b := v.Angle()
	return math.Abs(math.Mod(a-b, math.Pi)) < EPSILON
}
