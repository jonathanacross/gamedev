package main

import (
	"math/rand/v2"
	"slices"
)

type Level struct {
	// view in grid coords
	height int
	width  int

	// location of last generated column
	lastGenX int

	// ideal safe vertical coordinate for player
	lastSafeY int

	spriteSheet *SpriteSheet
}

const (
	grassUp          = 0
	block            = 1
	spike            = 2
	pipeUpLeft       = 3
	pipeUpCenter     = 4
	pipeUpRight      = 5
	grassDown        = 6
	pipeDownLeft     = 9
	pipeDownCenter   = 10
	pipeDownRight    = 11
	target           = 14
	pipeMiddleLeft   = 15
	pipeMiddleCenter = 16
	pipeMiddleRight  = 17
)

func NewLevel() *Level {

	height := ScreenHeight / TileSize
	width := ScreenWidth / TileSize

	return &Level{
		height:      height,
		width:       width,
		lastGenX:    0,
		lastSafeY:   height / 2,
		spriteSheet: NewSpriteSheet(TerrainImage, TileSize, TileSize, 6, 3),
	}
}

func (l *Level) makeTile(gridX int, gridY int, id int) *Tile {
	return &Tile{
		Location: Location{
			X: float64(gridX) * TileSize,
			Y: float64(gridY) * TileSize,
		},
		spriteSheet: l.spriteSheet,
		srcRect:     l.spriteSheet.Rect(id),
	}
}

func (l *Level) Update(camera *Camera, tiles *[]*Tile) {
	cameraMinX := camera.GetViewRect().Min.X
	cameraMaxX := camera.GetViewRect().Max.X

	// Remove any old stuff that has gone offscreen.
	*tiles = slices.DeleteFunc(*tiles, func(t *Tile) bool {
		return t.X < float64(cameraMinX)-2*TileSize
	})

	// Check camera
	requiredGenX := (cameraMaxX / TileSize) + 1
	if l.lastGenX >= requiredGenX {
		return
	}

	// Camera is close to edge of world; create more stuff.

	widthToGenerate := requiredGenX - l.lastGenX + 5
	newMaxGenX := l.lastGenX + widthToGenerate

	// add some floor/ceiling
	for x := l.lastGenX + 1; x <= newMaxGenX; x++ {
		floor := l.makeTile(x, l.height-1, grassUp)
		ceiling := l.makeTile(x, 0, grassDown)
		*tiles = append(*tiles, floor)
		*tiles = append(*tiles, ceiling)
	}

	// and an obstacle
	obstacle := l.makeNewObstacle(newMaxGenX - 3)
	*tiles = append(*tiles, obstacle...)

	l.lastGenX = newMaxGenX
}

// generates a number from lo to hi, inclusive
func uniformRand(lo int, hi int) int {
	return lo + rand.IntN(hi-lo+1)
}

func twoOrderedRandomNumbersInRange(lo int, hi int) (int, int) {
	if hi < lo {
		return 0, -1
	}
	if hi == lo {
		return lo, hi
	}

	x1 := uniformRand(lo, hi)
	x2 := uniformRand(lo, hi)
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	return x1, x2
}

