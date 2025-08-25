package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PlayerAction defines the type of action the player is requesting.
type PlayerAction int

const (
	NoAction PlayerAction = iota // Default, no specific action requested
	RespawnAction
	SwitchLevelAction
)

// PlayerActionEvent bundles the action type and any associated data.
type PlayerActionEvent struct {
	Action  PlayerAction
	Payload interface{} // e.g., LevelExit for SwitchLevelAction
}

// Game is the main game struct.
type Game struct {
	player          *Player
	currentLevelNum int // Store the number of the current level
	currentLevel    *Level
	gravity         float64
	allCheckpoints  map[int]*Checkpoint
	debug           bool
}

const (
	ScreenWidth  = 384
	ScreenHeight = 240
)

func (g *Game) Update() error {
	// Toggle debug mode on backtick key press
	if inpututil.IsKeyJustPressed(ebiten.KeyBackquote) {
		g.debug = !g.debug
		log.Printf("Debug mode is now: %v\n", g.debug)
	}

	if g.player.IsOnGround() && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.gravity *= -1
	}

	g.currentLevel.Update()

	// Pass gravity directly to the player's Update method
	actionEvent := g.player.Update(g.currentLevel, g.gravity)

	// Process player actions
	switch actionEvent.Action {
	case RespawnAction:
		g.Respawn()
	case SwitchLevelAction:
		if exit, ok := actionEvent.Payload.(LevelExit); ok {
			g.switchLevel(exit)
		} else {
			log.Println("Error: SwitchLevelAction payload is not of type LevelExit")
		}
	case NoAction:
		// Do nothing
	}

	PlayMusic()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentLevel.Draw(screen, g.debug)
	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) switchLevel(exit LevelExit) {
	g.currentLevelNum = exit.ToLevel
	g.currentLevel = LoadedLevels[g.currentLevelNum]

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

func (g *Game) Respawn() {
	if g.player.activeCheckpoint != nil {
		cp := g.player.activeCheckpoint
		if cp.LevelNum != g.currentLevelNum {
			g.currentLevelNum = cp.LevelNum
			g.currentLevel = LoadedLevels[g.currentLevelNum]
		}
		g.player.X, g.player.Y = cp.X, cp.Y
	} else {
		// This should only happen at the start of the game
		g.player.X, g.player.Y = g.currentLevel.startPoint.X, g.currentLevel.startPoint.Y
	}
}

func main() {
	spriteSheet := NewSpriteSheet(TileSet, TileSize, TileSize, 5, 7)

	// Pre-load all levels and their objects using the TilesetData.
	for levelNum, levelJSON := range Levels {
		LoadedLevels[levelNum] = NewLevel(levelJSON, TilesetData, spriteSheet, levelNum)
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
		allCheckpoints:  make(map[int]*Checkpoint),
	}

	for _, level := range LoadedLevels {
		for _, cp := range level.checkpoints {
			g.allCheckpoints[cp.Id] = cp
			if cp.Active {
				g.player.activeCheckpoint = cp
				g.player.X, g.player.Y = cp.X, cp.Y
				g.player.checkpointId = cp.Id
			}
		}
	}

	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
