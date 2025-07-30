package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 240
	TileSize     = 16
	ScoreOffset  = 3

	Gravity               = 0.1
	JumpVelocity          = -2.5
	PlayerSpeed           = 1
	BackgroundScrollSpeed = 0.25
	PlayerMaxHealth       = 3
	DropHeartsEveryNCoins = 20
)

func main() {
	game := NewGame()
	//ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)
	ebiten.SetWindowTitle("Flappy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
