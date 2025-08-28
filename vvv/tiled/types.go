package tiled

import "image"

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type Object struct {
	Name    string
	Type    string
	HitRect Rect

	SrcRect  Rect
	SrcImage *image.Image

	Properties *PropertySet
}

type MapLayer struct {
	Name     string
	TileData []int
	Objects  []Object
}

type Map struct {
	Name          string
	WidthInTiles  int
	HeightInTiles int
	TileWidth     int
	TileHeight    int

	Layers []MapLayer
	Tiles  []Tile
}

// func (m *Map) findLayer(name string) (layer, ok) {
//  ...
// }
