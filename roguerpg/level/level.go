package level

import (
	"fmt"
	"math"
	"math/rand"
	"roguerpg/core"
)

type Level struct {
	WidthInTiles  int
	HeightInTiles int
	Tiles         [][]*core.Tile
	Enemies       []core.Character
	Objects       []core.GameObject
	DownLevel     *Level
	UpLevel       *Level
}

func (level *Level) GetTile(x, y int) *core.Tile {
	if x < 0 || x >= level.WidthInTiles || y < 0 || y >= level.HeightInTiles {
		return nil
	}
	return level.Tiles[y][x]
}

// TileToWorld converts tile coordinates (tx, ty) to the center of the world coordinates.
func (level *Level) TileToWorld(tx, ty int) core.Location {
	// Return the center point of the tile
	return core.Location{
		X: float64(tx*core.TileSize) + float64(core.TileSize)/2,
		Y: float64(ty*core.TileSize) + float64(core.TileSize)/2,
	}
}

// WorldToTile converts world coordinates (Location) to tile coordinates.
// It uses floor to get the tile index.
func (level *Level) WorldToTile(l core.Location) (int, int) {
	tx := int(math.Floor(l.X / float64(core.TileSize)))
	ty := int(math.Floor(l.Y / float64(core.TileSize)))
	return tx, ty
}

// IsTileSolid checks if a tile at (tx, ty) is solid, with bounds checking.
func (level *Level) IsTileSolid(tx, ty int) bool {
	if tx < 0 || tx >= level.WidthInTiles || ty < 0 || ty >= level.HeightInTiles {
		// Treat out-of-bounds as solid to prevent enemies from escaping
		return true
	}
	// Tiles is [][]*core.Tile.

	// We'll trust level gen.
	return level.Tiles[ty][tx].Solid
}

func (level *Level) FindRandomFloorLocation() core.Location {
	for {
		x := rand.Intn(level.WidthInTiles)
		y := rand.Intn(level.HeightInTiles)
		tile := level.Tiles[y][x]
		if !tile.Solid {
			return level.TileToWorld(x, y)
		}
	}
}

// AddEnemies and AddObjects methods removed.
// They will be implemented as functions in `main` package or independent factories
// that take *Level as argument. We should move them later.

func (level *Level) GetUpstairsLocation() (core.Location, error) {
	// Problem: Stairs is in `objects` package.
	// `level` depends on `objects`? NO. Dependencies: level -> core.
	// `objects` -> core.
	// If we import `roguerpg/objects`, we make `level` depend on `objects`.
	// The plan says `level` depends on `core`.
	// So `level` CANNOT know about `Stairs` struct directly if `Stairs` is in `objects`.

	// Solutions:
	// 1. Check for `IsUpstairs` property via interface? `core.Interactable`?
	//    `core.Interactable` doesn't have `IsUpstairs`.
	// 2. Add `IsUpstairs() bool` to an interface in `core`?
	// 3. Keep `GetUpstairsLocation` but rely on some `Tag` or `Type`?
	//    `GameObject` doesn't have Type.
	// 4. Move `GetUpstairsLocation` to NOT be a method of Level, but a function in `main` where we know about both?
	//    This is cleaner.

	// I will REMOVE GetUpstairsLocation and GetDownstairsLocation from Level struct.
	// They are utility lookups.
	return core.Location{}, fmt.Errorf("method moved to main")
}

func (level *Level) GetDownstairsLocation() (core.Location, error) {
	return core.Location{}, fmt.Errorf("method moved to main")
}

// Returns the nearest enemy to the given location within the maxDist
// If no enemy is found within the maxDist, returns nil
func (level *Level) FindNearestEnemy(location core.Location, maxDist float64) core.Character {
	var nearestEnemy core.Character
	minDist := maxDist

	for _, enemy := range level.Enemies {
		if enemy.IsDead() {
			continue
		}
		dist := core.Vector(location).Minus(core.Vector(enemy.Location())).Length()
		if dist < minDist {
			minDist = dist
			nearestEnemy = enemy
		}
	}
	return nearestEnemy
}

func (level *Level) GetEnemies() []core.Character {
	return level.Enemies
}

func (level *Level) GetObjects() []core.GameObject {
	return level.Objects
}
