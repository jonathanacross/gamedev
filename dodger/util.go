package main

import ()

type Vector struct {
	X float64
	Y float64
}

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

func (r Rect) Intersects(other Rect) bool {
	return r.Left <= other.Right && other.Left <= r.Right &&
		r.Top <= other.Bottom && other.Top <= r.Bottom
}
