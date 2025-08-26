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
	spriteSheet *GridTileSet
	srcRect     image.Rectangle
	hitbox      Rect
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
	currImage := bs.spriteSheet.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type FlippableSprite struct {
	BaseSprite
	flipHoriz bool
	flipVert  bool
}

// Draw method for the flippable sprite. It handles the flipping logic.
func (f *FlippableSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// Apply horizontal flip
	if f.flipHoriz {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(f.spriteSheet.tileWidth), 0)
	}

	// Apply vertical flip
	if f.flipVert {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(f.spriteSheet.tileHeight))
	}

	// Translate to the sprite's position
	op.GeoM.Translate(f.X, f.Y)

	// Draw the sprite
	currImage := f.spriteSheet.image.SubImage(f.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)

	// show hitbox
	hb := f.FlippedHitbox()
	DrawRectFrame(screen, hb, color.RGBA{255, 255, 255, 255})
}

// FlippedHitbox returns the transformed hitbox based on the current
// flipping state.
func (f *FlippableSprite) FlippedHitbox() Rect {
	// TODO: should be using base hitbox
	// Start with the base hitbox
	hitbox := Rect{
		left:   f.X,
		top:    f.Y,
		right:  f.X + float64(f.spriteSheet.tileWidth),
		bottom: f.Y + float64(f.spriteSheet.tileHeight),
	}

	// TODO: flip around, keeping base hitbox intact
	// Apply horizontal flip if needed
	// if f.flipHoriz {
	// 	hitbox.left = f.X - float64(f.spriteSheet.tileWidth)
	// 	hitbox.right = f.X
	// }

	// // Apply vertical flip if needed
	// if f.flipVert {
	// 	hitbox.top = f.Y - float64(f.spriteSheet.tileHeight)
	// 	hitbox.bottom = f.Y
	// }

	return hitbox
}