func (l *Level) makeNewObstacle(gridX int) []*Tile {
	// Figure out the next safe area, by perturbing the previous
	// safe area.  This should keep the levels feasible for the player.
	const MaxChange = 5
	const SafeRadius = 2
	groundLevel := l.height - 3
	ceilingLevel := 2
	nextSafeY := -1
	if l.lastSafeY <= ceilingLevel+SafeRadius {
		nextSafeY = l.lastSafeY + MaxChange
	} else if l.lastSafeY >= groundLevel-SafeRadius {
		nextSafeY = l.lastSafeY - MaxChange
	} else {
		nextSafeY = l.lastSafeY + uniformRand(-MaxChange, MaxChange)
		if nextSafeY > groundLevel-SafeRadius {
			nextSafeY = groundLevel - SafeRadius
		}
		if nextSafeY < ceilingLevel+SafeRadius {
			nextSafeY = ceilingLevel + SafeRadius
		}
	}
	safeBottom := nextSafeY + SafeRadius
	safeTop := nextSafeY - SafeRadius

	// Add obstacles
	tiles := []*Tile{}

	obstacleLocs := uniformRand(0, 2)
	addCeilingObstacle := obstacleLocs == 0 || obstacleLocs == 2
	addFloorObstacle := obstacleLocs == 1 || obstacleLocs == 2

	if addCeilingObstacle {
		addPipe := uniformRand(0, 1) == 0
		if addPipe {
			ceilPipeLen := safeTop - ceilingLevel + 1
			ceilPipe := l.makeCeilingPipe(gridX, ceilPipeLen)
			tiles = append(tiles, ceilPipe...)
		} else {
			blockWidth := uniformRand(3, 5)
			blockLo, blockHi := twoOrderedRandomNumbersInRange(ceilingLevel, safeTop-1)
			ceilBlocks := l.makeBlocks(gridX, blockWidth, blockLo, blockHi)
			tiles = append(tiles, ceilBlocks...)
		}
	}

	if addFloorObstacle {
		addPipe := uniformRand(0, 1) == 0
		if addPipe {
			floorPipeLen := groundLevel - safeBottom + 1
			floorPipe := l.makeFloorPipe(gridX, floorPipeLen)
			tiles = append(tiles, floorPipe...)
		} else {
			blockWidth := uniformRand(3, 5)
			blockLo, blockHi := twoOrderedRandomNumbersInRange(safeBottom+1, groundLevel)
			floorBlocks := l.makeBlocks(gridX, blockWidth, blockLo, blockHi)
			tiles = append(tiles, floorBlocks...)
		}
	}

	// debug: show safe path
	// for j := safeTop; j <= safeBottom; j++ {
	// 	debugTile := l.makeTile(gridX+1, j, target)
	// 	tiles = append(tiles, debugTile)
	// }

	// Add coin to a safe location.
	// TODO: replace spike symbol with coin
	coinY := nextSafeY + uniformRand(-SafeRadius+1, SafeRadius-1)
	coin := l.makeTile(gridX+1, coinY, spike)
	tiles = append(tiles, coin)

	l.lastSafeY = nextSafeY
	return tiles
}

func (l *Level) makeFloorPipe(x int, height int) []*Tile {
	pipeTiles := []*Tile{}
	groundLevel := l.height - 2

	for j := range height - 1 {
		y := groundLevel - j
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+1, y, pipeMiddleCenter)
		right := l.makeTile(x+2, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	y := groundLevel - height + 1
	upLeft := l.makeTile(x, y, pipeUpLeft)
	upCenter := l.makeTile(x+1, y, pipeUpCenter)
	upRight := l.makeTile(x+2, y, pipeUpRight)
	pipeTiles = append(pipeTiles, upLeft, upCenter, upRight)

	return pipeTiles
}

func (l *Level) makeCeilingPipe(x int, height int) []*Tile {
	pipeTiles := []*Tile{}
	const ceilingLevel = 1

	for j := range height - 1 {
		y := ceilingLevel + j
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+1, y, pipeMiddleCenter)
		right := l.makeTile(x+2, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	y := ceilingLevel + height - 1
	downLeft := l.makeTile(x, y, pipeDownLeft)
	downCenter := l.makeTile(x+1, y, pipeDownCenter)
	downRight := l.makeTile(x+2, y, pipeDownRight)
	pipeTiles = append(pipeTiles, downLeft, downCenter, downRight)

	return pipeTiles
}

func (l *Level) makeBlocks(x int, width int, top int, bottom int) []*Tile {
	blockTiles := []*Tile{}

	for i := range width {
		for j := 0; j <= bottom-top; j++ {
			blk := l.makeTile(x+i, top+j, block)
			blockTiles = append(blockTiles, blk)
		}
	}

	return blockTiles
}
