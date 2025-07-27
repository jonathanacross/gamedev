package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player *Player

	background *Background
	tiles      []*Tile
	items      []*Item
	enemies    []*Enemy
}

var gameInstance = NewGame()

func (g *Game) Update() error {
	g.background.Update()
	g.player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.background.Draw(screen)
	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	return &Game{
		player:     NewPlayer(),
		background: NewBackground(),
		tiles:      nil,
		items:      nil,
		enemies:    nil,
	}
}
