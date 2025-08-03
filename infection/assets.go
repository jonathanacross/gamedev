package main

import (
	"embed"
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
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
