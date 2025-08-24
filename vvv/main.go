package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	player  Player
	level   Level
	gravity float64
}

const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

func (g *Game) Update() error {
	if g.player.IsOnGround() && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.gravity *= -1
	}
	g.level.Update()
	g.player.Update(g)

	PlayMusic()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.Draw(screen)
	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	spriteSheet := NewSpriteSheet(TileSet, TileSize, TileSize, 5, 7)
	g := &Game{
		player:  *NewPlayer(),
		level:   *NewLevel(LevelTilemapJSON, spriteSheet),
		gravity: 0.5,
	}
	ebiten.SetWindowSize(3*ScreenWidth, 3*ScreenHeight)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
