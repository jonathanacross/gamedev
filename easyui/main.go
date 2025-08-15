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
	demo := NewDemo()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("EasyUi Demo")
	if err := ebiten.RunGame(demo); err != nil {
		log.Fatal(err)
	}
}
