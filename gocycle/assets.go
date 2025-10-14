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
var ElaraSwimCharImage = loadImage("assets/elara_swim.png")
var HeatherGCharImage = loadImage("assets/heather_green.png")
var HeatherGSwimCharImage = loadImage("assets/heather_green_swim.png")
var HeatherVCharImage = loadImage("assets/heather_violet.png")
var HeatherVSwimCharImage = loadImage("assets/heather_violet_swim.png")
var EricaCharImage = loadImage("assets/erica.png")
var EricaSwimCharImage = loadImage("assets/erica_swim.png")
var MikeGCharImage = loadImage("assets/mike_green.png")
var MikeVCharImage = loadImage("assets/mike_violet.png")
var MiloCharImage = loadImage("assets/milo.png")
var SaraCharImage = loadImage("assets/sara.png")
var SaraSwimCharImage = loadImage("assets/sara_swim.png")

var MusicBytes = loadSoundBytes("assets/chiptune-grooving.mp3")

var MainFaceSource = loadFaceSource("assets/m5x7.ttf")

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
