package main

import (
	"gocycle/core"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	GameMode GameState = iota
	CharacterDisplay
)

type Game struct {
	Arena           core.Arena
	ArenaTimer      *Timer
	HumanController *core.HumanController
	State           GameState
	Selector        *CharacterSelector
}

func NewGame() *Game {
	humanController := core.NewHumanController()
	//turnProb := 0.1

	player1 := core.NewPlayer(1, core.Vector{X: 10, Y: 10}, core.Right, humanController)
	player2 := core.NewPlayer(2, core.Vector{X: 30, Y: 30}, core.Left, &core.RandomTurnerController{TurnProb: 0.01})
	player3 := core.NewPlayer(3, core.Vector{X: 10, Y: 30}, core.Down, &core.RandomTurnerController{TurnProb: 0.4})
	player4 := core.NewPlayer(4, core.Vector{X: 30, Y: 10}, core.Up, &core.AreaController{})
	players := []*core.Player{player1, player2, player3, player4}

	selector := NewCharacterSelector(
		16, 30, 74, 90, 2, 5, 10, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	return &Game{
		Arena:           *core.NewArena(ArenaWidth, ArenaHeight, players),
		ArenaTimer:      NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		HumanController: humanController,
		State:           CharacterDisplay,
		Selector:        selector,
	}
}

func (g *Game) UpdateInGameMode() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.HumanController.EnqueueDirection(core.Left)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.HumanController.EnqueueDirection(core.Right)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.HumanController.EnqueueDirection(core.Up)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.HumanController.EnqueueDirection(core.Down)
	}

	g.ArenaTimer.Update()
	if g.ArenaTimer.IsReady() {
		g.Arena.Update()
		g.ArenaTimer.Reset()
	}
	return nil
}

func (g *Game) UpdateInCharacterDisplayMode() error {
	g.Selector.Update()
	return nil
}

func (g *Game) Update() error {
	switch g.State {
	case GameMode:
		return g.UpdateInGameMode()
	case CharacterDisplay:
		return g.UpdateInCharacterDisplayMode()
	}
	return nil
}

func (g *Game) DrawInGameMode(screen *ebiten.Image) {
	g.DrawArena(screen)
}

func (g *Game) DrawInCharacterDisplayMode(screen *ebiten.Image) {
	g.Selector.Draw(screen)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case GameMode:
		g.DrawInGameMode(screen)
	case CharacterDisplay:
		g.DrawInCharacterDisplayMode(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
