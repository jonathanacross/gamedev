package main

import (
	"bytes"
	"embed"
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/*
var assets embed.FS

var PlayerImage = loadImage("assets/BirdSprite.png")
var BackgroundImage = loadImage("assets/background.png")
var TerrainImage = loadImage("assets/terrain.png")
var CoinImage = loadImage("assets/coin.png")
var OctoImage = loadImage("assets/octopus.png")
var BeeImage = loadImage("assets/bee.png")
var HeartImage = loadImage("assets/heart.png")
var ArcadeFaceSource = loadFaceSource("assets/pressstart2p.ttf")

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
