package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player *Player
	camera *Camera

	world *ebiten.Image

	background *Background
	tiles      []*Tile
	items      []*Item
	enemies    []*Enemy

	// TODO: update to level generator
	level *Level
}

func (g *Game) Update() error {
	g.camera.Center(Location{X: g.player.Location.X - ScreenWidth/4, Y: 0})
	g.level.Update(g.camera, &g.tiles)

	g.background.Update()
	g.player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.background.Draw(screen)

	for _, t := range g.tiles {
		t.Draw(g.camera, screen)
	}

	g.player.Draw(g.camera, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	return &Game{
		player:     NewPlayer(),
		camera:     NewCamera(),
		background: NewBackground(),
		tiles:      []*Tile{},
		items:      nil,
		enemies:    nil,
		level:      NewLevel(),
	}
}
