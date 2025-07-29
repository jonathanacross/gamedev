package main

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameState int

const (
	TitleScreen GameState = iota
	InGame
	GameOver
)

type Game struct {
	state GameState

	player *Player
	camera *Camera

	world *ebiten.Image

	background *Background
	tiles      []*Tile
	items      []*Item
	enemies    []Enemy

	// TODO: update to level generator
	level *Level

	score   int
	hiScore int
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

func (g *Game) updateTitleScreen() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.StartGame()
		g.state = InGame
	}
}

func (g *Game) updateGameOver() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state = TitleScreen

		if g.score > g.hiScore {
			g.hiScore = g.score
		}
	}
}

func (g *Game) updateInGame() {
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

	if g.player.health < 0 {
		g.state = GameOver
	}
}

func (g *Game) Update() error {
	switch g.state {
	case TitleScreen:
		g.updateTitleScreen()
	case InGame:
		g.updateInGame()
	case GameOver:
		g.updateGameOver()
	}

	PlayMusic()
	return nil
}

func drawTextAt(screen *ebiten.Image, message string, x float64, y float64, align text.Align) {
	op := &text.DrawOptions{}
	fontSize := float64(8)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = align

	text.Draw(screen, message, &text.GoTextFace{
		Source: ArcadeFaceSource,
		Size:   fontSize,
	}, op)
}

func (g *Game) drawTitleScreen(screen *ebiten.Image) {
	g.background.Draw(screen)
	scoreString := fmt.Sprintf("High score: %04d", g.hiScore)
	drawTextAt(screen, "Flappy", ScreenWidth/2, ScreenHeight/6, text.AlignCenter)
	drawTextAt(screen, "By Jonathan Cross", ScreenWidth/2, ScreenHeight/6+10, text.AlignCenter)
	drawTextAt(screen, scoreString, ScreenWidth/2, ScreenHeight/2, text.AlignCenter)
	drawTextAt(screen, "Press Space", ScreenWidth/2, ScreenHeight*4/5, text.AlignCenter)
}

func (g *Game) drawGame(screen *ebiten.Image) {
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

func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.drawGame(screen)
	drawTextAt(screen, "GAME OVER", ScreenWidth/2, ScreenHeight/2, text.AlignCenter)
	drawTextAt(screen, "Press Space", ScreenWidth/2, ScreenHeight*4/5, text.AlignCenter)
}

func (g *Game) drawScore(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("%04d", g.score)
	drawTextAt(screen, scoreText, ScreenWidth/2, ScoreOffset+1, text.AlignCenter)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case TitleScreen:
		g.drawTitleScreen(screen)
	case InGame:
		g.drawGame(screen)
	case GameOver:
		g.drawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) StartGame() {
	g.player = NewPlayer()
	g.camera = NewCamera()
	g.tiles = []*Tile{}
	g.items = nil
	g.enemies = nil
	g.level = NewLevel()
	g.score = 0
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
