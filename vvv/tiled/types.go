package tiled

// ImageProvider represents an image-like type,
// such as *image.Image or *ebiten.Image.
type ImageProvider interface{}

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type Property struct {
	Value interface{}
}

// Tile represents a single tile with its properties and image source.
type Tile struct {
	ID         int
	SrcRect    Rect
	SrcImage   ImageProvider
	HitRect    Rect
	Properties *PropertySet
}

// A property set is just a map of key value pairs.
// The values are Typed, and must be one of bool, int, float64, string,
// according to the setup in Tiled.
type PropertySet map[string]Property

// Object represents a single object layer element.
type Object struct {
	Name       string
	Type       string
	Properties *PropertySet
	Location   Rect
	GID        int
}

// MapLayer represents a single layer in the map.
type MapLayer struct {
	Name    string
	Type    string
	Width   int
	Height  int
	TileIds []int
	Objects []Object
}

// Map represents the entire Tiled map file.
type Map struct {
	Name          string
	WidthInTiles  int
	HeightInTiles int
	TileWidth     int
	TileHeight    int

	Layers []MapLayer
	Tiles  []Tile
}
