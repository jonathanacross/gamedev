package main

import (
	"embed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"

	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var PlayerSprite = mustLoadImage("assets/playerShip1_red.png")
var BulletSprite = mustLoadImage("assets/laserBlue16.png")
var MeteorSprites = mustLoadImages("assets/meteors/*.png")
var ScoreFont = mustLoadFont("assets/kenvector_future.ttf")

func mustLoadImage(name string) *ebiten.Image {
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

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}
