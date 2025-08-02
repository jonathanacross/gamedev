package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 300
	ScreenHeight = 400
	TileSize     = 32
)

func main() {
	game := NewGame()
	ebiten.SetWindowSize(2*ScreenWidth, 2*ScreenHeight)
	ebiten.SetWindowTitle("Infection")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
