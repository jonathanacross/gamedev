package main

import (
	"gocycle/core"
	"image/color"
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
	// TODO: move this inside the character selector state
	Selector *CharacterSelector
}

func NewGame() *Game {
	//humanController := core.NewHumanController()
	//turnProb := 0.1

	//player1 := core.NewPlayer(1, core.Vector{X: 10, Y: 10}, core.Right, humanController)
	//player2 := core.NewPlayer(2, core.Vector{X: 30, Y: 30}, core.Left, &core.RandomTurnerController{TurnProb: 0.01})
	//player3 := core.NewPlayer(3, core.Vector{X: 10, Y: 30}, core.Down, &core.RandomTurnerController{TurnProb: 0.4})
	//player4 := core.NewPlayer(4, core.Vector{X: 30, Y: 10}, core.Up, &core.AreaController{})
	//players := []*core.Player{player1, player2, player3, player4}

	selector := NewCharacterSelector(
		16, 30, 74, 90, 2, 5, 10, []int{6, 8, 0, 1, 2, 7, 9, 3, 4, 5})

	return &Game{
		//Arena:           *core.NewArena(ArenaWidth, ArenaHeight, players),
		//ArenaTimer:      NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		//HumanController: humanController,
		State:    &TitleScreenState{},
		Selector: selector,
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
		g.State = &CharacterPickerState{}
	}
	return nil
}

func (gs *TitleScreenState) Draw(g *Game, screen *ebiten.Image) {
	drawTextAt(screen, "GoCycle", ScreenWidth/2, ScreenHeight/2, text.AlignCenter, color.White)
	drawTextAt(screen, "Press Space", ScreenWidth/2, 3*ScreenHeight/4, text.AlignCenter, color.White)
}

// ------------------- Character Picker State

type CharacterPickerState struct{}

func (gs *CharacterPickerState) Update(g *Game) error {
	g.Selector.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State = NewGamePlayState(g.Selector.GetSelectedCharacters())
	}

	return nil
}

func (gs *CharacterPickerState) Draw(g *Game, screen *ebiten.Image) {
	g.Selector.Draw(screen)
}

// ------------------- Game Play State

type GamePlayState struct {
	ArenaView        *ArenaView
	ArenaTimer       *Timer
	HumanController1 *core.HumanController
	HumanController2 *core.HumanController
}

func NewGamePlayState(characters []*CharData) *GamePlayState {
	// humanController1 := core.NewHumanController()
	// humanController2 := core.NewHumanController()
	//turnProb := 0.1

	// TODO: these players should be passed in to constructor
	// player1 := core.NewPlayer(1, core.Vector{X: 10, Y: 10}, core.Right, humanController1)
	// player2 := core.NewPlayer(2, core.Vector{X: 30, Y: 30}, core.Left, &core.RandomTurnerController{TurnProb: 0.01})
	// player3 := core.NewPlayer(3, core.Vector{X: 10, Y: 30}, core.Down, &core.RandomTurnerController{TurnProb: 0.4})
	// player4 := core.NewPlayer(4, core.Vector{X: 30, Y: 10}, core.Up, &core.AreaController{})
	// players := []*core.Player{player1, player2, player3, player4}

	var human1 *core.HumanController
	var human2 *core.HumanController
	for _, char := range characters {
		if char.ControllerType == HumanFirstPlayer {
			human1 = char.Controller.(*core.HumanController)
		} else if char.ControllerType == HumanSecondPlayer {
			human2 = char.Controller.(*core.HumanController)
		}
	}

	return &GamePlayState{
		ArenaView:        NewArenaView(ArenaWidth, ArenaHeight, characters),
		ArenaTimer:       NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		HumanController1: human1,
		HumanController2: human2,
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
}
