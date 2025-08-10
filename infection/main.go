package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 508
	ScreenHeight = 680
	TileSize     = 64
)

func main() {
	game := NewGame()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Infection")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
