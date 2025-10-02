package main

import (
	"gocycle/core"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Arena           core.Arena
	ArenaTimer      *Timer
	HumanController *core.HumanController
}

func NewGame() *Game {
	humanController := core.NewHumanController(core.Right)
	turnProb := 0.1

	player1 := core.NewPlayer(1, core.Vector{X: 10, Y: 10}, core.Right, humanController)
	player2 := core.NewPlayer(2, core.Vector{X: 30, Y: 30}, core.Left, &core.RandomTurnerController{TurnProb: 0.01})
	player3 := core.NewPlayer(3, core.Vector{X: 10, Y: 30}, core.Down, &core.RandomTurnerController{TurnProb: turnProb})
	player4 := core.NewPlayer(4, core.Vector{X: 30, Y: 10}, core.Up, &core.AreaController{})
	return &Game{
		Arena:           *core.NewArena(ArenaWidth, ArenaHeight, []*core.Player{player1, player2, player3, player4}),
		ArenaTimer:      NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		HumanController: humanController,
	}
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.HumanController.RequestedDirection = core.Left
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.HumanController.RequestedDirection = core.Right
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.HumanController.RequestedDirection = core.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.HumanController.RequestedDirection = core.Down
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
