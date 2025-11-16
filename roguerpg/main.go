package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 240

	TileSize = 16

	ShowDebugInfo = false
)

type Game struct {
	level  *Level
	player *Player
	camera *Camera
}

func NewGame() *Game {
	level := BuildLevel(70, 50)
	player := NewPlayer()
	player.Location = level.FindRandomFloorLocation()
	return &Game{
		level:  level,
		player: player,
		camera: NewCamera(ScreenWidth, ScreenHeight),
	}
}

func (g *Game) Update() error {
	g.player.HandleUserInput()
	g.player.Update(g.level)
	g.camera.CenterOn(g.player.Location)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cameraMatrix := g.camera.WorldToScreen()
	viewRect := g.camera.GetViewRect()

	for _, row := range g.level.Tiles {
		for _, tile := range row {
			if tile.HitBox().Intersects(viewRect) {
				tile.Draw(screen, cameraMatrix)
				if tile.solid {
					tile.DrawDebugInfo(screen, cameraMatrix)
				}
			}
		}
	}
	g.player.Draw(screen, cameraMatrix)
	g.player.DrawDebugInfo(screen, cameraMatrix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(ScreenWidth*3, ScreenHeight*3)
	ebiten.SetWindowTitle("Rogue RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
