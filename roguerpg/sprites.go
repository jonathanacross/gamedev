package main

import (
	"image"

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
	return bs.hitbox.Offset(bs.X, bs.Y)
}

func (bs *BaseSprite) GetX() float64 { return bs.X }

func (bs *BaseSprite) GetY() float64 { return bs.Y }

func (bs *BaseSprite) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X, bs.Y)
	op.GeoM.Concat(cameraMatrix)
	currImage := bs.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type Tile struct {
	BaseSprite
	solid bool
}
