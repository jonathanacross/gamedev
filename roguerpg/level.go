package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Level struct {
	WidthInTiles  int
	HeightInTiles int
	Tiles         [][]*Tile
	Enemies       []Character
	Objects       []GameObject
	DownLevel     *Level
	UpLevel       *Level
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

func (level *Level) AddEnemies(depth int) {
	// TODO: scale type of enemies based on depth
	numEnemies := 25
	for range numEnemies {
		location := level.FindRandomFloorLocation()
		var enemy Character
		enemyType := rand.Intn(3)
		switch enemyType {
		case 0:
			enemy = NewBlobEnemy(location)
		case 1:
			enemy = NewBatEnemy(location)
		case 2:
			enemy = NewGoblinEnemy(location)
		}
		level.Enemies = append(level.Enemies, enemy)
	}
}

func (level *Level) AddObjects(isFirstLevel bool, isFinalLevel bool) {
	if !isFinalLevel {
		// no chests on final level
		numChests := 2
		for range numChests {
			chest := NewChest(level.FindRandomFloorLocation())
			level.Objects = append(level.Objects, chest)
		}

		downstairs := NewStairs(level.FindRandomFloorLocation(), false)
		level.Objects = append(level.Objects, downstairs)
	}
	if !isFirstLevel {
		upstairs := NewStairs(level.FindRandomFloorLocation(), true)
		level.Objects = append(level.Objects, upstairs)
	}
}

func (level *Level) GetUpstairsLocation() (Location, error) {
	for _, obj := range level.Objects {
		if stairs, ok := obj.(*Stairs); ok && stairs.IsUpstairs {
			return stairs.Location(), nil
		}
	}
	return Location{}, fmt.Errorf("no upstairs found in level")
}

func (level *Level) GetDownstairsLocation() (Location, error) {
	for _, obj := range level.Objects {
		if stairs, ok := obj.(*Stairs); ok && !stairs.IsUpstairs {
			return stairs.Location(), nil
		}
	}
	return Location{}, fmt.Errorf("no downstairs found in level")
}
