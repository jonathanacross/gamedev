package main

import (
	_ "image/png" // Import for image decoding
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameState int

const (
	GameNotStarted GameState = iota
	GameInProgress
	WaitingForHuman
	ComputerThinking
	AnimatingComputerMove
	GameOver
)

type Game struct {
	boardWidget      *BoardWidget
	moveTimer        *Timer
	engineWhite      Engine
	engineBlack      Engine
	spinner          *Spinner
	state            GameState
	prevPlayerToMove Player
}

func NewGame() *Game {
	return &Game{
		boardWidget:      NewBoardWidget(),
		moveTimer:        NewTimer(2 * time.Second),
		engineWhite:      &Human{},
		engineBlack:      &SlowEngine{},
		spinner:          NewSpinner(),
		state:            GameInProgress,
		prevPlayerToMove: Black,
	}
}

func (g *Game) getCurrentEngine() Engine {
	if g.boardWidget.gameBoard.playerToMove == White {
		return g.engineWhite
	} else {
		return g.engineBlack
	}
}

func (g *Game) Update() error {
	switch g.state {
	case GameInProgress:
		// Logic to determine whose turn it is
		g.boardWidget.allowUserInput = false // Ensure input is off by default
		g.spinner.SetVisible(false)          // Ensure spinner is off by default

		currEngine := g.getCurrentEngine()
		if currEngine.RequiresHumanInput() {
			g.state = WaitingForHuman
		} else {
			g.state = ComputerThinking
		}

	case WaitingForHuman:
		g.boardWidget.allowUserInput = true
		g.boardWidget.Update()

		// If the player has changed (a move was made), transition
		if g.boardWidget.gameBoard.playerToMove != g.prevPlayerToMove {
			g.prevPlayerToMove = g.boardWidget.gameBoard.playerToMove
			g.state = GameInProgress
		}

	case ComputerThinking:
		g.spinner.SetVisible(true)
		// Start the computer move in a goroutine
		go func() {
			currEngine := g.getCurrentEngine()
			move := currEngine.GenMove(g.boardWidget.gameBoard)
			g.boardWidget.DoComputerMove(move)
			g.state = AnimatingComputerMove
		}()

	case AnimatingComputerMove:
		g.boardWidget.UpdateComputerDragInfo()
		if !g.boardWidget.computerDragInfo.isAnimating {
			g.prevPlayerToMove = g.boardWidget.gameBoard.playerToMove
			g.state = GameInProgress
		}
	}
	g.spinner.Update() // Spinner update can be outside, as it's a general animation
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.boardWidget.Draw(screen)
	g.spinner.Draw(screen)
}

// Layout returns the game's logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
