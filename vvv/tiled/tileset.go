package tiled

import (
	"fmt"
	"image"
)

// --------------- Public interface -----------

type Tile struct {
	id         int
	srcRect    Rect
	srcImage   *image.Image
	hitRect    Rect
	properties *PropertySet
}

// ----------  Tiled JSON structs --------------

type ObjectJSON struct {
	Height     float64          `json:"height"`
	ID         int              `json:"id"`
	Gid        int              `json:"gid"`
	Name       string           `json:"name"`
	Properties []PropertiesJSON `json:"properties"`
	Type       string           `json:"type"`
	Width      float64          `json:"width"`
	X          float64          `json:"x"`
	Y          float64          `json:"y"`
}

type ObjectGroupJSON struct {
	DrawOrder string       `json:"draworder"`
	Name      string       `json:"name"`
	Objects   []ObjectJSON `json:"objects"`
	Opacity   float64      `json:"opacity"`
	Type      string       `json:"type"`
	Visible   bool         `json:"visible"`
	X         float64      `json:"x"`
	Y         float64      `json:"y"`
}

type TilesetTileJSON struct {
	ID          int              `json:"id"`
	Type        string           `json:"type"`
	Image       string           `json:"image"`
	ImageHeight int              `json:"imageheight"`
	ImageWidth  int              `json:"imagewidth"`
	Properties  []PropertiesJSON `json:"properties"`
	ObjectGroup *ObjectGroupJSON `json:"objectgroup,omitempty"`
}

type TilesetDataJSON struct {
	Columns     int               `json:"columns"`
	FirstGid    int               `json:"firstgid"`
	Image       string            `json:"image"`
	ImageHeight int               `json:"imageheight"`
	ImageWidth  int               `json:"imagewidth"`
	Name        string            `json:"name"`
	TileCount   int               `json:"tilecount"`
	TileHeight  int               `json:"tileheight"`
	Tiles       []TilesetTileJSON `json:"tiles"`
	TileWidth   int               `json:"tilewidth"`
	Type        string            `json:"type"`
}

// -----------  Internal conversion functions -----------

func isCollectionTileSet(tsj TilesetDataJSON) bool {
	return tsj.Image == ""
}

func getHitbox(tj *TilesetTileJSON, width float64, height float64) Rect {
	// See if there is a custom object defining
	// a hitbox.  (Note: there may be more than one
	// box defined in tiled, but we just look at the
	// first one.
	if tj.ObjectGroup != nil {
		return Rect{
			X:      tj.ObjectGroup.Objects[0].X,
			Y:      tj.ObjectGroup.Objects[0].Y,
			Width:  tj.ObjectGroup.Objects[0].Width,
			Height: tj.ObjectGroup.Objects[0].Height,
		}
	} else {
		// Build a default hitbox based on the size of the image
		return Rect{
			X:      0,
			Y:      0,
			Width:  width,
			Height: height,
		}
	}
}

func convertCollectionTilesetJSON(tsj TilesetDataJSON, gidStart int, images map[string]*image.Image) ([]Tile, error) {
	tiles := make([]Tile, tsj.TileCount)

	for i, tj := range tsj.Tiles {
		properties, err := GetProperties(tj.Properties)
		if err != nil {
			return nil, err
		}

		srcImage, ok := images[tj.Image]
		if !ok {
			return nil, fmt.Errorf("image %s not found", tj.Image)
		}

		srcRect := Rect{
			X:      0,
			Y:      0,
			Width:  float64(tj.ImageWidth),
			Height: float64(tj.ImageHeight),
		}

		t := Tile{
			id:         tj.ID + gidStart,
			srcRect:    srcRect,
			srcImage:   srcImage,
			hitRect:    getHitbox(&tj, float64(tj.ImageWidth), float64(tj.ImageHeight)),
			properties: properties,
		}
		tiles[i] = t
	}
	return tiles, nil
}

func convertGridTilesetJSON(tsj *TilesetDataJSON, gidStart int, images map[string]*image.Image) ([]Tile, error) {
	tiles := make([]Tile, tsj.TileCount)
	srcImage, ok := images[tsj.Image]
	if !ok {
		return nil, fmt.Errorf("image %s not found", tsj.Image)
	}
	tileWidth := float64(tsj.TileWidth)
	tileHeight := float64(tsj.TileHeight)
	columns := tsj.Columns

	for i, tj := range tsj.Tiles {
		x := (i % columns) * int(tileWidth)
		y := (i / columns) * int(tileHeight)
		srcRect := Rect{
			X:      float64(x),
			Y:      float64(y),
			Width:  tileWidth,
			Height: tileHeight,
		}

		properties, err := GetProperties(tj.Properties)
		if err != nil {
			return nil, err
		}

		t := Tile{
			id:         tj.ID + gidStart,
			srcRect:    srcRect,
			srcImage:   srcImage,
			hitRect:    getHitbox(&tj, tileWidth, tileHeight),
			properties: properties,
		}
		tiles[i] = t
	}
	return tiles, nil
}

// Main conversion method
func GetTiles(tsj *TilesetDataJSON, gidStart int, images map[string]*image.Image) ([]Tile, error) {
	if isCollectionTileSet(*tsj) {
		return convertCollectionTilesetJSON(*tsj, gidStart, images)
	} else {
		return convertGridTilesetJSON(tsj, gidStart, images)
	}
}
