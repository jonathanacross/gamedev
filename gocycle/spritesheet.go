package main

import (
	"image"
)

// TODO: rename
type SpriteSheet struct {
	// image         *ebiten.Image
	tileWidth     int
	tileHeight    int
	widthInTiles  int
	heightInTiles int
	gid           int
}

func NewSpriteSheet(tileWidth int, tileHeight int, widthInTiles int, heightInTiles int) *SpriteSheet {
	return &SpriteSheet{
		// image:         image,
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
