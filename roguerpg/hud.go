package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	HeartWidth        = 11
	HeartHeight       = 10
	HeartSubdivisions = 2
)

var heartInstance = newHeart()

type Heart struct {
	image       *ebiten.Image
	spriteSheet *SpriteSheet
}

func newHeart() Heart {
	return Heart{
		image:       HealthHeartImage,
		spriteSheet: NewSpriteSheet(HeartWidth, HeartHeight, 3, 1),
	}
}

func DrawHeart(screen *ebiten.Image, x float64, y float64, frame int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	subRect := heartInstance.spriteSheet.Rect(frame)
	currImage := heartInstance.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

func DrawPlayerHeath(screen *ebiten.Image, currHeath int, maxHeath int) {
	numHearts := (maxHeath + (HeartSubdivisions - 1)) / HeartSubdivisions

	for i := range numHearts {
		x := float64(20 + i*HeartWidth)
		y := float64(20)
		frame := clamp(currHeath-HeartSubdivisions*i, 0, HeartSubdivisions)
		DrawHeart(screen, x, y, frame)
	}
}

func DrawHeadsUpDisplay(screen *ebiten.Image, player *Player) {
	DrawPlayerHeath(screen, player.Health, player.MaxHealth)
}
