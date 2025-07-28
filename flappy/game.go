package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player *Player

	background *Background
	ground     []*Tile
	obstacles  []*Tile
	items      []*Item
	enemies    []*Enemy

	// TODO: update to level generator
	level *Level
}

func (g *Game) Update() error {
	g.level.Update(&g.ground, &g.obstacles)

	g.background.Update()
	g.player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.background.Draw(screen)

	for _, t := range g.ground {
		t.Draw(screen)
	}
	for _, t := range g.obstacles {
		t.Draw(screen)
	}

	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	return &Game{
		player:     NewPlayer(),
		background: NewBackground(),
		ground:     []*Tile{},
		obstacles:  []*Tile{},
		items:      nil,
		enemies:    nil,
		level:      NewLevel(),
	}
}
