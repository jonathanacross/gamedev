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

var WhiteSquare = loadImage("assets/yellowvirus.png")
var BlackSquare = loadImage("assets/redvirus.png")
var Empty1Square = loadImage("assets/empty1.png")
var Empty2Square = loadImage("assets/empty2.png")
var SpinnerImage = loadImage("assets/spinner.png")

// UI elements
var RightArrowIdleImage = loadImage("assets/rightarrowidle.png")
var RightArrowPressedImage = loadImage("assets/rightarrowpressed.png")
var LeftArrowIdleImage = loadImage("assets/leftarrowidle.png")
var LeftArrowPressedImage = loadImage("assets/leftarrowpressed.png")

var DisplayFont = loadFaceSource("assets/notosans-regular.ttf")

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
