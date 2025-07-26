package main

import (
	"bytes"
	"embed"
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/*
var assets embed.FS

var Picture = loadImage("assets/animecat.png")
var AltPicture = loadImage("assets/red.png")
var Background = loadImage("assets/background.png")
var NumberFont = loadGoRegularFont(50)
var ButtonFont = loadGoRegularFont(24)
var WinSound = loadSound("assets/tada.mp3")

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

func loadGoRegularFont(size float64) font.Face {
	tt, err := opentype.Parse(goregular.TTF) // Load font data
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
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
