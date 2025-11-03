package main

import (
	"fmt"
	"image"
)

type SpriteSheet struct {
	tileWidth     int
	tileHeight    int
	widthInTiles  int
	heightInTiles int
	gid           int
}

func NewSpriteSheet(tileWidth int, tileHeight int, widthInTiles int, heightInTiles int) *SpriteSheet {
	return &SpriteSheet{
		tileWidth:     tileWidth,
		tileHeight:    tileHeight,
		widthInTiles:  widthInTiles,
		heightInTiles: heightInTiles,
	}
}

func (ss *SpriteSheet) Rect(index int) image.Rectangle {
	index -= ss.gid
	x := (index % ss.widthInTiles) * ss.tileWidth
	y := (index / ss.widthInTiles) * ss.tileHeight
	return image.Rect(x, y, x+ss.tileWidth, y+ss.tileHeight)
}

// Spritesheet for a 7x7 blob tile set
// see https://www.boristhebrave.com/permanent/24/06/cr31/stagecast/art/atlas/blob/wangbl.png
type BlobSpriteSheet struct {
	tileWidth   int
	tileHeight  int
	lookupTable map[int]int
}

func NewBlobSpriteSheet(tileWidth int, tileHeight int) *BlobSpriteSheet {
	layout := []int{
		0, 4, 92, 124, 116, 80, 0,
		16, 20, 87, 223, 241, 21, 64,
		29, 117, 85, 71, 221, 125, 112,
		31, 253, 113, 28, 127, 247, 209,
		23, 199, 213, 95, 255, 245, 81,
		5, 84, 93, 119, 215, 193, 17,
		0, 1, 7, 197, 69, 68, 65,
	}

	lookup := make(map[int]int)
	for idx, blobID := range layout {
		lookup[blobID] = idx
	}

	return &BlobSpriteSheet{
		tileWidth:   tileWidth,
		tileHeight:  tileHeight,
		lookupTable: lookup,
	}
}

func (bss *BlobSpriteSheet) Rect(value int) image.Rectangle {
	index, ok := bss.lookupTable[value]
	if !ok {
		fmt.Printf("Warning: Blob value %d not found in lookup table, defaulting to 0\n", value)
		index = 0
	}
	x := (index % 7) * bss.tileWidth
	y := (index / 7) * bss.tileHeight
	return image.Rect(x, y, x+bss.tileWidth, y+bss.tileHeight)
}
