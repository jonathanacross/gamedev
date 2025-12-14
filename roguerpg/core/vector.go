package core

import (
	"math"
)

type Location Vector

type Vector struct {
	X float64
	Y float64
}

func ZeroVector() Vector {
	return Vector{X: 0, Y: 0}
}

func (v Vector) Minus(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector) Plus(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v Vector) Normalize() Vector {
	length := v.Length()
	if length == 0 {
		return Vector{0, 0}
	}
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
	}
}

func (v Vector) Scale(scalar float64) Vector {
	return Vector{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v Vector) Rotate(angleRadians float64) Vector {
	s := math.Sin(angleRadians)
	c := math.Cos(angleRadians)
	return Vector{
		X: v.X*c + v.Y*s,
		Y: -v.X*s + v.Y*c,
	}
}
