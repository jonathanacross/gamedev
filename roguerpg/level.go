package main

import (
	"math"
	"math/rand"
)

type Level struct {
	WidthInTiles  int
	HeightInTiles int
	Tiles         [][]*Tile
	Enemies       []*BlobEnemy
}

func (level *Level) GetTile(x, y int) *Tile {
	if x < 0 || x >= level.WidthInTiles || y < 0 || y >= level.HeightInTiles {
		return nil
	}
	return level.Tiles[y][x]
}

// TileToWorld converts tile coordinates (tx, ty) to the center of the world coordinates.
func (level *Level) TileToWorld(tx, ty int) Location {
	// Return the center point of the tile
	return Location{
		X: float64(tx*TileSize) + TileSize/2,
		Y: float64(ty*TileSize) + TileSize/2,
	}
}

// WorldToTile converts world coordinates (Location) to tile coordinates.
// It uses floor to get the tile index.
func (level *Level) WorldToTile(l Location) (int, int) {
	tx := int(math.Floor(l.X / TileSize))
	ty := int(math.Floor(l.Y / TileSize))
	return tx, ty
}

// IsTileSolid checks if a tile at (tx, ty) is solid, with bounds checking.
func (level *Level) IsTileSolid(tx, ty int) bool {
	if tx < 0 || tx >= level.WidthInTiles || ty < 0 || ty >= level.HeightInTiles {
		// Treat out-of-bounds as solid to prevent enemies from escaping
		return true
	}
	return level.Tiles[ty][tx].solid
}

func (level *Level) FindRandomFloorLocation() Location {
	for {
		x := rand.Intn(level.WidthInTiles)
		y := rand.Intn(level.HeightInTiles)
		tile := level.Tiles[y][x]
		if !tile.solid {
			return level.TileToWorld(x, y)
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
