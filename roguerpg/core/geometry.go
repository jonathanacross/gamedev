package core

import (
	"math"
)

type Rect struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (r Rect) Width() float64 {
	return r.Right - r.Left
}

func (r Rect) Height() float64 {
	return r.Bottom - r.Top
}

func (r Rect) Offset(x, y float64) Rect {
	return Rect{
		Left:   r.Left + x,
		Top:    r.Top + y,
		Right:  r.Right + x,
		Bottom: r.Bottom + y,
	}
}

func (r1 Rect) IntersectsX(r2 Rect) bool {
	return r1.Left < r2.Right && r1.Right > r2.Left
}

func (r1 Rect) IntersectsY(r2 Rect) bool {
	return r1.Top < r2.Bottom && r1.Bottom > r2.Top
}

func (r1 Rect) Intersects(r2 Rect) bool {
	return r1.IntersectsX(r2) && r1.IntersectsY(r2)
}

type CollisionAxis int

const (
	AxisX CollisionAxis = iota
	AxisY
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

func VectorToDirection(dirVector Vector) Direction {
	angle := math.Atan2(dirVector.Y, dirVector.X)
	if angle >= -math.Pi/4 && angle < math.Pi/4 {
		return Right
	} else if angle >= math.Pi/4 && angle < 3*math.Pi/4 {
		return Down
	} else if angle >= -3*math.Pi/4 && angle < -math.Pi/4 {
		return Up
	} else {
		return Left
	}
}

func DirectionToVector(dir Direction) Vector {
	switch dir {
	case Left:
		return Vector{X: -1, Y: 0}
	case Right:
		return Vector{X: 1, Y: 0}
	case Up:
		return Vector{X: 0, Y: -1}
	case Down:
		return Vector{X: 0, Y: 1}
	}
	return Vector{X: 0, Y: 0}
}
