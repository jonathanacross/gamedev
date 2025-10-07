package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CharacterMood int

const (
	CharacterNeutral CharacterMood = iota
	CharacterHappy
	CharacterSad
)

type CharacterState int

const (
	StateUnselected CharacterState = iota
	StateSelected
)

type CharacterFrame struct {
	BaseSprite
	Mood            CharacterMood
	State           CharacterState
	UnselectedColor color.Color
	SelectedColor   color.Color
	FrameColor      color.Color
	SpriteSheet     *SpriteSheet
}

type CharData struct {
	Name          string
	Image         *ebiten.Image
	SelectedColor color.Color
	FrameColor    color.Color
}

const (
	CharacterMilo int = iota
	CharacterSara
	CharacterDrQ
	CharacterErica
	CharacterBiff
	CharacterElara
	CharacterMikeG
	CharacterMikeV
	CharacterHeatherG
	CharacterHeatherV
)

var Characters []CharData = loadCharData()

func (cf *CharacterFrame) Draw(screen *ebiten.Image) {
	// Fill with background color
	var bgColor color.Color
	var frameColor color.Color
	if cf.State == StateSelected {
		bgColor = cf.SelectedColor
		frameColor = cf.FrameColor
	} else {
		bgColor = cf.UnselectedColor
		frameColor = cf.SelectedColor
	}
	vector.DrawFilledRect(screen,
		float32(cf.X), float32(cf.Y),
		float32(cf.hitbox.Width()), float32(cf.hitbox.Height()),
		bgColor, false)

	// Draw the character
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cf.X, cf.Y)

	srcRect := cf.SpriteSheet.Rect(int(cf.Mood))
	currImage := cf.image.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)

	// Draw the frame
	vector.StrokeRect(screen,
		float32(cf.X), float32(cf.Y),
		float32(cf.hitbox.Width()), float32(cf.hitbox.Height()),
		1, frameColor, false)
}

func NewCharacterFrame(CharacterIdx int, x, y float64, mood CharacterMood) *CharacterFrame {
	charData := Characters[CharacterIdx]

	frame := &CharacterFrame{
		BaseSprite: BaseSprite{
			Location: Location{
				X: x,
				Y: y,
			},
			image:   charData.Image,
			srcRect: charData.Image.Bounds(),
			hitbox:  Rect{left: 0, top: 0, right: 64, bottom: 80},
		},
		SpriteSheet:     NewSpriteSheet(64, 80, 3, 1),
		Mood:            mood,
		State:           StateUnselected,
		UnselectedColor: color.RGBA{0, 0, 0, 255},
		SelectedColor:   charData.SelectedColor,
		FrameColor:      charData.FrameColor,
	}

	return frame
}

func loadCharData() []CharData {
	return []CharData{
		{
			Name:          "Milo",
			Image:         MiloCharImage,
			SelectedColor: color.RGBA{22, 48, 83, 255},
			FrameColor:    color.RGBA{3, 166, 224, 255},
		},
		{
			Name:          "Sara",
			Image:         SaraCharImage,
			SelectedColor: color.RGBA{146, 132, 51, 255},
			FrameColor:    color.RGBA{248, 243, 79, 255},
		},
		{
			Name:          "Dr. Q",
			Image:         DrQCharImage,
			SelectedColor: color.RGBA{92, 92, 92, 255},
			FrameColor:    color.RGBA{220, 220, 220, 255},
		},
		{
			Name:          "Erica",
			Image:         EricaCharImage,
			SelectedColor: color.RGBA{104, 9, 13, 255},
			FrameColor:    color.RGBA{231, 64, 71, 255},
		},
		{
			Name:          "Biff",
			Image:         BiffCharImage,
			SelectedColor: color.RGBA{99, 26, 3, 255},
			FrameColor:    color.RGBA{238, 156, 50, 255},
		},
		{
			Name:          "Elara",
			Image:         ElaraCharImage,
			SelectedColor: color.RGBA{58, 59, 94, 255},
			FrameColor:    color.RGBA{121, 121, 203, 255},
		},
		{
			Name:          "Mike Green",
			Image:         MikeGCharImage,
			SelectedColor: color.RGBA{20, 104, 20, 255},
			FrameColor:    color.RGBA{156, 224, 42, 255},
		},
		{
			Name:          "Mike Violet",
			Image:         MikeVCharImage,
			SelectedColor: color.RGBA{68, 21, 51, 255},
			FrameColor:    color.RGBA{193, 92, 153, 255},
		},
		{
			Name:          "Heather Green",
			Image:         HeatherGCharImage,
			SelectedColor: color.RGBA{20, 104, 20, 255},
			FrameColor:    color.RGBA{156, 224, 42, 255},
		},
		{
			Name:          "Heather Violet",
			Image:         HeatherVCharImage,
			SelectedColor: color.RGBA{68, 21, 51, 255},
			FrameColor:    color.RGBA{193, 92, 153, 255},
		},
	}
}
