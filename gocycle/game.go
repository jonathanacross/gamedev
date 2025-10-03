package main

import (
	"gocycle/core"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Arena           core.Arena
	ArenaTimer      *Timer
	HumanController *core.HumanController
}

func NewGame() *Game {
	humanController := core.NewHumanController()
	//turnProb := 0.1

	player1 := core.NewPlayer(1, core.Vector{X: 10, Y: 10}, core.Right, humanController)
	player2 := core.NewPlayer(2, core.Vector{X: 30, Y: 30}, core.Left, &core.RandomTurnerController{TurnProb: 0.01})
	player3 := core.NewPlayer(3, core.Vector{X: 10, Y: 30}, core.Down, &core.RandomAvoidingController{})
	player4 := core.NewPlayer(4, core.Vector{X: 30, Y: 10}, core.Up, &core.AreaController{})
	players := []*core.Player{player1, player2, player3, player4}

	return &Game{
		Arena:           *core.NewArena(ArenaWidth, ArenaHeight, players),
		ArenaTimer:      NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		HumanController: humanController,
	}
}

func (g *Game) Update() error {

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

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawArena(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
