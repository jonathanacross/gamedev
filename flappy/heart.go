package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	HeartWidth  = 11
	HeartHeight = 10
)

var heartInstance = newHeart()

type Heart struct {
	spriteSheet *SpriteSheet
}

func newHeart() Heart {
	return Heart{
		spriteSheet: NewSpriteSheet(HeartImage, HeartWidth, HeartHeight, 2, 1),
	}
}

func DrawHeart(screen *ebiten.Image, x float64, y float64, filled bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	frame := 0
	if !filled {
		frame = 1
	}
	subRect := heartInstance.spriteSheet.Rect(frame)
	currImage := heartInstance.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}
