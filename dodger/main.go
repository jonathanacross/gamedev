package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Dodger")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
