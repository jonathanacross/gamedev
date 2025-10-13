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

// GetSquareColor returns the appropriate color.RGBA for a square at (x, y)
// based on its state and the current player positions.
func (av *ArenaView) GetSquareColor(x, y int, square core.Square) color.Color {

	// If ANY player's current position matches (x, y), draw the Head color.
	currentPos := core.Vector{X: x, Y: y}
	for _, player := range av.Arena.Players {
		if player.IsAlive && player.Position.Equals(currentPos) {
			// Safety check: ensure the ID is within the bounds of our color map
			if player.ID >= 1 && player.ID < NumCharacters {
				return av.Characters[player.ID-1].BrightColor
			}
		}
	}

	switch square {
	case core.Open:
		// Background color for open space
		return color.RGBA{R: 34, G: 32, B: 52, A: 255}
	case core.Wall:
		// Color for the arena border walls
		return color.RGBA{R: 200, G: 200, B: 200, A: 255}
	default:
		// --- Player Path / Head Logic ---
		playerID := int(square)

		// Safety check to ensure the ID is within the bounds of our color map
		if playerID < 1 || playerID >= NumCharacters {
			// Use a fallback color for unexpected IDs
			return color.RGBA{R: 255, G: 0, B: 255, A: 255}
		}

		headColor := av.Characters[playerID-1].BrightColor
		pathColor := av.Characters[playerID-1].DarkColor

		// Check if the player is alive and the square is their current "head" position.
		// Note: We access players using ID-1 because the Players slice is 0-indexed.
		if playerID <= len(av.Arena.Players) {
			player := av.Arena.Players[playerID-1]

			// Check if the current grid coordinate matches the player's head position.
			currentPos := core.Vector{X: x, Y: y}

			// Use player.IsAlive check to ensure dead players don't have a "Head"
			if player.IsAlive && player.Position.Equals(currentPos) {
				return headColor
			}
		}

		// If it's not the head (or the player is dead), return the path color.
		return pathColor
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
