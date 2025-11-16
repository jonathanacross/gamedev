package main

import (
	"math/rand"
)

type Level struct {
	WidthInTiles  int
	HeightInTiles int
	Tiles         [][]*Tile
}

func (level *Level) FindRandomFloorLocation() Location {
	for {
		x := rand.Intn(level.WidthInTiles)
		y := rand.Intn(level.HeightInTiles)
		tile := level.Tiles[y][x]
		if !tile.solid {
			return Location{
				X: float64(x*TileSize + TileSize/2),
				Y: float64(y*TileSize + TileSize/2),
			}
		}
	}
}
