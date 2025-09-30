package main

import (
	"gocycle/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// PlayerColorMap maps Player IDs (1, 2, 3, 4) to their specific rendering colors.
// Index 0 is left empty since Player IDs start at 1.
var PlayerColorMap = []struct {
	Head color.RGBA
	Path color.RGBA
}{
	// Index 0: Unused
	{},
	// Player 1
	{
		Head: color.RGBA{R: 219, G: 65, B: 97, A: 255},
		Path: color.RGBA{R: 178, G: 16, B: 48, A: 255},
	},
	// Player 2
	{
		Head: color.RGBA{R: 113, G: 243, B: 65, A: 255},
		Path: color.RGBA{R: 73, G: 170, B: 16, A: 255},
	},
	// Player 3
	{
		Head: color.RGBA{R: 195, G: 178, B: 255, A: 255},
		Path: color.RGBA{R: 146, G: 65, B: 243, A: 255},
	},
	// Player 4
	{
		Head: color.RGBA{R: 255, G: 243, B: 146, A: 255},
		Path: color.RGBA{R: 235, G: 211, B: 32, A: 255},
	},
}

// GetSquareColor returns the appropriate color.RGBA for a square at (x, y)
// based on its state and the current player positions.
func (g *Game) GetSquareColor(x, y int, square core.Square) color.RGBA {

	// If ANY player's current position matches (x, y), draw the Head color.
	currentPos := core.Vector{X: x, Y: y}
	for _, player := range g.Arena.Players {
		if player.IsAlive && player.Position.Equals(currentPos) {
			// Safety check: ensure the ID is within the bounds of our color map
			if player.ID >= 1 && player.ID < len(PlayerColorMap) {
				return PlayerColorMap[player.ID].Head
			}
		}
	}

	switch square {
	case core.Open:
		// Background color for open space
		return color.RGBA{R: 34, G: 32, B: 52, A: 255}
	case core.Wall:
		// Color for the arena border walls
		return color.RGBA{R: 63, G: 80, B: 151, A: 255}
	default:
		// --- Player Path / Head Logic ---
		playerID := int(square)

		// Safety check to ensure the ID is within the bounds of our color map
		if playerID < 1 || playerID >= len(PlayerColorMap) {
			// Use a fallback color for unexpected IDs
			return color.RGBA{R: 255, G: 0, B: 255, A: 255}
		}

		renderData := PlayerColorMap[playerID]

		// Check if the player is alive and the square is their current "head" position.
		// Note: We access players using ID-1 because the Players slice is 0-indexed.
		if playerID <= len(g.Arena.Players) {
			player := g.Arena.Players[playerID-1]

			// Check if the current grid coordinate matches the player's head position.
			currentPos := core.Vector{X: x, Y: y}

			// Use player.IsAlive check to ensure dead players don't have a "Head"
			if player.IsAlive && player.Position.Equals(currentPos) {
				return renderData.Head
			}
		}

		// If it's not the head (or the player is dead), return the path color.
		return renderData.Path
	}
}

func (g *Game) DrawArena(screen *ebiten.Image) {
	for y := 0; y < g.Arena.Height; y++ {
		for x := 0; x < g.Arena.Width; x++ {
			square := g.Arena.Grid[y][x]

			color := g.GetSquareColor(x, y, square)

			r32, g32, b32, _ := color.RGBA()
			r := float64(r32) / 0xFFFF
			g := float64(g32) / 0xFFFF
			b := float64(b32) / 0xFFFF

			var cm colorm.ColorM
			cm.Scale(r, g, b, 1.0)

			op := &colorm.DrawImageOptions{}
			op.GeoM.Translate(float64(x*SquareSize+10), float64(y*SquareSize+10))

			colorm.DrawImage(screen, SquareImage, cm, op)
		}
	}
}
