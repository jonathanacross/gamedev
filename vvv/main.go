package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 240
	TileSize     = 16

	Gravity      = 0.5
	RunSpeed     = 1.67
	MaxFallSpeed = 6

	StartLevelId = 1
	NumCrystals  = 3
)

// PlayerAction defines the type of action the player is requesting.
type PlayerAction int

const (
	NoAction PlayerAction = iota // Default, no specific action requested
	RespawnAction
	SwitchLevelAction
	CheckpointReachedAction
	WinGameAction
)

// PlayerActionEvent bundles the action type and any associated data.
type PlayerActionEvent struct {
	Action  PlayerAction
	Payload interface{} // e.g., LevelExit for SwitchLevelAction
}

type GameState int

const (
	StateTitleScreen GameState = iota
	StateInGame
	StateWinScreen
)

// Game is the main game struct.
type Game struct {
	player           *Player
	currentLevelNum  int // Store the number of the current level
	currentLevel     *Level
	gravity          float64
	allCheckpoints   map[int]*Checkpoint
	activeCheckpoint *Checkpoint
	debug            bool
	state            GameState
}

func (g *Game) UpdateInGame() {
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
		exit := actionEvent.Payload.(LevelExit)
		g.switchLevel(exit)
	case CheckpointReachedAction:
		newCheckpoint := actionEvent.Payload.(*Checkpoint)
		g.SetActiveCheckpoint(newCheckpoint)
	case WinGameAction:
		g.state = StateWinScreen

	case NoAction:
		// Do nothing
	}
}

func (g *Game) UpdateTitleScreen() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.StartNewGame()
	}
}

func (g *Game) Update() error {
	switch g.state {
	case StateTitleScreen:
		g.UpdateTitleScreen()
	case StateInGame:
		g.UpdateInGame()
	case StateWinScreen:
		g.UpdateTitleScreen()
	}

	PlayMusic()

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

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateTitleScreen:
		g.DrawTitleScreen(screen)
	case StateInGame:
		g.currentLevel.Draw(screen, g.debug)
		g.player.Draw(screen, g.debug)
	case StateWinScreen:
		g.DrawWinScreen(screen)
	}
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
		g.player.X = float64(ScreenWidth) - g.player.HitBox().Width() - 10.0 // Start at the right of the new screen
	} else if exit.bottom >= float64(ScreenHeight-1) { // Exit at the bottom of the screen
		g.player.Y = 10.0 // Start at the top of the new screen
	} else if exit.top <= 1 { // Exit at the top of the screen
		g.player.Y = float64(ScreenHeight) - g.player.HitBox().Height() - 10.0 // Start at the bottom of the new screen
	}
}

func (g *Game) Respawn() {
	if g.activeCheckpoint != nil {
		cp := g.activeCheckpoint
		if cp.LevelNum != g.currentLevelNum {
			g.currentLevelNum = cp.LevelNum
			g.currentLevel = LoadedLevels[g.currentLevelNum]
		}
		g.player.X, g.player.Y = cp.X, cp.Y
		g.player.numDeaths++
	} else {
		// This should only happen at the start of the game
		g.player.X, g.player.Y = g.currentLevel.startPoint.X, g.currentLevel.startPoint.Y
		g.player.numDeaths = 0
	}
}

func (g *Game) SetActiveCheckpoint(cp *Checkpoint) {
	// Deactivate all checkpoints first
	for _, checkpoint := range g.allCheckpoints {
		checkpoint.SetActive(false)
	}

	// Now activate the new checkpoint
	cp.SetActive(true)
	g.activeCheckpoint = cp
}

func (g *Game) DrawTitleScreen(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(StartScreen, op)

	drawTextAt(screen, "VVV", ScreenWidth/2, ScreenHeight/6, text.AlignCenter)
	drawTextAt(screen, "By Jonathan Cross", ScreenWidth/2, ScreenHeight/6+10, text.AlignCenter)
	drawTextAt(screen, "Collect the crystals.", 40, 74, text.AlignStart)
	drawTextAt(screen, "Use left/right arrows to move", 40, 90, text.AlignStart)
	drawTextAt(screen, "Use Space to reverse gravity", 40, 106, text.AlignStart)
	drawTextAt(screen, "Press Return to start game", 40, 122, text.AlignStart)
}

func (g *Game) DrawWinScreen(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(WinScreen, op)

	deathMessage := fmt.Sprintf("Number of deaths: %d", g.player.numDeaths)
	drawTextAt(screen, "You Win!", ScreenWidth/2, ScreenHeight/6+10, text.AlignCenter)
	drawTextAt(screen, deathMessage, 40, 90, text.AlignStart)
	drawTextAt(screen, "Press Return to play again", 40, 122, text.AlignStart)
}

func (g *Game) StartNewGame() {
	// Reload the levels to restore the initial state
	for levelNum, tiledMap := range Levels {
		LoadedLevels[levelNum] = NewLevel(tiledMap, levelNum)
	}
	startLevel, ok := LoadedLevels[StartLevelId]
	if !ok {
		panic("starting level not found")
	}

	g.player.Reset()
	g.currentLevel = startLevel
	g.gravity = Gravity

	// Set the initial player position
	for _, level := range LoadedLevels {
		for _, obj := range level.objects {
			if cp, ok := obj.(*Checkpoint); ok {
				g.allCheckpoints[cp.Id] = cp
				if cp.Active {
					g.activeCheckpoint = cp
					g.player.X, g.player.Y = cp.X, cp.Y
				}
			}
		}
	}

	g.state = StateInGame
}

// NewGame creates and initializes a new Game struct.
func NewGame() *Game {
	g := &Game{
		player:           NewPlayer(),
		currentLevelNum:  1,
		currentLevel:     nil,
		gravity:          Gravity,
		allCheckpoints:   make(map[int]*Checkpoint),
		debug:            false,
		state:            StateTitleScreen,
		activeCheckpoint: nil,
	}

	return g
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
