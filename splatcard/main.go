package main

import (
	"image/color"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ScreenWidth  = 576
	ScreenHeight = 360

	// Layout
	LetterWidth = 28
	PlatformY   = 250
	TileStartX  = 40
)

// Game is the main game struct.
type Game struct {
	CardSet *CardSet
	Card    *Card
	Frog    *Frog

	// Entities for current word
	Platforms      []*Platform
	currentAnswer  string // The characters the player has typed so far
	currentIndex   int    // The index of the next character to be typed
	backspaceCount int    // New field
}

// toLower converts a rune to its lowercase equivalent.
func toLower(r rune) rune {
	return unicode.ToLower(r)
}

func (g *Game) Update() error {
	// Handle regular character input
	var chars []rune
	chars = ebiten.AppendInputChars(chars)
	for _, r := range chars {
		if g.currentIndex < len(g.Card.Value) {
			expectedChar := rune(g.Card.Value[g.currentIndex])
			if toLower(r) == toLower(expectedChar) {
				// Correct character typed: update answer and advance
				g.currentAnswer += string(r)
				g.currentIndex++
				// Update frog position
				g.Frog.X = g.Platforms[g.currentIndex].X
			}
		}
	}

	// Handle backspace input
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(g.currentAnswer) > 0 {
		// Prevent rapid backspacing
		if g.backspaceCount < 1 { // Use a counter to limit rapid fire
			g.currentAnswer = g.currentAnswer[:len(g.currentAnswer)-1]
			g.currentIndex--
			g.backspaceCount = 10 // A delay of 10 ticks, adjust as needed
			// Update frog position
			g.Frog.X = g.Platforms[g.currentIndex].X
		}
		g.backspaceCount--
	} else {
		g.backspaceCount = 0
	}

	// Check if the answer is complete
	if g.currentIndex == len(g.Card.Value) {
		// The player has completed the word, start a new card
		g.StartNewCard()
	}

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

	for _, platform := range g.Platforms {
		platform.Draw(screen)
	}

	// Draw the letters the player has typed
	for i, ch := range g.currentAnswer {
		drawTextAt(screen, string(ch),
			TileStartX+float64((float64(i)+0.5)*LetterWidth), PlatformY,
			text.AlignCenter)
	}

	g.Frog.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) StartNewCard() {
	g.Card = g.CardSet.GetCard()

	g.Platforms = []*Platform{}
	for i := range len(g.Card.Value) + 1 {
		x := TileStartX + float64(i*LetterWidth)
		y := float64(PlatformY)
		g.Platforms = append(g.Platforms, NewPlatform(x, y))
	}

	g.Frog.X = TileStartX
	g.Frog.Y = float64(PlatformY - 32)

	g.currentAnswer = ""
	g.currentIndex = 0
}

func NewGame() *Game {
	g := Game{
		CardSet: NewCardSet(),
		Frog:    NewFrog(),
		Card:    nil,
	}
	g.StartNewCard()
	return &g
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(2*ScreenWidth, 2*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
