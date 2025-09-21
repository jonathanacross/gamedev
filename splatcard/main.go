package main

import (
	"image/color"
	"math/rand"
	"time"

	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ScreenWidth  = 576
	ScreenHeight = 360

	// Layout
	LetterWidth     = 28
	FallingItemTopY = 40
	PlatformY       = 270
	TileStartX      = 40
	FrogOffsetY     = 20

	// Game mechanics
	FallUpVelocity          = -0.5
	FallDownVelocity        = 1.0
	CrocodileSpeed          = 1.0
	CrocodileUp             = 9
	CrocodileOffsetY        = 33
	NumMistakesForCrocodile = 10
)

// GameState represents the current state of the game.
type GameState int

const (
	Playing GameState = iota
	FlashAnswer
)

// Game is the main game struct.
type Game struct {
	CardSet *CardSet
	Card    *Card
	Frog    *Frog

	// Entities for current word
	Platforms   []*Platform
	Boots       []*Boot
	ThreatHeron *Heron
	Crocodile   *Crocodile

	// Game state fields
	gameState        GameState
	currentAnswer    string
	currentIndex     int
	numMistakes      int
	backspaceTimer   *Timer
	surprisedTimer   *Timer
	flashAnswerTimer *Timer
	typingTimer      *Timer // New timer for typing inactivity
}

func (g *Game) Update() error {
	PlayMusic()

	// Update all components
	g.backspaceTimer.Update()
	g.surprisedTimer.Update()
	g.typingTimer.Update() // Update the new timer

	switch g.gameState {
	case Playing:
		g.Frog.Update()
		for _, boot := range g.Boots {
			boot.Update()
		}
		// Update the flying heron threat
		if g.ThreatHeron != nil {
			g.ThreatHeron.Update()
		}

		g.Crocodile.Update()

		g.handleFrogState()
		g.checkCollisions()
		g.handleInput()

	case FlashAnswer:
		g.flashAnswerTimer.Update()
		if g.flashAnswerTimer.IsReady() {
			g.StartNewCard()
		}
	}

	return nil
}

func (g *Game) handleFrogState() {
	if g.Frog.state == Surprised && g.surprisedTimer.IsReady() {
		g.Frog.state = Idle
	}

	// If the frog has just completed the final jump, start a new card
	if g.Frog.IsJumping() && g.Frog.IsJumpFinished() {
		g.Frog.Land()
		g.Card.ConsecutiveCorrect++
		g.reinsertCard(g.Card, false)
		PlaySound(ClearSoundBytes)
		g.StartNewCard()
	}

	if g.Frog.state == Dying && g.Frog.IsDyingFinished() {
		g.resetCurrentWord()
	}
}

func (g *Game) handleInput() {
	// Check for inactivity to spawn the heron threat
	if g.typingTimer.IsReady() && g.ThreatHeron == nil {
		PlaySound(ErrorSoundBytes) // Play a warning sound
		g.ThreatHeron = NewHeron(g.Frog.X, g.Frog.Y)
		g.typingTimer.Reset()
	}

	// Handle backspace input
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && g.backspaceTimer.IsReady() && len(g.currentAnswer) > 0 && g.Frog.state != Jumping {
		g.currentAnswer = g.currentAnswer[:len(g.currentAnswer)-1]
		g.currentIndex--
		// Instantly move the frog back
		g.Frog.X = g.Platforms[g.currentIndex].X
		g.backspaceTimer.Reset()
		g.typingTimer.Reset() // Reset timer on backspace
	}

	// Handle new character input
	var chars []rune
	chars = ebiten.AppendInputChars(chars)
	for _, r := range chars {
		if g.currentIndex < len(g.Card.Value) {
			expectedChar := rune(g.Card.Value[g.currentIndex])
			if toLower(r) == toLower(expectedChar) {
				g.currentAnswer += string(r)
				g.currentIndex++
				g.typingTimer.Reset() // Reset timer on correct key press

				// Check if this is the final character
				if g.currentIndex == len(g.Card.Value) {
					// Initiate the final jump animation
					targetX := g.Platforms[g.currentIndex].X
					g.Frog.Jump(targetX)
				} else {
					// Instantly move the frog to the next platform
					g.Frog.X = g.Platforms[g.currentIndex].X
				}
			} else {
				g.Frog.state = Surprised
				g.surprisedTimer.Reset()
				PlaySound(ErrorSoundBytes)
				g.numMistakes++
				g.Card.ConsecutiveCorrect = 0
				if g.numMistakes > NumMistakesForCrocodile {
					g.Crocodile.Bite()
				} else {
					g.Crocodile.Y -= CrocodileUp
				}
				g.typingTimer.Reset() // Reset timer on incorrect key press
			}
		}
	}
}

