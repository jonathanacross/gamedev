package main

import (
	"encoding/json"
	"os"
)

type PropertiesJSON struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// IntValue attempts to convert the property's value to an int.
func (p *PropertiesJSON) IntValue() (int, bool) {
	if v, ok := p.Value.(float64); ok {
		return int(v), true
	}
	return 0, false
}

// BoolValue attempts to convert the property's value to a bool.
func (p *PropertiesJSON) BoolValue() (bool, bool) {
	if v, ok := p.Value.(bool); ok {
		return v, true
	}
	// Tiled can sometimes export booleans as 0 or 1.
	if v, ok := p.Value.(float64); ok {
		return v == 1, true
	}
	return false, false
}

func getStringProperty(properties []PropertiesJSON, name string) (string, bool) {
	for _, prop := range properties {
		if prop.Name == name {
			if v, ok := prop.Value.(string); ok {
				return v, true
			}
		}
	}
	return "", false
}

func getBoolProperty(properties []PropertiesJSON, name string) (bool, bool) {
	for _, prop := range properties {
		if prop.Name == name {
			if v, ok := prop.BoolValue(); ok {
				return v, true
			}
		}
	}
	return false, false
}

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

type TilesetTileJSON struct {
	ID          int                    `json:"id"`
	ObjectGroup TilesetObjectGroupJSON `json:"objectgroup"`
	Properties  []PropertiesJSON       `json:"properties"`
}

type TilesetDataJSON struct {
	Columns     int               `json:"columns"`
	Image       string            `json:"image"`
	Imageheight int               `json:"imageheight"`
	Imagewidth  int               `json:"imagewidth"`
	Margin      int               `json:"margin"`
	Name        string            `json:"name"`
	Spacing     int               `json:"spacing"`
	Tilecount   int               `json:"tilecount"`
	Tileheight  int               `json:"tileheight"`
	Tiles       []TilesetTileJSON `json:"tiles"`
	Tilewidth   int               `json:"tilewidth"`
	Type        string            `json:"type"`
}

type ObjectJSON struct {
	Height     float64          `json:"height"`
	ID         int              `json:"id"`
	Gid        int              `json:"gid"`
	Name       string           `json:"name"`
	Properties []PropertiesJSON `json:"properties"`
	Rotation   int              `json:"rotation"`
	Type       string           `json:"type"`
	Visible    bool             `json:"visible"`
	Width      float64          `json:"width"`
	X          float64          `json:"x"`
	Y          float64          `json:"y"`
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
