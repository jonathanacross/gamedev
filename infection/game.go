package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameState int

// Layout
const (
	ScreenWidth  = 508
	ScreenHeight = 680
	TileSize     = 64

	Margin = 30

	Player1ViewX = Margin
	Player1ViewY = Margin
	Player2ViewX = 300
	Player2ViewY = Margin

	BoardWidgetX = Margin
	BoardWidgetY = 135

	NewGameButtonX      = Margin
	NewGameButtonY      = 608
	NewGameButtonWidth  = 150
	NewGameButtonHeight = 40

	SpinnerSize = 48
	SpinnerX    = (ScreenWidth) / 2
	SpinnerY    = 628
)

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
		boardWidget:   NewBoardWidget(BoardWidgetX, BoardWidgetY),
		spinner:       NewSpinner(SpinnerX, SpinnerY, 0.03),
		newGameButton: nil,
		player1View:   NewPlayerView(Player1ViewX, Player1ViewY, 0, color.RGBA{255, 255, 0, 255}),
		player2View:   NewPlayerView(Player2ViewX, Player2ViewY, 0, color.RGBA{255, 0, 0, 255}),
	}

	g.newGameButton = NewButton(NewGameButtonX, NewGameButtonY, NewGameButtonWidth, NewGameButtonHeight,
		GenerateButtonImage(NewGameButtonWidth, NewGameButtonHeight, "New Game", color.RGBA{55, 148, 110, 255}, color.RGBA{0, 0, 0, 255}),
		GenerateButtonImage(NewGameButtonWidth, NewGameButtonHeight, "New Game", color.RGBA{153, 229, 80, 255}, color.RGBA{0, 0, 0, 255}),
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
		if g.gameBoard.IsGameOver() {
			g.state = GameOver
		} else {
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

	case GameOver:
		if !g.gameBoard.IsGameOver() {
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
