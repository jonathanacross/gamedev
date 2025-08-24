package main

import (
	"bytes"
	"embed"
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed assets/*
var assets embed.FS

var TileSet = loadImage("assets/tileset.png")
var PlayerSprite = loadImage("assets/player.png")
var LevelTilemapJSON = NewTilemapJson("assets/level1-1.json")
var Music = loadSound("assets/bach-prelude.mp3")

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
