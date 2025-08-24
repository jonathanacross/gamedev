package main

import (
	"bytes"
	"embed"
	"image"
	"path/filepath"
	"strconv"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed assets/*
var assets embed.FS

var TileSet = loadImage("assets/tileset.png")
var PlayerSprite = loadImage("assets/player.png")
var Music = loadSound("assets/bach-prelude.mp3")
var Levels = loadLevels("assets/levels")

// You will also need a global map to store the loaded levels once the game is initialized
var LoadedLevels = make(map[int]*Level)

func loadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func loadSound(name string) *mp3.Stream {
	content, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	soundStream, err := mp3.DecodeWithoutResampling(bytes.NewReader(content))
	if err != nil {
		panic(err)
	}

	return soundStream
}

func loadLevels(dir string) map[int]TilemapJSON {
	levels := make(map[int]TilemapJSON)
	dirEntries, err := assets.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			// Extract level number from the filename (e.g., "level1.json" -> 1)
			name := entry.Name()
			if len(name) > len("level.json") {
				numStr := name[len("level") : len(name)-len(".json")]
				if num, err := strconv.Atoi(numStr); err == nil {
					// Read the file and parse the JSON
					filePath := filepath.Join(dir, name)
					levels[num] = NewTilemapJson(filePath)
				}
			}
		}
	}

	return levels
}
