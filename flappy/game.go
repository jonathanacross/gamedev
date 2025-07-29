package main

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	player *Player
	camera *Camera

	world *ebiten.Image

	background *Background
	tiles      []*Tile
	items      []*Item
	enemies    []Enemy

	// TODO: update to level generator
	level *Level

	score int
}

func (g *Game) CheckCollisions() {
	for _, t := range g.tiles {
		if t.HitRect().Overlaps(g.player.HitRect()) {
			g.player.DoHit()
		}
	}
	for _, e := range g.enemies {
		if e.HitRect().Overlaps(g.player.HitRect()) {
			g.player.DoHit()
		}
	}

	for j, item := range g.items {
		if item.HitRect().Overlaps(g.player.HitRect()) {
			g.items = append(g.items[:j], g.items[j+1:]...)
			g.score++
		}
	}

}

func (g *Game) Update() error {
	g.camera.Center(Location{X: g.player.Location.X - ScreenWidth/4, Y: 0})
	g.level.Update(g.camera, &g.tiles, &g.items, &g.enemies)

	g.CheckCollisions()

	g.background.Update()
	g.player.Update()
	for _, i := range g.items {
		i.Update()
	}
	for _, e := range g.enemies {
		e.Update()
	}
	return nil
}

func (g *Game) drawScore(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	fontSize := float64(8)
	op.GeoM.Translate(ScreenWidth/2, ScoreOffset+1)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter

	scoreText := fmt.Sprintf("%05d", g.score)
	text.Draw(screen, scoreText, &text.GoTextFace{
		Source: ArcadeFaceSource,
		Size:   fontSize,
	}, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.background.Draw(screen)

	for _, t := range g.tiles {
		t.Draw(g.camera, screen)
	}
	for _, i := range g.items {
		i.Draw(g.camera, screen)
	}
	for _, e := range g.enemies {
		e.Draw(g.camera, screen)
	}

	g.player.Draw(g.camera, screen)
	g.drawScore(screen)
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
		score:      0,
	}
}
