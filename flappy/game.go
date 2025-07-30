package main

import (
	"fmt"
	"image/color"
	_ "image/png" // Import for image decoding

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// GameState defines the different states of the game.
type GameState int

const (
	TitleScreen GameState = iota
	InGame
	GameOver
)

// Game holds all the game's state and logic.
type Game struct {
	state GameState

	player *Player
	camera *Camera

	background *Background
	tiles      []*Tile // Static level geometry
	items      []Item  // Collectible items
	enemies    []Enemy // Enemies

	level *Level // Level generator

	score   int
	hiScore int
}

// CheckCollisions handles collision detection for player with tiles, items, and enemies.
func (g *Game) CheckCollisions() {
	// Player-Tile collisions
	for _, t := range g.tiles {
		if t.HitRect().Overlaps(g.player.HitRect()) {
			g.player.DoHit()
		}
	}
	// Player-Enemy collisions
	for _, e := range g.enemies {
		if e.HitRect().Overlaps(g.player.HitRect()) {
			g.player.DoHit()
		}
	}

	// Player-Item collisions (iterate backwards for safe deletion during loop)
	for j := len(g.items) - 1; j >= 0; j-- {
		item := g.items[j]
		if item.HitRect().Overlaps(g.player.HitRect()) {
			item.UseItem(g)
			// Remove the collected item from the slice
			g.items = append(g.items[:j], g.items[j+1:]...)
		}
	}
}

// updateTitleScreen handles input and transitions from the title screen.
func (g *Game) updateTitleScreen() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.StartGame()
		g.state = InGame
	}
}

// updateGameOver handles input and transitions from the game over screen.
func (g *Game) updateGameOver() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state = TitleScreen

		if g.score > g.hiScore {
			g.hiScore = g.score
		}
	}
}

// updateInGame handles the main game loop logic when the game is active.
func (g *Game) updateInGame() {
	// Center camera on player, offsetting the player to the left side
	g.camera.Center(Location{X: g.player.Location.X - ScreenWidth/4, Y: 0})
	// Update level generation based on camera position
	g.level.Update(g.camera, &g.tiles, &g.items, &g.enemies)

	// Update all dynamic entities
	g.background.Update()
	g.player.Update()
	for _, i := range g.items {
		i.Update()
	}
	for _, e := range g.enemies {
		e.Update()
	}

	g.CheckCollisions()

	if g.player.health < 0 {
		g.state = GameOver
	}
}

// Update is the main update loop for the entire game, delegating to state-specific updates.
func (g *Game) Update() error {
	switch g.state {
	case TitleScreen:
		g.updateTitleScreen()
	case InGame:
		g.updateInGame()
	case GameOver:
		g.updateGameOver()
	}

	PlayMusic() // Play background music
	return nil
}

// drawTextAt is a helper function to draw text on the screen with alignment.
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

// drawTitleScreen renders the title screen.
func (g *Game) drawTitleScreen(screen *ebiten.Image) {
	g.background.Draw(screen) // Draw background behind title text
	scoreString := fmt.Sprintf("High score: %04d", g.hiScore)
	drawTextAt(screen, "Flappy", ScreenWidth/2, ScreenHeight/6, text.AlignCenter)
	drawTextAt(screen, "By Jonathan Cross", ScreenWidth/2, ScreenHeight/6+10, text.AlignCenter)
	drawTextAt(screen, scoreString, ScreenWidth/2, ScreenHeight/2, text.AlignCenter)
	drawTextAt(screen, "Press Space", ScreenWidth/2, ScreenHeight*4/5, text.AlignCenter)
}

// drawGame renders all game elements during active gameplay.
func (g *Game) drawGame(screen *ebiten.Image) {
	g.background.Draw(screen)

	// Draw tiles (static level geometry)
	for _, t := range g.tiles {
		t.Draw(g.camera, screen)
	}
	// Draw items
	for _, i := range g.items {
		i.Draw(g.camera, screen)
	}
	// Draw enemies
	for _, e := range g.enemies {
		e.Draw(g.camera, screen)
	}

	g.player.Draw(g.camera, screen)
	g.drawScore(screen)
}

// drawGameOver renders the game over screen, showing the last game state.
func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.drawGame(screen) // Draw the last game frame
	drawTextAt(screen, "GAME OVER", ScreenWidth/2, ScreenHeight/2, text.AlignCenter)
	drawTextAt(screen, "Press Space", ScreenWidth/2, ScreenHeight*4/5, text.AlignCenter)
}

// drawScore renders the current score on the screen.
func (g *Game) drawScore(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("%04d", g.score)
	drawTextAt(screen, scoreText, ScreenWidth/2, ScoreOffset+1, text.AlignCenter)
}

// Draw is the main drawing loop for the entire game, delegating to state-specific draws.
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

// Layout returns the game's logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// StartGame initializes or resets the game state for a new play session.
func (g *Game) StartGame() {
	g.player = NewPlayer()
	g.camera = NewCamera()
	g.tiles = []*Tile{}
	g.items = []Item{}
	g.enemies = []Enemy{}
	g.level = NewLevel()
	g.score = 0
}

// NewGame creates and initializes a new Game instance.
func NewGame() *Game {
	return &Game{
		player:     NewPlayer(),
		camera:     NewCamera(),
		background: NewBackground(),
		tiles:      []*Tile{},
		items:      []Item{},
		enemies:    []Enemy{},
		level:      NewLevel(),
		score:      0,
	}
}
