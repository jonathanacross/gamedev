package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PropertiesJSON struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// IntValue decodes the Value field into an integer based on the Type field.
// It returns the integer value and a boolean indicating success.
func (p *PropertiesJSON) IntValue() (int, bool) {
	if p.Type == "int" {
		// Use a type switch to safely get the integer value from the interface{}.
		switch v := p.Value.(type) {
		case float64:
			// JSON numbers are often unmarshaled as float64.
			return int(v), true
		case int:
			return v, true
		default:
			fmt.Printf("Warning: Property '%s' with type '%s' has an unexpected value type: %T\n", p.Name, p.Type, v)
			return 0, false
		}
	}
	return 0, false
}

type ObjectJSON struct {
	Height     float64          `json:"height"`
	Id         int              `json:"id"`
	Name       string           `json:"name"`
	Properties []PropertiesJSON `json:"properties"`
	Type       string           `json:"type"`
	Width      float64          `json:"width"`
	X          float64          `json:"x"`
	Y          float64          `json:"y"`
}

type TilemapLayerJSON struct {
	Data    []int        `json:"data"`
	Width   int          `json:"width"`
	Height  int          `json:"height"`
	Type    string       `json:"type"`
	Objects []ObjectJSON `json:"objects"`
}

type TilemapJSON struct {
	Layers     []TilemapLayerJSON `json:"layers"`
	Properties []PropertiesJSON   `json:"properties"`
	Width      int                `json:"width"`
	Height     int                `json:"height"`
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
