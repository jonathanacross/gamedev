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
	BootOffsetX     = 20
	BootOffsetY     = 20

	// Game mechanics
	HeronSpeed              = 1.2
	BootFallSpeed           = 1.0
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
	Platforms []*Platform
	Herons    []*Heron
	Boots     []*Boot
	Crocodile *Crocodile

	// Game state fields
	gameState        GameState
	currentAnswer    string
	currentIndex     int
	numMistakes      int
	backspaceTimer   *Timer
	surprisedTimer   *Timer
	flashAnswerTimer *Timer
	heronSpawnTimer  *Timer
}

func (g *Game) updateHerons() {
	// Remove herons that have gone offscreen
	var remainingHerons []*Heron
	for _, h := range g.Herons {
		if !h.IsOffscreen() {
			remainingHerons = append(remainingHerons, h)
		}
	}
	g.Herons = remainingHerons

	// Spawn new herons if needed
	if g.heronSpawnTimer.IsReady() {
		// pick a random platform, but not the ending spot
		platform := rand.Intn(len(g.Platforms) - 1)
		heightOffset := rand.Float64() * 60
		heron := NewHeron(ScreenWidth, FallingItemTopY+heightOffset, g.Platforms[platform].X)
		g.Herons = append(g.Herons, heron)
		g.Boots = append(g.Boots, heron.GetBoot())
		// Reset the timer with a new random duration between 0.5 and 1.5 seconds.
		g.heronSpawnTimer = NewTimer(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	}

	// process the current list of herons
	for _, h := range g.Herons {
		h.Update()
	}
}

func (g *Game) updateBoots() {
	for _, b := range g.Boots {
		b.Update()
	}
	g.removeFallenBoots()
}

func (g *Game) removeFallenBoots() {
	var remainingBoots []*Boot
	for _, b := range g.Boots {
		b.Update()
		if b.Y < PlatformY {
			remainingBoots = append(remainingBoots, b)
		}
	}
	g.Boots = remainingBoots
}

func (g *Game) Update() error {
	PlayMusic()

	// Update all components
	g.backspaceTimer.Update()
	g.surprisedTimer.Update()
	g.heronSpawnTimer.Update()

	switch g.gameState {
	case Playing:
		g.Frog.Update()
		for _, h := range g.Herons {
			h.Update()
		}

		g.updateHerons()
		g.updateBoots()
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
		g.StartNewCard()
	}

	if g.Frog.state == Dying && g.Frog.IsDyingFinished() {
		g.resetCurrentWord()
	}
}

func (g *Game) handleInput() {
	// Handle backspace input
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && g.backspaceTimer.IsReady() && len(g.currentAnswer) > 0 && g.Frog.state != Jumping {
		g.currentAnswer = g.currentAnswer[:len(g.currentAnswer)-1]
		g.currentIndex--
		// Instantly move the frog back
		g.Frog.X = g.Platforms[g.currentIndex].X
		g.backspaceTimer.Reset()
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

				// Reset the frog's state to Idle to allow the jump
				g.Frog.state = Idle

				// Move the frog to the next platform.
				// This handles both intermediate and final letters.
				targetX := g.Platforms[g.currentIndex].X
				g.Frog.X = targetX

				// Check if the word is complete and trigger the jump
				if g.currentIndex == len(g.Card.Value) {
					g.Frog.Jump(targetX)
					PlaySound(ClearSoundBytes)
				}
			} else {
				g.Frog.state = Surprised
				g.surprisedTimer.Reset()
				PlaySound(ErrorSoundBytes)
				g.numMistakes++
				g.Card.ConsecutiveCorrect = 0
				if g.numMistakes <= NumMistakesForCrocodile {
					g.Crocodile.Y -= CrocodileUp
				}
				if g.numMistakes >= NumMistakesForCrocodile {
					g.Crocodile.Bite()
				}
			}
		}
	}
}

func (g *Game) checkCollisions() {
	if g.Frog.state == Dying {
		return // Do not check for collisions if the frog is already dying
	}

	// Check for collisions with boots
	for _, b := range g.Boots {
		if g.Frog.HasCollided(&b.BaseSprite) {
			PlaySound(SplatSoundBytes)
			g.Frog.Hit()
			return
		}
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
	g.Herons = []*Heron{}
	g.Boots = []*Boot{}
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

	for _, h := range g.Herons {
		h.Draw(screen)
	}

	for _, b := range g.Boots {
		b.Draw(screen)
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
	g.Herons = []*Heron{}
	g.Boots = []*Boot{}
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
		heronSpawnTimer:  NewTimer(2 * time.Second),
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
