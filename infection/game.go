package main

import (
	_ "image/png" // Import for image decoding

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	gameBoard *Board
}

func NewGame() *Game {
	return &Game{
		gameBoard: NewBoard(),
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	margin := 40
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize; c++ {
			idx := GetIndex(r, c)
			x := float64(c*TileSize + margin)
			y := float64(r*TileSize + margin)

			backgroundImge := Empty1Square
			if (r+c)%2 == 0 {
				backgroundImge = Empty2Square
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			screen.DrawImage(backgroundImge, op)

			if g.gameBoard.white.Get(idx) {
				screen.DrawImage(WhiteSquare, op)
			} else if g.gameBoard.black.Get(idx) {
				screen.DrawImage(BlackSquare, op)
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
}

// Layout returns the game's logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
