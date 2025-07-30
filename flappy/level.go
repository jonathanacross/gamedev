package main

import (
	"math/rand/v2"
	"slices"
)

// Level manages the procedural generation of game obstacles, items, and enemies.
type Level struct {
	// view in grid coords
	height int
	width  int

	// location of last generated column
	lastGenX int

	// ideal safe vertical coordinate for player
	lastSafeY int

	spriteSheet *SpriteSheet

	// how many coins we've generated
	coinCounter int
}

// Constants for tile sprite sheet indices.
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
	target           = 14 // For debug or special markers
	pipeMiddleLeft   = 15
	pipeMiddleCenter = 16
	pipeMiddleRight  = 17
)

// NewLevel creates and initializes a new Level generator.
func NewLevel() *Level {
	height := ScreenHeight / TileSize
	width := ScreenWidth / TileSize

	return &Level{
		height:      height,
		width:       width,
		lastGenX:    0,
		lastSafeY:   height / 2, // Start with a safe vertical position in the middle
		spriteSheet: NewSpriteSheet(TerrainImage, TileSize, TileSize, 6, 3),
		coinCounter: 0,
	}
}

// makeTile creates a new Tile entity at the given grid coordinates with the specified sprite ID.
func (l *Level) makeTile(gridX int, gridY int, id int) *Tile {
	ss := l.spriteSheet
	return &Tile{
		BaseSprite: NewBaseSprite(
			float64(gridX)*TileSize,
			float64(gridY)*TileSize,
			ss,
			ss.Rect(id),
		),
	}
}

// makeCoin creates a new CoinItem entity at the given grid coordinates.
func makeCoin(gridX int, gridY int) Item {
	coinSS := NewSpriteSheet(CoinImage, TileSize, TileSize, 10, 1)
	anim := NewAnimation(0, 9, 10)
	return &CoinItem{
		AnimatedSprite: NewAnimatedSprite(
			float64(gridX)*TileSize,
			float64(gridY)*TileSize,
			coinSS,
			anim,
		),
	}
}

// makeHeart creates a new HeartItem entity at the given grid coordinates.
func makeHeart(gridX int, gridY int) Item {
	heartSS := NewSpriteSheet(HeartItemImage, TileSize, TileSize, 2, 1)
	anim := NewAnimation(0, 1, 10)
	return &HeartItem{
		AnimatedSprite: NewAnimatedSprite(
			float64(gridX)*TileSize,
			float64(gridY)*TileSize,
			heartSS,
			anim,
		),
	}
}

// makeOcto creates a new Octo enemy entity at the given grid X coordinate.
func (l *Level) makeOcto(gridX int) Enemy {
	octoSS := NewSpriteSheet(OctoImage, TileSize, TileSize, 2, 1)
	anim := NewAnimation(0, 1, 8)
	return &Octo{
		AnimatedSprite: NewAnimatedSprite(
			float64(gridX)*TileSize,
			ScreenHeight/2, // Initial Y, will be adjusted by Octo's update logic
			octoSS,
			anim,
		),
		minY:  float64(1 * TileSize),
		maxY:  float64(l.height-2) * TileSize,
		t:     rand.Float64() * 1000, // Random starting point for sine wave
		speed: 0.01,
	}
}

// makeBee creates a new Bee enemy entity at the given grid coordinates.
func (l *Level) makeBee(gridX int, gridY int) Enemy {
	beeSS := NewSpriteSheet(BeeImage, TileSize, TileSize, 2, 1)
	anim := NewAnimation(0, 1, 8)
	return &Bee{
		AnimatedSprite: NewAnimatedSprite(
			float64(gridX)*TileSize,
			float64(gridY)*TileSize,
			beeSS,
			anim,
		),
		speed: 1, // Horizontal speed
	}
}

