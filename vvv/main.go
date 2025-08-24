package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game is the main game struct.
type Game struct {
	player          *Player
	currentLevelNum int // Store the number of the current level
	currentLevel    *Level
	gravity         float64
}

const (
	ScreenWidth  = 384
	ScreenHeight = 240
)

func (g *Game) Update() error {
	if g.player.IsOnGround() && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.gravity *= -1
	}

	g.currentLevel.Update()
	g.player.Update(g)

	// Check for collision with level exits
	g.checkLevelExits()

	PlayMusic()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentLevel.Draw(screen)
	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// checkLevelExits checks for a collision between the player and any level exits.
func (g *Game) checkLevelExits() {
	playerHitRect := g.player.HitRect()

	for _, exit := range g.currentLevel.exits {
		if playerHitRect.Intersects(&exit.Rect) {
			log.Printf("Player collided with exit to level %d", exit.ToLevel)
			g.switchLevel(exit)
			return // Exit after the first collision is found and processed
		}
	}
}

func (g *Game) switchLevel(exit LevelExit) {
	newLevel, ok := LoadedLevels[exit.ToLevel]
	if !ok {
		log.Printf("Level %d not found in loaded levels.", exit.ToLevel)
		return
	}
	g.currentLevel = newLevel
	g.currentLevelNum = exit.ToLevel

	// Determine the transition direction and adjust the player's position.
	// You can add properties to the LevelExit object to specify the entry point,
	// but for now, we'll assume a simple screen-to-screen transition.
	if exit.right >= float64(ScreenWidth-1) { // Exit on the right side of the screen
		g.player.X = 10.0 // Start at the left of the new screen
	} else if exit.left <= 1 { // Exit on the left side of the screen
		g.player.X = float64(ScreenWidth) - g.player.HitRect().Width() - 10.0 // Start at the right of the new screen
	} else if exit.bottom >= float64(ScreenHeight-1) { // Exit at the bottom of the screen
		g.player.Y = 10.0 // Start at the top of the new screen
	} else if exit.top <= 1 { // Exit at the top of the screen
		g.player.Y = float64(ScreenHeight) - g.player.HitRect().Height() - 10.0 // Start at the bottom of the new screen
	}
}

func main() {
	spriteSheet := NewSpriteSheet(TileSet, TileSize, TileSize, 5, 7)

	// Pre-load all levels and their objects using the TilesetData.
	for levelNum, levelJSON := range Levels {
		LoadedLevels[levelNum] = NewLevel(levelJSON, TilesetData, spriteSheet)
	}

	// Set the initial level to 1
	startLevel, ok := LoadedLevels[1]
	if !ok {
		panic("starting level not found")
	}

	g := &Game{
		player:          NewPlayer(),
		currentLevelNum: 1,
		currentLevel:    startLevel,
		gravity:         0.5,
	}

	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
