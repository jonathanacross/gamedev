package main

type Location struct {
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

func (r Rect) Offset(x, y float64) Rect {
	return Rect{
		Left:   r.Left + x,
		Top:    r.Top + y,
		Right:  r.Right + x,
		Bottom: r.Bottom + y,
	}
}

func (r1 Rect) Intersects(r2 Rect) bool {
	return r1.Left < r2.Right && r1.Right > r2.Left &&
		r1.Top < r2.Bottom && r1.Bottom > r2.Top
}

// GameObject is an interface for any interactive entity in the game world.
type GameObject interface {
	HitBox() Rect
	Update()
}