// Update manages the level's procedural generation, adding new content as the camera moves.
func (l *Level) Update(camera *Camera, tiles *[]*Tile, items *[]Item, enemies *[]Enemy) {
	cameraMinX := camera.GetViewRect().Min.X
	cameraMaxX := camera.GetViewRect().Max.X

	// Remove entities that have gone off-screen to save memory and processing.
	*tiles = slices.DeleteFunc(*tiles, func(t *Tile) bool {
		return t.X < float64(cameraMinX)-2*TileSize
	})
	*items = slices.DeleteFunc(*items, func(i Item) bool {
		return i.GetX() < float64(cameraMinX)-2*TileSize
	})
	*enemies = slices.DeleteFunc(*enemies, func(e Enemy) bool {
		return e.GetX() < float64(cameraMinX)-2*TileSize
	})

	// Determine if new content needs to be generated.
	requiredGenX := (cameraMaxX / TileSize) + 1
	if l.lastGenX >= requiredGenX {
		return // No generation needed yet
	}

	// Generate new content for the level.
	widthToGenerate := requiredGenX - l.lastGenX + 5
	newMaxGenX := l.lastGenX + widthToGenerate

	// Add floor and ceiling tiles for the newly generated width.
	for x := l.lastGenX + 1; x <= newMaxGenX; x++ {
		floor := l.makeTile(x, l.height-1, grassUp)
		ceiling := l.makeTile(x, 0, grassDown)
		*tiles = append(*tiles, floor, ceiling)
	}

	// Generate obstacles, coins, and enemies for the new section.
	obstacles, coins, octos := l.makeNewObstacle(newMaxGenX - 3) // Adjust X for obstacle placement
	*tiles = append(*tiles, obstacles...)
	*items = append(*items, coins...)
	*enemies = append(*enemies, octos...)

	l.lastGenX = newMaxGenX // Update the last generated column
}

// uniformRandInt generates a random integer between lo and hi, inclusive.
func uniformRandInt(lo int, hi int) int {
	return lo + rand.IntN(hi-lo+1)
}

// twoOrderedRandomIntsInRange generates two random integers within a range, ensuring x1 <= x2.
func twoOrderedRandomIntsInRange(lo int, hi int) (int, int) {
	if hi < lo {
		return 0, -1
	}
	if hi == lo {
		return lo, hi
	}

	x1 := uniformRandInt(lo, hi)
	x2 := uniformRandInt(lo, hi)
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	return x1, x2
}

// makeNewObstacle generates a section of the level including obstacles, items, and enemies.
func (l *Level) makeNewObstacle(gridX int) ([]*Tile, []Item, []Enemy) {
	// Figure out the next safe area for the player by perturbing the previous safe area.
	// This ensures the level remains traversable.
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
		nextSafeY = l.lastSafeY + uniformRandInt(-MaxChange, MaxChange)
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

	obstacleLocs := uniformRandInt(0, 2) // 0: ceiling, 1: floor, 2: both
	addCeilingObstacle := obstacleLocs == 0 || obstacleLocs == 2
	addFloorObstacle := obstacleLocs == 1 || obstacleLocs == 2

	if addCeilingObstacle {
		addPipe := uniformRandInt(0, 1) == 0 // Randomly choose between pipe or blocks
		if addPipe {
			ceilPipeLen := safeTop - ceilingLevel + 1
			ceilPipe := l.makeCeilingPipe(gridX, ceilPipeLen)
			tiles = append(tiles, ceilPipe...)
		} else {
			blockWidth := uniformRandInt(1, 4)
			blockLo, blockHi := twoOrderedRandomIntsInRange(ceilingLevel, safeTop-1)
			ceilBlocks := l.makeBlocks(gridX, blockWidth, blockLo, blockHi)
			tiles = append(tiles, ceilBlocks...)
		}
	}

	if addFloorObstacle {
		addPipe := uniformRandInt(0, 1) == 0 // Randomly choose between pipe or blocks
		if addPipe {
			floorPipeLen := groundLevel - safeBottom + 1
			floorPipe := l.makeFloorPipe(gridX, floorPipeLen)
			tiles = append(tiles, floorPipe...)
		} else {
			blockWidth := uniformRandInt(1, 4)
			blockLo, blockHi := twoOrderedRandomIntsInRange(safeBottom+1, groundLevel)
			floorBlocks := l.makeBlocks(gridX, blockWidth, blockLo, blockHi)
			tiles = append(tiles, floorBlocks...)
		}
	}

	// Add item (coin or heart) to a safe location.
	items := []Item{}
	itemX := gridX + 1
	itemY := nextSafeY + uniformRandInt(-SafeRadius+1, SafeRadius-1)
	if l.coinCounter == DropHeartsEveryNCoins {
		l.coinCounter = 0
		heart := makeHeart(itemX, itemY)
		items = append(items, heart)
	} else {
		l.coinCounter++
		coin := makeCoin(itemX, itemY)
		items = append(items, coin)
	}

	// Add enemies
	enemies := []Enemy{}
	makeEnemy := uniformRandInt(0, 1) == 0 // 50% chance to make an enemy
	if makeEnemy {
		makeBee := uniformRandInt(0, 1) == 0 // Randomly choose between Bee or Octo
		if makeBee {
			beeX := gridX + 5
			beeY := uniformRandInt(safeTop, safeBottom)
			bee := l.makeBee(beeX, beeY)
			enemies = append(enemies, bee)
		} else {
			octoX := uniformRandInt(gridX, gridX+4)
			octo := l.makeOcto(octoX)
			enemies = append(enemies, octo)
		}
	}

	l.lastSafeY = nextSafeY // Update the last safe Y coordinate for the next generation cycle
	return tiles, items, enemies
}

