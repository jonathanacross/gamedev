package main

type Location struct {
	X float64
	Y float64
}

type Rect struct {
	left   float64
	top    float64
	right  float64
	bottom float64
}

func (r Rect) Width() float64 {
	return r.right - r.left
}

func (r Rect) Height() float64 {
	return r.bottom - r.top
}

func (r1 *Rect) Intersects(r2 *Rect) bool {
	return r1.left < r2.right && r1.right > r2.left &&
		r1.top < r2.bottom && r1.bottom > r2.top
}
