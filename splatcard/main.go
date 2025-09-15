package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 240
)

// Game is the main game struct.
type Game struct {
	CardSet *CardSet
	Card    *Card
}

func (g *Game) Update() error {
	// PlayMusic()

	return nil
}

// drawTextAt is a helper function to draw text on the screen with alignment.
func drawTextAt(screen *ebiten.Image, message string, x float64, y float64, align text.Align) {
	fontSize := float64(16)
	fontFace := &text.GoTextFace{
		Source: MainFaceSource,
		Size:   fontSize,
	}

	// Manually handle alignment to ensure pixel-perfect rendering
	textWidth, _ := text.Measure(message, fontFace, 1.0)
	if align == text.AlignCenter {
		x -= float64(textWidth) / 2
	} else if align == text.AlignEnd {
		x -= float64(textWidth)
	}
	x = float64(int(x))
	y = float64(int(y))

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignStart

	text.Draw(screen, message, fontFace, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})

	drawTextAt(screen, g.Card.Key, ScreenWidth/2, 50, text.AlignCenter)

	for i, ch := range g.Card.Value {
		drawTextAt(screen, string(ch), 50+float64(i*16), 150, text.AlignCenter)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	g := Game{
		CardSet: NewCardSet(),
		Card:    nil,
	}
	g.Card = g.CardSet.GetCard()
	return &g
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
