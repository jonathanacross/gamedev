package main

import (
	"image/color"
	"time"
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
	Platforms []*Platform
	// Game state fields
	currentAnswer  string
	currentIndex   int
	backspaceTimer *Timer
	surprisedTimer *Timer
	jumpTargetX    float64
	jumpQueue      []rune
}

func (g *Game) Update() error {
	// Update all timers and the frog's animation on every tick
	g.backspaceTimer.Update()
	g.surprisedTimer.Update()
	g.Frog.Update(g)

	// Check for state transitions and queued jumps
	if g.Frog.state == Surprised && g.surprisedTimer.IsReady() {
		g.Frog.state = Idle
	}

	if g.Frog.state == Jumping && g.Frog.animations[Jumping].IsFinished() {
		g.Frog.state = Idle
		g.Frog.X = g.jumpTargetX
		g.Frog.Y = float64(PlatformY - 32)
	}

	// If the frog has just landed, process the next queued jump
	if g.Frog.state == Idle && len(g.jumpQueue) > 0 {
		g.Frog.state = Jumping
		g.Frog.animations[Jumping].Reset()
		g.Frog.jumpStartX = g.Frog.X
		g.jumpTargetX = g.Platforms[g.currentIndex].X
		g.jumpQueue = g.jumpQueue[1:]
	}

	// Always handle backspace, regardless of frog state
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && g.backspaceTimer.IsReady() && len(g.currentAnswer) > 0 && g.Frog.state != Jumping {
		g.currentAnswer = g.currentAnswer[:len(g.currentAnswer)-1]
		g.currentIndex--
		g.Frog.X = g.Platforms[g.currentIndex].X
		g.backspaceTimer.Reset()
	}

	// Always process new input
	var chars []rune
	chars = ebiten.AppendInputChars(chars)
	for _, r := range chars {
		if g.currentIndex < len(g.Card.Value) {
			expectedChar := rune(g.Card.Value[g.currentIndex])
			if toLower(r) == toLower(expectedChar) {
				// Correct character logic: queue the jump and update the answer.
				g.currentAnswer += string(r)
				g.currentIndex++

				if g.currentIndex < len(g.Platforms) {
					g.jumpQueue = append(g.jumpQueue, r)
				}
			} else {
				// Incorrect character logic: still enter surprised state.
				g.Frog.state = Surprised
				g.surprisedTimer.Reset()
			}
		}
	}

	if g.currentIndex == len(g.Card.Value) {
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
	g.Frog.state = Idle
	g.jumpTargetX = g.Platforms[0].X
	g.jumpQueue = []rune{}
}

func NewGame() *Game {
	g := Game{
		CardSet:        NewCardSet(),
		Frog:           NewFrog(),
		Card:           nil,
		backspaceTimer: NewTimer(100 * time.Millisecond),
		surprisedTimer: NewTimer(500 * time.Millisecond),
	}
	g.StartNewCard()
	return &g
}

func toLower(r rune) rune {
	return unicode.ToLower(r)
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(2*ScreenWidth, 2*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
