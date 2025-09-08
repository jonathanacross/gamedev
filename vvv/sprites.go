package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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

// GetX returns the X coordinate of the BaseSprite.
func (bs *BaseSprite) GetX() float64 { return bs.X }

// GetY returns the Y coordinate of the BaseSprite.
func (bs *BaseSprite) GetY() float64 { return bs.Y }

func (bs *BaseSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X, bs.Y)
	currImage := bs.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type FlippableSprite struct {
	BaseSprite
	flipHoriz bool
	flipVert  bool
}

// Draw method for the flippable sprite. It handles the flipping logic.
func (f *FlippableSprite) Draw(screen *ebiten.Image, debug bool) {
	op := &ebiten.DrawImageOptions{}

	// Apply horizontal flip
	if f.flipHoriz {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(f.srcRect.Dx()), 0)
	}

	// Apply vertical flip
	if f.flipVert {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(f.srcRect.Dy()))
	}

	// Translate to the sprite's position
	op.GeoM.Translate(f.X, f.Y)

	// Draw the sprite
	currImage := f.image.SubImage(f.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)

	if debug {
		// show hitbox
		hb := f.FlippedHitbox()
		DrawRectFrame(screen, hb, color.RGBA{255, 255, 255, 255})
	}
}

// FlippedHitbox returns the transformed hitbox based on the current
// flipping state.
func (f *FlippableSprite) FlippedHitbox() Rect {
	box := Rect{
		left:   f.X + f.hitbox.left,
		top:    f.Y + f.hitbox.top,
		right:  f.X + f.hitbox.right,
		bottom: f.Y + f.hitbox.bottom,
	}

	// Apply horizontal flip if needed
	if f.flipHoriz {
		box.left = f.X + float64(f.srcRect.Dx()) - f.hitbox.right
		box.right = f.X + float64(f.srcRect.Dx()) - f.hitbox.left
	}

	// Apply vertical flip if needed
	if f.flipVert {
		box.top = f.Y + float64(f.srcRect.Dy()) - f.hitbox.bottom
		box.bottom = f.Y + float64(f.srcRect.Dy()) - f.hitbox.top
	}

	return box
}
