package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func NewRect(ir image.Rectangle) Rect {
	return Rect{
		left:   float64(ir.Min.X),
		top:    float64(ir.Min.Y),
		right:  float64(ir.Max.X),
		bottom: float64(ir.Max.Y),
	}
}

// Intersects checks if two rectangles overlap.
func (r Rect) Intersects(other Rect) bool {
	return r.right > other.left && r.left < other.right && r.bottom > other.top && r.top < other.bottom
}

// BaseSprite provides common fields and methods for any visible game entity.
// It handles drawing a single sprite or the current frame of an animation.
type BaseSprite struct {
	Location
	image   *ebiten.Image
	srcRect image.Rectangle
	hitbox  Rect
}

func (bs *BaseSprite) HitBox() Rect {
	return Rect{
		left:   bs.X + bs.hitbox.left,
		top:    bs.Y + bs.hitbox.top,
		right:  bs.X + bs.hitbox.right,
		bottom: bs.Y + bs.hitbox.bottom,
	}
}

func (bs *BaseSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X, bs.Y)
	currImage := bs.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

// HasCollided checks for collision with another BaseSprite
func (bs *BaseSprite) HasCollided(other *BaseSprite) bool {
	return bs.HitBox().Intersects(other.HitBox())
}
