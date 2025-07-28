package main

import (
	"math"
	"math/rand/v2"
	"slices"
	"time"
)

type Level struct {
	height   int
	width    int
	tileSize int

	spriteSheet   *SpriteSheet
	obstacleTimer *Timer
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
		height:        height,
		width:         width,
		tileSize:      TileSize,
		spriteSheet:   NewSpriteSheet(TerrainImage, TileSize, TileSize, 6, 3),
		obstacleTimer: NewTimer(1500 * time.Millisecond),
	}
}

func (l *Level) makeTile(x, y float64, id int) *Tile {
	return &Tile{
		Location: Location{
			X: x,
			Y: y,
		},
		spriteSheet: l.spriteSheet,
		srcRect:     l.spriteSheet.Rect(id),
	}
}

func (l *Level) Update(ground *[]*Tile, obstacles *[]*Tile) {
	// move everything as time goes on
	const speed = 1
	for _, t := range *ground {
		t.X -= speed
	}
	for _, t := range *obstacles {
		t.X -= speed
	}

	// remove any old stuff that has gone offscreen
	*ground = slices.DeleteFunc(*ground, func(t *Tile) bool {
		return t.X < float64(-l.tileSize)
	})
	*obstacles = slices.DeleteFunc(*obstacles, func(t *Tile) bool {
		return t.X < float64(-l.tileSize)
	})

	// generate new ground, as needed
	maxTileRight := 0.0
	for _, t := range *ground {
		if t.X > maxTileRight {
			maxTileRight = t.X
		}
	}

	for maxTileRight < float64(l.width*l.tileSize) {
		floor := l.makeTile(maxTileRight, float64((l.height-1)*l.tileSize), grassUp)
		ceiling := l.makeTile(maxTileRight, 0, grassDown)
		*ground = append(*ground, floor)
		*ground = append(*ground, ceiling)

		maxTileRight += float64(l.tileSize)
	}

	// generate new obstacles, as needed
	l.obstacleTimer.Update()
	if l.obstacleTimer.IsReady() {
		l.obstacleTimer.Reset()

		obstacle := l.makeNewObstacle(maxTileRight)
		*obstacles = append(*obstacles, obstacle...)
	}
}

func (l *Level) makeNewObstacle(x float64) []*Tile {
	thing := rand.IntN(3)
	switch thing {
	case 0:
		size := rand.IntN(3) + 3
		return l.makeUpPipe(x, size)
	case 1:
		size := rand.IntN(3) + 3
		return l.makeDownPipe(x, size)
	case 2:
		width := rand.IntN(3) + 2
		height := rand.IntN(2) + 1
		position := rand.Float64()
		return l.makeBlocks(x, width, height, position)
	}

	return []*Tile{}
}

func (l *Level) makeUpPipe(x float64, height int) []*Tile {
	pipeTiles := []*Tile{}
	ts := float64(l.tileSize)
	groundLevel := float64(l.height-2) * ts

	for j := range height {
		y := groundLevel - float64(j)*ts
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+ts, y, pipeMiddleCenter)
		right := l.makeTile(x+2*ts, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	y := groundLevel - float64(height)*ts
	upLeft := l.makeTile(x, y, pipeUpLeft)
	upCenter := l.makeTile(x+ts, y, pipeUpCenter)
	upRight := l.makeTile(x+2*ts, y, pipeUpRight)
	pipeTiles = append(pipeTiles, upLeft, upCenter, upRight)

	return pipeTiles
}

func (l *Level) makeDownPipe(x float64, height int) []*Tile {
	pipeTiles := []*Tile{}
	ts := float64(l.tileSize)
	ceilingLevel := ts

	for j := range height {
		y := ceilingLevel + float64(j)*ts
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+ts, y, pipeMiddleCenter)
		right := l.makeTile(x+2*ts, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	y := ceilingLevel + float64(height)*ts
	downLeft := l.makeTile(x, y, pipeDownLeft)
	downCenter := l.makeTile(x+ts, y, pipeDownCenter)
	downRight := l.makeTile(x+2*ts, y, pipeDownRight)
	pipeTiles = append(pipeTiles, downLeft, downCenter, downRight)

	return pipeTiles
}

func (l *Level) makeBlocks(x float64, width, height int, position float64) []*Tile {
	blockTiles := []*Tile{}
	ts := float64(l.tileSize)
	lo := 2
	hi := l.height - 3 - height

	startY := math.Floor(float64(lo)+position*float64(hi-lo)) * ts

	for i := range width {
		for j := range height {
			tx := x + float64(i)*ts
			ty := startY + float64(j)*ts
			blk := l.makeTile(tx, ty, block)
			blockTiles = append(blockTiles, blk)
		}
	}

	return blockTiles
}
