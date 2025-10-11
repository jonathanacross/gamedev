package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ArenaWidth   = 50
	ArenaHeight  = 50
	ArenaOffsetX = 90
	ArenaOffsetY = 30
	SquareSize   = 4

	GameUpdateSpeedMillis = 100
	NumRounds             = 5

	ScreenWidth  = 384
	ScreenHeight = 240
	TileSize     = 16
	ScoreOffset  = 3

	CharPortraitWidth       = 64
	CharPortraitBigHeight   = 80
	CharPortraitSmallHeight = 64
)

func main() {
	game := NewGame()
	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)
	ebiten.SetWindowTitle("GoCycle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
