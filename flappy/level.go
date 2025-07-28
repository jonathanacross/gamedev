package main

import (
	"math"
	"math/rand/v2"
	"slices"
)

type Level struct {
	// view in grid coords
	height int
	width  int

	// location of last generated column
	lastGenX int

	spriteSheet *SpriteSheet
}

const (
	grassUp          = 0
	block            = 1
	pipeUpLeft       = 3
	pipeUpCenter     = 4
	pipeUpRight      = 5
	grassDown        = 6
	pipeDownLeft     = 9
	pipeDownCenter   = 10
	pipeDownRight    = 11
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

func (l *Level) makeNewObstacle(gridX int) []*Tile {
	maxPipeLength := 2 * (l.height - 6) / 3

	thing := rand.IntN(3)
	switch thing {
	case 0:
		size := rand.IntN(maxPipeLength) + 3
		return l.makeUpPipe(gridX, size)
	case 1:
		size := rand.IntN(maxPipeLength) + 3
		return l.makeDownPipe(gridX, size)
	case 2:
		width := rand.IntN(3) + 2
		height := rand.IntN(2) + 1
		position := rand.Float64()
		return l.makeBlocks(gridX, width, height, position)
	}

	return []*Tile{}
}

func (l *Level) makeUpPipe(x int, height int) []*Tile {
	pipeTiles := []*Tile{}
	groundLevel := l.height - 2

	for j := range height {
		y := groundLevel - j
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+1, y, pipeMiddleCenter)
		right := l.makeTile(x+2, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	y := groundLevel - height
	upLeft := l.makeTile(x, y, pipeUpLeft)
	upCenter := l.makeTile(x+1, y, pipeUpCenter)
	upRight := l.makeTile(x+2, y, pipeUpRight)
	pipeTiles = append(pipeTiles, upLeft, upCenter, upRight)

	return pipeTiles
}

func (l *Level) makeDownPipe(x int, height int) []*Tile {
	pipeTiles := []*Tile{}
	const ceilingLevel = 1

	for j := range height {
		y := ceilingLevel + j
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+1, y, pipeMiddleCenter)
		right := l.makeTile(x+2, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	y := ceilingLevel + height
	downLeft := l.makeTile(x, y, pipeDownLeft)
	downCenter := l.makeTile(x+1, y, pipeDownCenter)
	downRight := l.makeTile(x+2, y, pipeDownRight)
	pipeTiles = append(pipeTiles, downLeft, downCenter, downRight)

	return pipeTiles
}

func (l *Level) makeBlocks(x int, width int, height int, position float64) []*Tile {
	blockTiles := []*Tile{}
	lo := 2
	hi := l.height - 3 - height

	// TODO: change to pure int, no floating point
	startY := int(math.Floor(float64(lo) + position*float64(hi-lo)))

	for i := range width {
		for j := range height {
			tx := x + i
			ty := startY + j
			blk := l.makeTile(tx, ty, block)
			blockTiles = append(blockTiles, blk)
		}
	}

	return blockTiles
}
