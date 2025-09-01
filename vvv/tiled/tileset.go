package tiled

import (
	"fmt"
)

// ConvertTileset converts an intermediate tiledTileset struct into a slice of
// game-ready Tile structs. It receives the loaded image(s), so it doesn't have
// to worry about file I/O.
func ConvertTileset(tsData *tiledTileset, images map[string]ImageProvider, firstGID int) ([]Tile, error) {
	// Check if this is a collection tileset (individual images) or a sprite sheet.
	isCollection := tsData.Image == ""

	if isCollection {
		return convertCollectionTileset(tsData, images, firstGID)
	}
	return convertSpriteSheetTileset(tsData, images[tsData.Image], firstGID)
}

// convertCollectionTileset handles tilesets with individual tile images.
func convertCollectionTileset(tsData *tiledTileset, images map[string]ImageProvider, firstGID int) ([]Tile, error) {
	tiles := make([]Tile, 0, len(tsData.Tiles))
	for _, tiledTile := range tsData.Tiles {
		properties, err := GetProperties(tiledTile.Properties)
		if err != nil {
			return nil, err
		}
		hitRect := getHitbox(&tiledTile, float64(tiledTile.ImageWidth), float64(tiledTile.ImageHeight))

		// Get the correct image from the map
		img, ok := images[tiledTile.Image]
		if !ok {
			return nil, fmt.Errorf("image not found for tile %d: %s", tiledTile.ID, tiledTile.Image)
		}

		tile := Tile{
			ID:         tiledTile.ID + firstGID,
			SrcRect:    Rect{X: 0, Y: 0, Width: float64(tiledTile.ImageWidth), Height: float64(tiledTile.ImageHeight)},
			SrcImage:   img,
			HitRect:    hitRect,
			Properties: &properties,
		}
		tiles = append(tiles, tile)
	}
	return tiles, nil
}

// convertSpriteSheetTileset handles tilesets that use a single sprite sheet image.
func convertSpriteSheetTileset(tsData *tiledTileset, srcImage ImageProvider, firstGID int) ([]Tile, error) {
	tiles := make([]Tile, 0, tsData.TileCount)
	tileWidth := float64(tsData.TileWidth)
	tileHeight := float64(tsData.TileHeight)
	columns := tsData.Columns

	for _, tiledTile := range tsData.Tiles {
		x := (tiledTile.ID % columns) * int(tileWidth)
		y := (tiledTile.ID / columns) * int(tileHeight)
		srcRect := Rect{X: float64(x), Y: float64(y), Width: tileWidth, Height: tileHeight}

		properties, err := GetProperties(tiledTile.Properties)
		if err != nil {
			return nil, err
		}
		hitRect := getHitbox(&tiledTile, tileWidth, tileHeight)

		tile := Tile{
			ID:         tiledTile.ID + firstGID,
			SrcRect:    srcRect,
			SrcImage:   srcImage,
			HitRect:    hitRect,
			Properties: &properties,
		}
		tiles = append(tiles, tile)
	}
	return tiles, nil
}

// getHitbox calculates the hitbox for a tile based on its object group.
func getHitbox(tiledTile *tiledTile, width float64, height float64) Rect {
	// Tiled allows a custom hitbox to be defined via a single object in the object group.
	if tiledTile.ObjectGroup.Objects != nil && len(tiledTile.ObjectGroup.Objects) > 0 {
		obj := tiledTile.ObjectGroup.Objects[0]
		return Rect{
			X:      obj.X,
			Y:      obj.Y,
			Width:  obj.Width,
			Height: obj.Height,
		}
	}

	// If no custom hitbox, use the full tile dimensions.
	return Rect{X: 0, Y: 0, Width: width, Height: height}
}
