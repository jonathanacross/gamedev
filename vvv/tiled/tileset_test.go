package tiled

import (
	"image"
	"testing"
)

func TestGetTiles_CollectionTileset(t *testing.T) {
	// Mock a TilesetDataJSON object based on platforms.json
	tsj := TilesetDataJSON{
		Columns: 0,
		Name:    "platforms",
		Tiles: []TilesetTileJSON{
			{
				ID:          0,
				Image:       "platform-small.png",
				ImageHeight: 7,
				ImageWidth:  16,
				ObjectGroup: &ObjectGroupJSON{
					Objects: []ObjectJSON{
						{Height: 7, Width: 8, X: 2, Y: 0},
					},
				},
				Properties: []PropertiesJSON{
					{Name: "intproperty", Type: "int", Value: 8},
				},
				Type: "Platform",
			},
			{
				ID:          1,
				Image:       "platform-medium.png",
				ImageHeight: 7,
				ImageWidth:  32,
				Type:        "Platform",
			},
		},
		TileCount:  2,
		TileHeight: 7,
		TileWidth:  48,
		Type:       "tileset",
	}

	// Mock images for the tileset
	images := map[string]*image.Image{
		"platform-small.png":  new(image.Image),
		"platform-medium.png": new(image.Image),
	}

	tiles, err := GetTiles(&tsj, 1, images)

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if len(tiles) != 2 {
		t.Fatalf("expected 2 tiles, got %d", len(tiles))
	}

	// Test the first tile (ID 0)
	t.Run("Tile 0", func(t *testing.T) {
		tile := tiles[0]
		if tile.id != 1 {
			t.Errorf("expected ID 1, got %d", tile.id)
		}
		if tile.srcRect.Width != 16 || tile.srcRect.Height != 7 {
			t.Errorf("expected srcRect {0,0,16,7}, got %+v", tile.srcRect)
		}
		if tile.hitRect.X != 2 || tile.hitRect.Width != 8 {
			t.Errorf("expected hitbox {2,0,8,7}, got %+v", tile.hitRect)
		}
		if prop, _ := tile.properties.GetPropertyInt("intproperty"); prop != 8 {
			t.Errorf("expected intproperty to be 8, got %v", prop)
		}
	})

	// Test the second tile (ID 1)
	t.Run("Tile 1", func(t *testing.T) {
		tile := tiles[1]
		if tile.id != 2 {
			t.Errorf("expected ID 2, got %d", tile.id)
		}
		if tile.srcRect.Width != 32 || tile.srcRect.Height != 7 {
			t.Errorf("expected srcRect {0,0,32,7}, got %+v", tile.srcRect)
		}
		// The second tile has no custom hitbox, so it should have a default
		if tile.hitRect.Width != 32 || tile.hitRect.Height != 7 {
			t.Errorf("expected default hitbox {0,0,32,7}, got %+v", tile.hitRect)
		}
		if tile.properties != nil && len(*tile.properties) > 0 {
			t.Errorf("expected no properties, but got %d", len(*tile.properties))
		}
	})

	// Error case: Missing image
	t.Run("Missing Image Error", func(t *testing.T) {
		tsjWithMissingImage := tsj
		tsjWithMissingImage.Tiles[0].Image = "missing.png"
		_, err := GetTiles(&tsjWithMissingImage, 1, images)
		if err == nil {
			t.Errorf("expected error for missing image, but got nil")
		}
	})
}

func TestGetTiles_GridTileset(t *testing.T) {
	// Mock a TilesetDataJSON object based on tileset.json
	tsj := TilesetDataJSON{
		Columns:     2,
		Image:       "tileset.png",
		ImageHeight: 32,
		ImageWidth:  32,
		Name:        "tileset",
		TileCount:   4,
		TileHeight:  16,
		TileWidth:   16,
		Type:        "tileset",
		Tiles: []TilesetTileJSON{
			{ID: 0},
			{ID: 1, Properties: []PropertiesJSON{{Name: "solid", Type: "bool", Value: true}}},
			{ID: 2, ObjectGroup: &ObjectGroupJSON{Objects: []ObjectJSON{{Height: 9, Width: 16, X: 0, Y: 7}}}, Type: "Spikes"},
			{ID: 3},
		},
	}

	// Mock the single image for the grid tileset
	images := map[string]*image.Image{
		"tileset.png": new(image.Image),
	}

	tiles, err := GetTiles(&tsj, 1, images)

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if len(tiles) != len(tsj.Tiles) {
		t.Fatalf("expected %d tiles, got %d", len(tsj.Tiles), len(tiles))
	}

	// Test the first tile (ID 0)
	t.Run("Tile 0", func(t *testing.T) {
		tile := tiles[0]
		// Note that all the ids are shifted by 1 because GetTiles was called with gid=1
		if tile.id != 1 {
			t.Errorf("expected ID 1, got %d", tile.id)
		}
		if tile.srcRect.X != 0 || tile.srcRect.Y != 0 {
			t.Errorf("expected srcRect {0,0,16,16}, got %+v", tile.srcRect)
		}
		if tile.hitRect.Width != 16 || tile.hitRect.Height != 16 {
			t.Errorf("expected default hitbox {0,0,16,16}, got %+v", tile.hitRect)
		}
	})

	// Test the second tile (ID 1)
	t.Run("Tile 1", func(t *testing.T) {
		tile := tiles[1]
		if tile.id != 2 {
			t.Errorf("expected ID 2, got %d", tile.id)
		}
		if tile.srcRect.X != 16 || tile.srcRect.Y != 0 {
			t.Errorf("expected srcRect {16,0,16,16}, got %+v", tile.srcRect)
		}
	})

	// Test a tile with a custom hitbox (ID 30)
	t.Run("Tile 2 with custom hitbox", func(t *testing.T) {
		tile := tiles[2]
		if tile.id != 3 {
			t.Errorf("expected ID 3, got %d", tile.id)
		}
		// Based on the Tiled JSON, this tile is at column 0, row 1
		if tile.srcRect.X != 0 || tile.srcRect.Y != 16 {
			t.Errorf("expected srcRect {0,16,16,16}, got %+v", tile.srcRect)
		}
		if tile.hitRect.Y != 7 {
			t.Errorf("expected custom hitbox Y to be 7, got %v", tile.hitRect.Y)
		}
	})

	// Error case: Missing image
	t.Run("Missing Image Error", func(t *testing.T) {
		tsjWithMissingImage := tsj
		tsjWithMissingImage.Image = "missing.png"
		_, err := GetTiles(&tsjWithMissingImage, 1, images)
		if err == nil {
			t.Errorf("expected error for missing image, but got nil")
		}
	})
}
