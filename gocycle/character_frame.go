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
	Mood        CharacterMood
	State       CharacterState
	SpriteSheet *SpriteSheet
	CharData    CharData
}

func (cf *CharacterFrame) Draw(screen *ebiten.Image) {
	// Fill with background color
	var bgColor color.Color
	var frameColor color.Color
	if cf.State == StateSelected {
		bgColor = cf.CharData.SelectedColor
		frameColor = cf.CharData.FrameColor
	} else {
		bgColor = color.Black
		frameColor = cf.CharData.SelectedColor
	}
	vector.DrawFilledRect(screen,
		float32(cf.X), float32(cf.Y),
		float32(cf.hitbox.Width()), float32(cf.hitbox.Height()),
		bgColor, false)

	// Draw the character
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cf.X, cf.Y)

	// Get the source rcdt from the spritesheet, but
	// update based on the hitbox in case we're showing small
	// portraits.
	srcRect := cf.SpriteSheet.Rect(int(cf.Mood))
	srcRect.Max.Y = srcRect.Min.Y + int(cf.hitbox.Height())

	currImage := cf.image.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)

	// Draw the frame
	vector.StrokeRect(screen,
		float32(cf.X), float32(cf.Y),
		float32(cf.hitbox.Width()), float32(cf.hitbox.Height()),
		1, frameColor, false)
}

func NewCharacterFrame(CharacterIdx int, x, y float64, mood CharacterMood, smallPortrait bool) *CharacterFrame {
	charData := Characters[CharacterIdx]
	width := CharPortraitWidth
	height := CharPortraitBigHeight
	if smallPortrait {
		height = CharPortraitSmallHeight
	}

	frame := &CharacterFrame{
		BaseSprite: BaseSprite{
			Location: Location{
				X: x,
				Y: y,
			},
			image:   charData.Image,
			srcRect: charData.Image.Bounds(),
			hitbox:  Rect{left: 0, top: 0, right: float64(width), bottom: float64(height)},
		},
		SpriteSheet: NewSpriteSheet(width, CharPortraitBigHeight, 3, 1),
		Mood:        mood,
		State:       StateUnselected,
		CharData:    charData,
	}

	return frame
}
