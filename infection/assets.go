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

// Go style
// var WhiteSquare = loadImage("assets/gogui-white-64x64.png")
// var BlackSquare = loadImage("assets/gogui-black-64x64.png")
// var Empty1Square = loadImage("assets/woodlight.png")
// var Empty2Square = loadImage("assets/wooddark.png")

// UI elements
// var ButtonImage = loadImage("assets/button.png")
// var ButtonPressedImage = loadImage("assets/buttonpressed.png")
var PlayPausePlayImage = loadImage("assets/playpause-play.png")
var PlayPausePauseImage = loadImage("assets/playpause-pause.png")
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
