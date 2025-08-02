package main

import (
	_ "image/png" // Import for image decoding
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	boardWidget *BoardWidget
	moveTimer   *Timer
	engine      Engine
}

func NewGame() *Game {
	return &Game{
		boardWidget: NewBoardWidget(),
		moveTimer:   NewTimer(2 * time.Second),
		engine:      &RandomEngine{},
	}
}

func (g *Game) Update() error {
	g.boardWidget.Update()
	// g.moveTimer.Update()
	// if g.moveTimer.IsReady() {
	// 	move := g.engine.GenMove(g.gameBoard)
	// 	g.gameBoard.Move(move)
	// 	g.moveTimer.Reset()
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.boardWidget.Draw(screen)
}

// Layout returns the game's logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
