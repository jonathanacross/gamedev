package tiled

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

// A mock file system implementation for testing.
type mockFS struct {
	files map[string][]byte
}

func (m *mockFS) Open(name string) (fs.File, error) {
	data, ok := m.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &mockFile{data: data}, nil
}

type mockFile struct {
	data []byte
	pos  int
}

func (m *mockFile) Stat() (fs.FileInfo, error) { return nil, nil }
func (m *mockFile) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	n = copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}
func (m *mockFile) Close() error { return nil }

func newMockFS() *mockFS {
	// Create a simple 16x16 PNG image in memory.
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})

	// Encode the image to PNG bytes to create a valid image file.
	var buf bytes.Buffer
	png.Encode(&buf, img)
	imageBytes := buf.Bytes()

	return &mockFS{
		files: map[string][]byte{
			"assets/levels/level4.json": []byte(`
				{
					"height": 15,
					"layers": [
						{
							"data": [1, 2, 3],
							"height": 15,
							"name": "Tile Layer 1",
							"type": "tilelayer",
							"width": 20
						},
						{
							"name": "Object Layer 1",
							"objects": [
								{
									"gid": 25,
									"height": 16,
									"id": 1,
									"name": "Object 1",
									"properties": [
										{ "name": "damage", "type": "int", "value": 10 }
									],
									"type": "Spikes",
									"width": 16,
									"x": 100,
									"y": 100
								}
							],
							"type": "objectgroup"
						}
					],
					"name": "level4.json",
					"tilesets": [
						{
							"firstgid": 1,
							"source": "../tilesets/tileset.json"
						}
					],
					"tileheight": 16,
					"tilewidth": 16,
					"width": 20
				}
			`),
			"assets/tilesets/tileset.json": []byte(`
				{
					"columns": 5,
					"image": "../images/tileset.png",
					"imageheight": 128,
					"imagewidth": 80,
					"name": "tileset",
					"tilecount": 3,
					"tileheight": 16,
					"tilewidth": 16,
					"tiles": [
						{
							"id": 0,
							"properties": [
								{ "name": "solid", "type": "bool", "value": true }
							]
						},
						{
							"id": 1,
							"properties": [
								{ "name": "solid", "type": "bool", "value": true }
							]
						}
					]
				}
			`),
			"assets/images/tileset.png": imageBytes,
		},
	}
}

func TestLoadMap(t *testing.T) {
	t.Run("Successfully load a map", func(t *testing.T) {
		mockFS := newMockFS()
		loader := NewFsLoader(mockFS)

		gameMap, err := loader.LoadMap("assets/levels/level4.json")
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		// Verify map properties
		if gameMap.Name != "level4.json" {
			t.Errorf("Expected map name 'level4.json', got '%s'", gameMap.Name)
		}
		if gameMap.WidthInTiles != 20 || gameMap.HeightInTiles != 15 {
			t.Errorf("Expected map dimensions 20x15, got %dx%d", gameMap.WidthInTiles, gameMap.HeightInTiles)
		}
		if len(gameMap.Layers) != 2 {
			t.Errorf("Expected 2 layers, got %d", len(gameMap.Layers))
		}
		if len(gameMap.Tiles) != 3 {
			t.Errorf("Expected 3 tiles, got %d", len(gameMap.Tiles))
		}

		// Verify layer data
		tileLayer := gameMap.Layers[0]
		if tileLayer.Name != "Tile Layer 1" || tileLayer.Type != "tilelayer" {
			t.Errorf("Expected tile layer, got name '%s' and type '%s'", tileLayer.Name, tileLayer.Type)
		}
		if !reflect.DeepEqual(tileLayer.TileIds, []int{1, 2, 3}) {
			t.Errorf("Expected tile data [1, 2, 3], got %v", tileLayer.TileIds)
		}

		objectLayer := gameMap.Layers[1]
		if objectLayer.Name != "Object Layer 1" || objectLayer.Type != "objectgroup" {
			t.Errorf("Expected object group layer, got name '%s' and type '%s'", objectLayer.Name, objectLayer.Type)
		}
		if len(objectLayer.Objects) != 1 {
			t.Fatalf("Expected 1 object, got %d", len(objectLayer.Objects))
		}

		// Verify object properties
		obj := objectLayer.Objects[0]
		if obj.Name != "Object 1" {
			t.Errorf("Expected object name 'Object 1', got '%s'", obj.Name)
		}
		if obj.Type != "Spikes" {
			t.Errorf("Expected object type 'Spikes', got '%s'", obj.Type)
		}
		if obj.Location.X != 100 || obj.Location.Y != 100 {
			t.Errorf("Expected object location {100, 100}, got {%v, %v}", obj.Location.X, obj.Location.Y)
		}

		damage, err := obj.Properties.GetPropertyInt("damage")
		if err != nil {
			t.Errorf("Expected damage property, but got an error: %v", err)
		}
		if damage != 10 {
			t.Errorf("Expected damage property with value 10, got %d", damage)
		}
	})

	t.Run("Map file not found", func(t *testing.T) {
		mockFS := newMockFS()
		loader := NewFsLoader(mockFS)
		_, err := loader.LoadMap("assets/levels/non-existent.json")
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}
		if !strings.Contains(err.Error(), "failed to load map file") {
			t.Errorf("Expected 'file not found' error, got: %v", err)
		}
	})

	t.Run("Invalid map JSON", func(t *testing.T) {
		mockFS := newMockFS()
		mockFS.files["assets/levels/invalid.json"] = []byte(`{"invalid"}`)
		loader := NewFsLoader(mockFS)
		_, err := loader.LoadMap("assets/levels/invalid.json")
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}
		if !strings.Contains(err.Error(), "failed to parse map JSON") {
			t.Errorf("Expected 'JSON parse' error, got: %v", err)
		}
	})
}
