package main

import (
	"gocycle/core"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Arena      core.Arena
	ArenaTimer *Timer
}

func NewGame() *Game {
	player1 := core.NewPlayer(1, core.Vector{X: 10, Y: 10}, core.Right, &core.RandomController{})
	player2 := core.NewPlayer(2, core.Vector{X: 30, Y: 30}, core.Left, &core.RandomController{})
	player3 := core.NewPlayer(3, core.Vector{X: 10, Y: 30}, core.Down, &core.RandomController{})
	player4 := core.NewPlayer(4, core.Vector{X: 30, Y: 10}, core.Up, &core.RandomController{})
	return &Game{
		Arena:      *core.NewArena(ArenaWidth, ArenaHeight, []*core.Player{player1, player2, player3, player4}),
		ArenaTimer: NewTimer(1000 * time.Millisecond),
	}
}

func (g *Game) Update() error {
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
