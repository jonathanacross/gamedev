package main

import (
	"encoding/json"
	"os"
)

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers     []TilemapLayerJSON `json:"layers"`
	Properties []PropertiesJSON   `json:"properties"`
}

type PropertiesJSON struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewTilemapJson(filepath string) TilemapJSON {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var tilemapjson TilemapJSON
	err = json.Unmarshal(contents, &tilemapjson)
	if err != nil {
		panic(err)
	}

	return tilemapjson
}
