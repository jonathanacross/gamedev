package main

import (
	"embed"

	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var TerrainTileset = loadImage("assets/terrain.png")
var WallBlobTileset = loadImage("assets/walls_blob.png")
var PlayerIdleSpritesImage = loadImage("assets/player_idle.png")
var PlayerWalkSpritesImage = loadImage("assets/player_walk.png")
var PlayerDeathSpritesImage = loadImage("assets/player_death.png")

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
