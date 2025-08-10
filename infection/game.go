package main

import (
	"image/color"

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
	state         GameState
	gameBoard     *Board
	boardWidget   *BoardWidget
	spinner       *Spinner
	newGameButton *Button

	player1View *PlayerView
	player2View *PlayerView
}

func NewGame() *Game {
	g := Game{
		state:         GameInProgress,
		gameBoard:     NewBoard(),
		boardWidget:   NewBoardWidget(30, 130),
		spinner:       NewSpinner(255, 628, 0.03),
		newGameButton: nil,
		player1View:   NewPlayerView(30, 30, 0, color.RGBA{255, 255, 0, 255}),
		player2View:   NewPlayerView(300, 30, 0, color.RGBA{255, 0, 0, 255}),
	}

	buttonWidth := 150
	buttonHeight := 40
	g.newGameButton = NewButton(30, 608, buttonWidth, buttonHeight,
		GenerateButtonImage(buttonWidth, buttonHeight, "New Game", color.RGBA{55, 148, 110, 255}, color.RGBA{0, 0, 0, 255}),
		GenerateButtonImage(buttonWidth, buttonHeight, "New Game", color.RGBA{153, 229, 80, 255}, color.RGBA{0, 0, 0, 255}),
		func() { g.gameBoard = NewBoard() },
	)
	return &g
}

func (g *Game) getCurrentEngine() Engine {
	if g.gameBoard.playerToMove == White {
		return g.player1View.GetEngine()
	}
	return g.player2View.GetEngine()
}

func (g *Game) Update() error {
	g.spinner.Update()
	g.player1View.Update()
	g.player2View.Update()
	g.newGameButton.Update()

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
		if !g.getCurrentEngine().RequiresHumanInput() {
			// player switched the engine during their turn
			g.state = GameInProgress
		} else if move, ok := g.boardWidget.GetAndClearHumanMove(); ok {
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
	w, b := g.gameBoard.Score()
	screen.Fill(color.RGBA{10, 40, 30, 255})
	g.player1View.Draw(screen, w, g.gameBoard.playerToMove == White)
	g.player2View.Draw(screen, b, g.gameBoard.playerToMove == Black)
	g.boardWidget.Draw(screen, g.gameBoard)
	if g.spinner.IsVisible() {
		g.spinner.Draw(screen)
	}
	g.newGameButton.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
