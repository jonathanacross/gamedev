package main

import (
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
	state       GameState
	gameBoard   *Board
	boardWidget *BoardWidget
	whiteEngine Engine
	blackEngine Engine
	spinner     *Spinner
}

func NewGame() *Game {
	g := Game{
		state:       GameInProgress,
		gameBoard:   NewBoard(),
		boardWidget: NewBoardWidget(),
		whiteEngine: &HumanEngine{},
		blackEngine: &MinimaxEngine{maxDepth: 5},
		spinner:     NewSpinner(),
	}
	return &g
}

func (g *Game) getCurrentEngine() Engine {
	if g.gameBoard.playerToMove == White {
		return g.whiteEngine
	}
	return g.blackEngine
}

func (g *Game) Update() error {
	g.spinner.Update()

	switch g.state {
	case GameInProgress:
		g.spinner.SetVisible(false)
		currEngine := g.getCurrentEngine()
		if currEngine.RequiresHumanInput() {
			g.state = WaitingForHuman
		} else {
			g.state = ComputerThinking
			// Launch the goroutine to generate the move;
			// not called in ComputerThinking state to avoid
			// multiple goroutines being launched as Update is called repeatedly.
			go func() {
				currEngine := g.getCurrentEngine()
				move := currEngine.GenMove(g.gameBoard)
				g.boardWidget.DoComputerMove(move, g.gameBoard.playerToMove)
				g.state = AnimatingComputerMove
			}()
		}

	case WaitingForHuman:
		g.boardWidget.Update(g.gameBoard)
		if move, ok := g.boardWidget.GetAndClearHumanMove(); ok {
			if valid, _ := move.IsValid(g.gameBoard); valid {
				g.gameBoard.Move(move)
			}
			g.state = GameInProgress
		}

	case ComputerThinking:
		g.spinner.SetVisible(true)

	case AnimatingComputerMove:
		g.boardWidget.UpdateComputerDragInfo()
		if move, ok := g.boardWidget.GetAndClearComputerMove(); ok {
			g.gameBoard.Move(move)
			g.state = GameInProgress
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.boardWidget.Draw(screen, g.gameBoard)
	if g.spinner.IsVisible() {
		g.spinner.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
