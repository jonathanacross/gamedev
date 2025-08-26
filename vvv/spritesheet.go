package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileSet interface {
	Image(int) *ebiten.Image
}

type ImageGroupTileSet struct {
	images []*ebiten.Image
	gid    int
}

func (s *ImageGroupTileSet) Image(index int) *ebiten.Image {
	return s.images[index-s.gid]
}

type GridTileSet struct {
	image         *ebiten.Image
	tileWidth     int
	tileHeight    int
	widthInTiles  int
	heightInTiles int
	gid           int
}

func NewGridTileSet(image *ebiten.Image, tileWidth int, tileHeight int, widthInTiles int, heightInTiles int) *GridTileSet {
	return &GridTileSet{
		image:         image,
		tileWidth:     tileWidth,
		tileHeight:    tileHeight,
		widthInTiles:  widthInTiles,
		heightInTiles: heightInTiles,
	}
}

func (ss *GridTileSet) Rect(index int) image.Rectangle {
	index -= ss.gid
	x := (index % ss.widthInTiles) * ss.tileWidth
	y := (index / ss.widthInTiles) * ss.tileHeight
	return image.Rect(x, y, x+ss.tileWidth, y+ss.tileHeight)
}

func (s *GridTileSet) Image(index int) *ebiten.Image {
	return s.image.SubImage(s.Rect(index)).(*ebiten.Image)
}
