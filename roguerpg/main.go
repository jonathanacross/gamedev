package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 384 * 3
	ScreenHeight = 900 // 240 * 3
)

type Game struct {
	terrain [][]*Tile
}

func NewGame() *Game {
	terrain := BuildLevel(70, 50)
	return &Game{
		terrain: terrain,
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, row := range g.terrain {
		for _, tile := range row {
			tile.Draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("GoCycle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
