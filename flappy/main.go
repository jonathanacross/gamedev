package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 216
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