// makeFloorPipe generates a vertical pipe structure extending from the floor upwards.
func (l *Level) makeFloorPipe(x int, height int) []*Tile {
	pipeTiles := []*Tile{}
	groundLevel := l.height - 2 // Ground level is 2 tiles from bottom

	// Middle sections of the pipe
	for j := 0; j < height-1; j++ {
		y := groundLevel - j
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+1, y, pipeMiddleCenter)
		right := l.makeTile(x+2, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	// Top cap of the pipe
	y := groundLevel - height + 1
	upLeft := l.makeTile(x, y, pipeUpLeft)
	upCenter := l.makeTile(x+1, y, pipeUpCenter)
	upRight := l.makeTile(x+2, y, pipeUpRight)
	pipeTiles = append(pipeTiles, upLeft, upCenter, upRight)

	return pipeTiles
}

// makeCeilingPipe generates a vertical pipe structure extending from the ceiling downwards.
func (l *Level) makeCeilingPipe(x int, height int) []*Tile {
	pipeTiles := []*Tile{}
	const ceilingLevel = 1 // Ceiling level is 1 tile from top

	// Middle sections of the pipe
	for j := 0; j < height-1; j++ {
		y := ceilingLevel + j
		left := l.makeTile(x, y, pipeMiddleLeft)
		center := l.makeTile(x+1, y, pipeMiddleCenter)
		right := l.makeTile(x+2, y, pipeMiddleRight)
		pipeTiles = append(pipeTiles, left, center, right)
	}

	// Bottom cap of the pipe
	y := ceilingLevel + height - 1
	downLeft := l.makeTile(x, y, pipeDownLeft)
	downCenter := l.makeTile(x+1, y, pipeDownCenter)
	downRight := l.makeTile(x+2, y, pipeDownRight)
	pipeTiles = append(pipeTiles, downLeft, downCenter, downRight)

	return pipeTiles
}

// makeBlocks generates a rectangular block obstacle.
func (l *Level) makeBlocks(x int, width int, top int, bottom int) []*Tile {
	blockTiles := []*Tile{}

	for i := 0; i < width; i++ {
		for j := 0; j <= bottom-top; j++ {
			blk := l.makeTile(x+i, top+j, block)
			blockTiles = append(blockTiles, blk)
		}
	}

	return blockTiles
}
