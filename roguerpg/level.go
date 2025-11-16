package main

import (
	"math/rand"
)

type Level struct {
	WidthInTiles  int
	HeightInTiles int
	Tiles         [][]*Tile
	Enemies       []*BlobEnemy
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

func (level *Level) AddEnemies() {
	numEnemies := 25
	for range numEnemies {
		enemy := NewBlobEnemy()
		enemy.Location = level.FindRandomFloorLocation()
		level.Enemies = append(level.Enemies, enemy)
	}
}
