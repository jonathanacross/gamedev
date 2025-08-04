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
		blackEngine: &GreedyEngine{}, // Use GreedyEngine for the computer player
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
		}

	case WaitingForHuman:
		g.boardWidget.Update(g.gameBoard)
		if move, ok := g.boardWidget.GetAndClearHumanMove(); ok {
			if valid, _ := IsValidMove(g.gameBoard, move); valid {
				g.gameBoard.Move(move)
			}
			g.state = GameInProgress
		}

	case ComputerThinking:
		g.spinner.SetVisible(true)
		go func() {
			currEngine := g.getCurrentEngine()
			move := currEngine.GenMove(g.gameBoard)
			g.boardWidget.DoComputerMove(move, g.gameBoard.playerToMove)
			g.state = AnimatingComputerMove
		}()

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
