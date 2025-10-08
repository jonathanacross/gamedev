package main

import (
	"gocycle/core"
	"image/color"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameState interface {
	Update(g *Game) error
	Draw(g *Game, screen *ebiten.Image)
}

type Game struct {
	State GameState
}

func NewGame() *Game {
	return &Game{
		State: &TitleScreenState{},
	}
}

func (g *Game) Update() error {
	g.State.Update(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.State.Draw(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func drawTextAt(screen *ebiten.Image, message string, x float64, y float64, align text.Align, c color.Color) {
	fontSize := float64(8)
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

// ------------------- Title Screen State

type TitleScreenState struct{}

func (gs *TitleScreenState) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State = NewCharacterPickerState()
	}
	return nil
}

func (gs *TitleScreenState) Draw(g *Game, screen *ebiten.Image) {
	drawTextAt(screen, "GoCycle", ScreenWidth/2, ScreenHeight/2, text.AlignCenter, color.White)
	drawTextAt(screen, "Press Space", ScreenWidth/2, 3*ScreenHeight/4, text.AlignCenter, color.White)
}

// ------------------- Character Picker State

type CharacterPickerState struct {
	Selector *CharacterSelector
}

func NewCharacterPickerState() *CharacterPickerState {
	return &CharacterPickerState{
		Selector: NewCharacterSelector(
			16, 30, 74, 90, 2, 5, 10, []int{6, 8, 0, 1, 2, 7, 9, 3, 4, 5}),
	}
}

func (gs *CharacterPickerState) Update(g *Game) error {
	gs.Selector.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		selectedChars := gs.Selector.GetSelectedCharacters()
		// Shuffle so that chars don't always start in the same place.
		rand.Shuffle(len(selectedChars), func(i, j int) {
			selectedChars[i], selectedChars[j] = selectedChars[j], selectedChars[i]
		})

		g.State = NewGamePlayState(selectedChars)
	}

	return nil
}

func (gs *CharacterPickerState) Draw(g *Game, screen *ebiten.Image) {
	gs.Selector.Draw(screen)
}

// ------------------- Game Play State

type GamePlayState struct {
	ArenaView        *ArenaView
	ArenaTimer       *Timer
	HumanController1 *core.HumanController
	HumanController2 *core.HumanController
	CharacterCards   []*CharacterFrame
}

func NewGamePlayState(characters []*CharData) *GamePlayState {
	var human1 *core.HumanController
	var human2 *core.HumanController
	for _, char := range characters {
		switch char.ControllerType {
		case HumanFirstPlayer:
			human1 = char.Controller.(*core.HumanController)
		case HumanSecondPlayer:
			human2 = char.Controller.(*core.HumanController)
		}
	}

	positionData := PositionDataByNumChars[len(characters)]

	cards := []*CharacterFrame{}
	for i, char := range characters {
		cards = append(cards, NewCharacterFrame(char,
			positionData[i].CardX, positionData[i].CardY, CharacterNeutral, false))
	}

	players := []*core.Player{}
	initialDirections := []core.Vector{core.Right, core.Left, core.Down, core.Up}
	for i, char := range characters {
		players = append(players, core.NewPlayer(i+1,
			positionData[i].ArenaLoc, initialDirections[i], char.Controller))
	}
	var arena = core.NewArena(ArenaWidth, ArenaHeight, players)

	return &GamePlayState{
		ArenaView:        NewArenaView(arena, characters),
		ArenaTimer:       NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		HumanController1: human1,
		HumanController2: human2,
		CharacterCards:   cards,
	}
}

func (gs *GamePlayState) Update(g *Game) error {
	if gs.HumanController1 != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			gs.HumanController1.EnqueueDirection(core.Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			gs.HumanController1.EnqueueDirection(core.Right)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			gs.HumanController1.EnqueueDirection(core.Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			gs.HumanController1.EnqueueDirection(core.Down)
		}
	}
	if gs.HumanController2 != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			gs.HumanController2.EnqueueDirection(core.Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			gs.HumanController2.EnqueueDirection(core.Right)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			gs.HumanController2.EnqueueDirection(core.Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			gs.HumanController2.EnqueueDirection(core.Down)
		}
	}

	gs.ArenaTimer.Update()
	if gs.ArenaTimer.IsReady() {
		gs.ArenaView.Update()
		gs.ArenaTimer.Reset()
	}
	return nil
}

func (gs *GamePlayState) Draw(g *Game, screen *ebiten.Image) {
	gs.ArenaView.Draw(screen)
	for _, card := range gs.CharacterCards {
		card.Draw(screen)
	}
}

var PositionDataByNumChars = getPositionData()

type PositionData struct {
	ArenaLoc core.Vector
	CardX    float64
	CardY    float64
}

// Positions of where to put players in the arena and where to draw the player
// cards.  This depends on the number of players.
func getPositionData() [][]PositionData {
	return [][]PositionData{
		{},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
		},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: 30, Y: 30}, CardX: 300, CardY: 10},
		},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: 30, Y: 30}, CardX: 300, CardY: 10},
			{ArenaLoc: core.Vector{X: 10, Y: 30}, CardX: 10, CardY: 120},
		},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: 30, Y: 30}, CardX: 300, CardY: 120},
			{ArenaLoc: core.Vector{X: 10, Y: 30}, CardX: 10, CardY: 120},
			{ArenaLoc: core.Vector{X: 30, Y: 10}, CardX: 300, CardY: 10},
		},
	}
}