func (g *Game) checkCollisions() {
	if g.Frog.state == Dying {
		return // Do not check for collisions if the frog is already dying
	}

	for _, boot := range g.Boots {
		if g.Frog.HasCollided(&boot.BaseSprite) {
			PlaySound(SplatSoundBytes)
			g.Frog.Hit()
			return
		}
	}

	// Check for collision with the timed heron threat
	if g.ThreatHeron != nil && g.Frog.HasCollided(&g.ThreatHeron.BaseSprite) {
		PlaySound(SplatSoundBytes)
		g.Frog.Hit()
		return
	}

	if g.Crocodile.state == Biting && g.Frog.HasCollided(&g.Crocodile.BaseSprite) {
		PlaySound(MunchSoundBytes)
		g.gameState = FlashAnswer
		g.flashAnswerTimer.Reset()
		g.reinsertCard(g.Card, true)
		return
	}
}

func (g *Game) reinsertCard(card *Card, wasBitten bool) {
	if wasBitten {
		g.Card.ConsecutiveCorrect = 0
		g.CardSet.ReinsertCard(card, rand.Intn(2)+1)
	} else if g.numMistakes > 2 {
		g.CardSet.ReinsertCard(card, 4)
	} else {
		g.CardSet.ReinsertCard(card, 8*g.Card.ConsecutiveCorrect)
	}
}

func (g *Game) resetCurrentWord() {
	g.Frog.state = Idle
	g.Frog.X = g.Platforms[0].X // move frog to first platform
	g.currentAnswer = ""
	g.currentIndex = 0
	g.ThreatHeron = nil // Reset the heron threat
}

// drawTextAt is a helper function to draw text on the screen with alignment.
func drawTextAt(screen *ebiten.Image, message string, x float64, y float64, align text.Align, c color.Color) {
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
	op.ColorScale.ScaleWithColor(c)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignStart

	text.Draw(screen, message, fontFace, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(PondSprite, op)

	for _, platform := range g.Platforms {
		platform.Draw(screen)
	}

	for _, boot := range g.Boots {
		boot.Draw(screen)
	}

	// Draw the timed heron threat if it exists
	if g.ThreatHeron != nil {
		g.ThreatHeron.Draw(screen)
	}

	g.Frog.Draw(screen)
	g.Crocodile.Draw(screen)

	switch g.gameState {
	case Playing:
		drawTextAt(screen, g.Card.Key, ScreenWidth/2, 40, text.AlignCenter, color.Black)
		for i, ch := range g.currentAnswer {
			drawTextAt(screen, string(ch),
				TileStartX+float64((float64(i)+0.5)*LetterWidth), PlatformY-10,
				text.AlignCenter, color.White)
		}
	case FlashAnswer:
		// Flash the correct answer
		c := color.Color(color.White)
		if (g.flashAnswerTimer.currentTicks/10)%2 == 0 {
			c = color.Color(color.RGBA{255, 0, 0, 255})
		}
		for i, ch := range g.Card.Value {
			drawTextAt(screen, string(ch),
				TileStartX+float64((float64(i)+0.5)*LetterWidth), PlatformY-10,
				text.AlignCenter, c)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) StartNewCard() {
	g.Card = g.CardSet.GetCard()
	if g.Card == nil {
		// No more cards, you can handle this as a game end state
		return
	}

	g.Platforms = []*Platform{}
	for i := range len(g.Card.Value) + 1 {
		endPlatform := i == len(g.Card.Value)
		x := TileStartX + float64(i*LetterWidth)
		y := float64(PlatformY)
		g.Platforms = append(g.Platforms, NewPlatform(x, y, endPlatform))
	}

	g.gameState = Playing
	g.Frog.state = Idle
	g.Frog.X = TileStartX
	g.Frog.Y = float64(PlatformY - FrogOffsetY)
	g.Crocodile.state = Floating
	g.Crocodile.X = ScreenWidth
	g.Crocodile.Y = PlatformY - CrocodileOffsetY + CrocodileUp*NumMistakesForCrocodile
	g.numMistakes = 0
	g.currentAnswer = ""
	g.currentIndex = 0

	// Reset heron threat
	g.ThreatHeron = nil
	g.typingTimer.Reset()

	// Create random boots
	g.Boots = []*Boot{}
	// numBoots := rand.Intn(2) + 1 // 1 or 2 boots
	// indices := rand.Perm(len(g.Card.Value) - 1)[0:numBoots]
	// for _, i := range indices {
	// 	x := TileStartX + float64((i+1)*LetterWidth)
	// 	y := float64(rand.Intn(10) - 100)
	// 	g.Boots = append(g.Boots, NewBoot(x, y))
	// }
}

func NewGame() *Game {
	g := Game{
		CardSet:          NewCardSet(),
		Frog:             NewFrog(),
		Crocodile:        NewCrocodile(),
		Card:             nil,
		numMistakes:      0,
		currentAnswer:    "",
		currentIndex:     0,
		backspaceTimer:   NewTimer(100 * time.Millisecond),
		surprisedTimer:   NewTimer(500 * time.Millisecond),
		flashAnswerTimer: NewTimer(2 * time.Second),
		typingTimer:      NewTimer(4 * time.Second), // Initialize new timer
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
