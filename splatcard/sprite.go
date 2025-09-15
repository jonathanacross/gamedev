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

// BaseSprite provides common fields and methods for any visible game entity.
// It handles drawing a single sprite or the current frame of an animation.
type BaseSprite struct {
	Location
	image   *ebiten.Image
	srcRect image.Rectangle
	hitbox  Rect
}

// HitBox returns the collision rectangle for the BaseSprite.
func (bs *BaseSprite) HitBox() Rect {
	return bs.hitbox
}

func (bs *BaseSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X, bs.Y)
	currImage := bs.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}
