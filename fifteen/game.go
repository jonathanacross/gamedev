package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Layout
const (
	ScreenWidth  = 552
	ScreenHeight = 600
	TileSize     = 128
	GridMargin   = 20
	ButtonWidth  = 150
	ButtonHeight = 40
	ButtonX      = 200
	ButtonY      = 545
)

type Game struct {
	grid          *TileGrid
	shuffleButton *Button
}

func getPicture() *ebiten.Image {
	if len(os.Args) == 2 && os.Args[1] == "--red" {
		return AltPicture
	} else {
		return Picture
	}
}

func NewGame() *Game {
	grid := NewTileGrid(getPicture())
	grid.Randomize()

	return &Game{
		grid:          grid,
		shuffleButton: NewButton(ButtonX, ButtonY, ButtonWidth, ButtonHeight, grid.Randomize),
	}
}

func (g *Game) Update() error {
	g.grid.Update()
	g.shuffleButton.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.grid.Draw(screen)
	g.shuffleButton.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
