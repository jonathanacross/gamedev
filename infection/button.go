package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	bounds image.Rectangle

	idleImage    *ebiten.Image
	pressedImage *ebiten.Image
	onclick      func()
	isPressed    bool
}

func NewButton(x, y, width, height int, idleImage, pressedImage *ebiten.Image, onclick func()) *Button {
	bounds := image.Rectangle{
		Min: image.Point{x, y},
		Max: image.Point{x + width, y + height},
	}

	return &Button{
		bounds:       bounds,
		idleImage:    idleImage,
		pressedImage: pressedImage,
		onclick:      onclick,
		isPressed:    false,
	}
}

func Contains(r image.Rectangle, x, y int) bool {
	return r.Min.X <= x && x <= r.Max.X && r.Min.Y <= y && y <= r.Max.Y
}

func (b *Button) Update() {
	cx, cy := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.isPressed = Contains(b.bounds, cx, cy)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if b.isPressed && Contains(b.bounds, cx, cy) {
			b.isPressed = false
			b.onclick()
		}
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.bounds.Min.X), float64(b.bounds.Min.Y))
	if b.isPressed {
		screen.DrawImage(b.pressedImage, op)
	} else {
		screen.DrawImage(b.idleImage, op)
	}
}

func GenerateButtonImage(width, height int, buttonText string, buttonColor, textColor color.RGBA) *ebiten.Image {
	strokewidth := 7

	img := ebiten.NewImage(width, height)
	img.Fill(buttonColor)
	vector.StrokeRect(img, 0, 0, float32(width), float32(height), float32(strokewidth), textColor, false)

	face := text.GoTextFace{
		Source: DisplayFont,
		Size:   float64(16),
	}
	_, textHeight := text.Measure(buttonText, &face, 1.0)
	op := &text.DrawOptions{}
	fontSize := float64(16)
	op.GeoM.Translate(0.5*float64(width), 0.3*float64(textHeight))
	op.ColorScale.ScaleWithColor(textColor)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(img, buttonText, &face, op)

	return img
}
