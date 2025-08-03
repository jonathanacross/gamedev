package main

import (
	"fmt"
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

func (g *Game) doComputerMove(engine Engine) {
	move := engine.GenMove(g.boardWidget.gameBoard)
	g.boardWidget.DoComputerMove(move)
}

func (g *Game) getCurrentEngine() Engine {
	if g.boardWidget.gameBoard.playerToMove == White {
		return g.engineWhite
	} else {
		return g.engineBlack
	}
}

func (g *Game) Update() error {
	g.boardWidget.Update()
	g.spinner.Update()

	// check for change at end of move
	if g.state == ComputerThinking || g.state == WaitingForHuman {
		if g.prevPlayerToMove != g.boardWidget.gameBoard.playerToMove {
			fmt.Printf("Player to move: %v\n", g.boardWidget.gameBoard.playerToMove)
			g.prevPlayerToMove = g.boardWidget.gameBoard.playerToMove
			g.state = GameInProgress
		}
	}

	if g.state == GameInProgress {
		currEngine := g.getCurrentEngine()
		if currEngine.RequiresHumanInput() {
			g.state = WaitingForHuman
			g.boardWidget.allowUserInput = true
			g.spinner.SetVisible(false)
		} else {
			g.state = ComputerThinking
			g.boardWidget.allowUserInput = false
			g.spinner.SetVisible(true)
			go g.doComputerMove(currEngine)
		}
	}
	g.moveTimer.Update()

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
