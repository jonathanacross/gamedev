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

// var TitlePageImage = loadImage("assets/title.png")
var SquareImage = loadImage("assets/square.png")

// Characters
var BiffCharImage = loadImage("assets/biff.png")
var DrQCharImage = loadImage("assets/drq_teal.png")
var ElaraCharImage = loadImage("assets/elara.png")
var EricaCharImage = loadImage("assets/erica.png")
var HeatherGCharImage = loadImage("assets/heather_green.png")
var HeatherVCharImage = loadImage("assets/heather_violet.png")
var MikeGCharImage = loadImage("assets/mike_green.png")
var MikeVCharImage = loadImage("assets/mike_violet.png")
var MiloCharImage = loadImage("assets/milo.png")
var SaraCharImage = loadImage("assets/sara.png")

var MusicBytes = loadSoundBytes("assets/8bit-canon.mp3")

//var MainFaceSource = loadFaceSource("assets/ByteBounce.ttf")

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
