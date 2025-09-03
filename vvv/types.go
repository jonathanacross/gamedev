package main

import "github.com/hajimehoshi/ebiten/v2"

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

func (r Rect) Offset(x, y float64) Rect {
	return Rect{
		left:   r.left + x,
		top:    r.top + y,
		right:  r.right + x,
		bottom: r.bottom + y,
	}
}

func (r1 Rect) Intersects(r2 Rect) bool {
	return r1.left < r2.right && r1.right > r2.left &&
		r1.top < r2.bottom && r1.bottom > r2.top
}

// GameObject is an interface for any interactive entity in the game world.
type GameObject interface {
	HitBox() Rect
	Update()
}

// Drawable is for any GameObject that needs to be drawn every frame.
type Drawable interface {
	Draw(screen *ebiten.Image)
}
