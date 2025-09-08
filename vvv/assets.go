package main

import (
	"bytes"
	"embed"

	"image"
	"path/filepath"
	"strconv"
	"vvv/tiled"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/*
var assets embed.FS

var TileSetImage = loadImage("assets/images/tileset.png")
var PlayerSprite = loadImage("assets/images/player.png")
var MonsterSprite = loadImage("assets/images/helicopterguy.png")
var CheckpointSprite = loadImage("assets/images/checkpoint.png")
var BreakingFloorSprite = loadImage("assets/images/breakingfloor.png")
var StartScreen = loadImage("assets/images/titlescreen.png")
var WinScreen = loadImage("assets/images/winscreen.png")
var Levels = loadLevels("assets/levels")
var Music = loadSound("assets/sounds/bach-prelude.mp3")
var ArcadeFaceSource = loadFaceSource("assets/fonts/pressstart2p.ttf")

// Store the loaded levels once the game is initialized
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

func loadFaceSource(name string) *text.GoTextFaceSource {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	face, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		panic(err)
	}
	return face
}

func loadLevels(dir string) map[int]*tiled.Map {
	levels := make(map[int]*tiled.Map)
	dirEntries, err := assets.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	ebitenImageConverter := func(img image.Image) (tiled.ImageProvider, error) {
		return ebiten.NewImageFromImage(img), nil
	}

	loader := tiled.NewFsLoaderWithImageConverter(assets, ebitenImageConverter)
	for _, entry := range dirEntries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()
			if len(name) > len("level.json") && name[:len("level")] == "level" {
				numStr := name[len("level") : len(name)-len(".json")]
				if num, err := strconv.Atoi(numStr); err == nil {
					filePath := filepath.Join(dir, name)
					level, err := loader.LoadMap(filePath)
					if err != nil {
						panic(err)
					}
					levels[num] = level
				}
			}
		}
	}

	return levels
}
