package main

import (
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
