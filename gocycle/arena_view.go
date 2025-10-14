package main

import (
	"gocycle/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type ArenaView struct {
	Arena      *core.Arena
	Characters []*CharData
}

func NewArenaView(arena *core.Arena, characters []*CharData) *ArenaView {
	return &ArenaView{
		Arena:      arena,
		Characters: characters,
	}
}

// GetSquareColor returns the appropriate color.Color for a square at (x, y)
// based on its state and the current player positions.
func (av *ArenaView) GetSquareColor(x, y int, square core.Square) color.Color {
	currentPos := core.Vector{X: x, Y: y}

	// Check for Player Head Position.
	// If any player's current position matches (x, y), draw their Head color.
	for _, player := range av.Arena.Players {
		if player.IsAlive && player.Position.Equals(currentPos) {
			// characters is 0-indexed, player.ID is 1-indexed
			charData := av.Characters[player.ID-1]
			return charData.BrightColor
		}
	}

	// Check for Arena Square Type
	switch square {
	case core.Open:
		return color.RGBA{R: 34, G: 32, B: 52, A: 255}
	case core.Wall:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255}
	default:
		// Player Path / Trail
		playerID := int(square)
		// characters is 0-indexed, player.ID is 1-indexed
		charData := av.Characters[playerID-1]
		return charData.DarkColor
	}
}

func (av *ArenaView) Draw(screen *ebiten.Image) {
	for y := 0; y < av.Arena.Height; y++ {
		for x := 0; x < av.Arena.Width; x++ {
			square := av.Arena.Grid[y][x]

			color := av.GetSquareColor(x, y, square)

			r32, g32, b32, _ := color.RGBA()
			r := float64(r32) / 0xFFFF
			g := float64(g32) / 0xFFFF
			b := float64(b32) / 0xFFFF

			var cm colorm.ColorM
			cm.Scale(r, g, b, 1.0)

			op := &colorm.DrawImageOptions{}
			op.GeoM.Translate(float64(x*SquareSize+ArenaOffsetX), float64(y*SquareSize+ArenaOffsetY))

			colorm.DrawImage(screen, SquareImage, cm, op)
		}
	}
}

func (av *ArenaView) Update() {
	av.Arena.Update()
}
