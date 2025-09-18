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

var FrogSpriteSheet = loadImage("assets/frog2.png")
var PlatformSprite = loadImage("assets/platform.png")
var BootSprite = loadImage("assets/boot.png")

// Load sound files as byte slices so they can be reused
var MusicBytes = loadSoundBytes("assets/8bit-canon.mp3")
var ClearSoundBytes = loadSoundBytes("assets/clear-sound.mp3")
var ErrorSoundBytes = loadSoundBytes("assets/error-sound.mp3")
var MainFaceSource = loadFaceSource("assets/ByteBounce.ttf")

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

func loadSoundBytes(name string) []byte {
	content, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return content
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
