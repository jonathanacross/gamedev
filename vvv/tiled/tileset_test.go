package tiled

import (
	"reflect"
	"testing"
)

// Implements the ImageProvider interface for testing.
type MockImage struct{}

var mockImage = &MockImage{}

func TestConvertTileset(t *testing.T) {
	// Test case for a sprite sheet tileset.
	t.Run("SpriteSheetTileset", func(t *testing.T) {
		tsData := &tiledTileset{
			Image:      "tileset.png",
			TileWidth:  16,
			TileHeight: 16,
			Columns:    5,
			Tiles: []tiledTile{
				{ID: 0, Properties: []tiledProperty{{Name: "solid", Type: "bool", Value: true}}},
				{ID: 1},
			},
			TileCount: 2,
		}

		firstGID := 10
		images := map[string]ImageProvider{"tileset.png": mockImage}
		tiles, err := ConvertTileset(tsData, images, firstGID)
		if err != nil {
			t.Fatalf("ConvertTileset failed: %v", err)
		}

		if len(tiles) != 2 {
			t.Fatalf("expected 2 tiles, got %d", len(tiles))
		}

		// Check first tile (ID 0)
		expectedTile0 := Tile{
			ID:         10,
			SrcRect:    Rect{X: 0, Y: 0, Width: 16, Height: 16},
			SrcImage:   mockImage,
			HitRect:    Rect{X: 0, Y: 0, Width: 16, Height: 16},
			Properties: &PropertySet{"solid": {Value: true}},
		}
		if !reflect.DeepEqual(tiles[0], expectedTile0) {
			t.Errorf("expected tile %+v, got %+v", expectedTile0, tiles[0])
		}

		// Check second tile (ID 1)
		expectedTile1 := Tile{
			ID:         11,
			SrcRect:    Rect{X: 16, Y: 0, Width: 16, Height: 16},
			SrcImage:   mockImage,
			HitRect:    Rect{X: 0, Y: 0, Width: 16, Height: 16},
			Properties: &PropertySet{},
		}
		if !reflect.DeepEqual(tiles[1], expectedTile1) {
			t.Errorf("expected tile %+v, got %+v", expectedTile1, tiles[1])
		}
	})

	// Test case for a collection tileset.
	t.Run("CollectionTileset", func(t *testing.T) {
		tsData := &tiledTileset{
			Image:      "",
			Columns:    0,
			TileWidth:  48,
			TileHeight: 7,
			Tiles: []tiledTile{
				{
					ID:          0,
					Image:       "platform-small.png",
					ImageWidth:  16,
					ImageHeight: 7,
					ObjectGroup: tiledObjectGroup{
						Objects: []tiledObject{
							{X: 2, Y: 0, Width: 8, Height: 7},
						},
					},
					Properties: []tiledProperty{
						{Name: "intproperty", Type: "int", Value: 8},
					},
				},
			},
			TileCount: 1,
		}

		firstGID := 41
		images := map[string]ImageProvider{
			"platform-small.png": mockImage,
		}
		tiles, err := ConvertTileset(tsData, images, firstGID)
		if err != nil {
			t.Fatalf("ConvertTileset failed: %v", err)
		}

		if len(tiles) != 1 {
			t.Fatalf("expected 1 tile, got %d", len(tiles))
		}

		// Check the tile
		expectedTile := Tile{
			ID:         41,
			SrcRect:    Rect{X: 0, Y: 0, Width: 16, Height: 7},
			SrcImage:   images["platform-small.png"],
			HitRect:    Rect{X: 2, Y: 0, Width: 8, Height: 7},
			Properties: &PropertySet{"intproperty": {Value: 8}},
		}
		if !reflect.DeepEqual(tiles[0], expectedTile) {
			t.Errorf("expected tile %+v, got %+v", expectedTile, tiles[0])
		}
	})
}

func TestGetHitbox(t *testing.T) {
	t.Run("CustomHitbox", func(t *testing.T) {
		tiledTile := &tiledTile{
			ObjectGroup: tiledObjectGroup{
				Objects: []tiledObject{
					{X: 5, Y: 5, Width: 10, Height: 10},
				},
			},
		}
		hitbox := getHitbox(tiledTile, 32, 32)
		expected := Rect{X: 5, Y: 5, Width: 10, Height: 10}
		if !reflect.DeepEqual(hitbox, expected) {
			t.Errorf("expected hitbox %+v, got %+v", expected, hitbox)
		}
	})

	t.Run("DefaultHitbox", func(t *testing.T) {
		tiledTile := &tiledTile{
			ObjectGroup: tiledObjectGroup{
				Objects: []tiledObject{},
			},
		}
		hitbox := getHitbox(tiledTile, 32, 32)
		expected := Rect{X: 0, Y: 0, Width: 32, Height: 32}
		if !reflect.DeepEqual(hitbox, expected) {
			t.Errorf("expected hitbox %+v, got %+v", expected, hitbox)
		}
	})
}
