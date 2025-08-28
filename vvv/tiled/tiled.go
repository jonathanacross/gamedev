package tiled

import (
	"encoding/json"
	"os"
)

type ObjectLayerJSON struct {
	Draworder string       `json:"draworder"`
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Objects   []ObjectJSON `json:"objects"`
	Opacity   float64      `json:"opacity"`
	Type      string       `json:"type"`
	Visible   bool         `json:"visible"`
	X         int          `json:"x"`
	Y         int          `json:"y"`
}

type TilemapLayerJSON struct {
	Data    []int        `json:"data"`
	Width   int          `json:"width"`
	Height  int          `json:"height"`
	Objects []ObjectJSON `json:"objects"`
	Type    string       `json:"type"`
	Name    string       `json:"name"`
}

type TilesetJSON struct {
	Firstgid int    `json:"firstgid"`
	Source   string `json:"source"`
}

type CollisionRectangleJSON struct {
	Height float64 `json:"height"`
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
}

type TilesetObjectGroupJSON struct {
	Name    string                   `json:"name"`
	Objects []CollisionRectangleJSON `json:"objects"`
}

type LevelJSON struct {
	Layers   []TilemapLayerJSON `json:"layers"`
	Height   int                `json:"height"`
	Width    int                `json:"width"`
	Tilesets []TilesetJSON      `json:"tilesets"`
}

func NewLevelJSON(filepath string) LevelJSON {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var leveljson LevelJSON
	err = json.Unmarshal(contents, &leveljson)
	if err != nil {
		panic(err)
	}

	return leveljson
}

func NewTilesetJSON(filepath string) TilesetDataJSON {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var tilesetjson TilesetDataJSON
	err = json.Unmarshal(contents, &tilesetjson)
	if err != nil {
		panic(err)
	}

	return tilesetjson
}

func (leveljson LevelJSON) FindTileset(gid int) *TilesetJSON {
	// Tiled stores tilesets in descending order of firstgid.
	// So we can iterate backwards to find the correct tileset.
	for i := len(leveljson.Tilesets) - 1; i >= 0; i-- {
		ts := leveljson.Tilesets[i]
		if gid >= ts.Firstgid {
			return &ts
		}
	}
	return nil
}

// FindTile takes a local tile ID (not a global GID) and returns the tile.
func (td TilesetDataJSON) FindTile(id int) *TilesetTileJSON {
	for _, tile := range td.Tiles {
		if tile.ID == id {
			return &tile
		}
	}
	return nil
}
