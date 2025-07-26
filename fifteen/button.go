package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Button struct {
	bounds Rect

	idleImage    *ebiten.Image
	pressedImage *ebiten.Image
	onclick      func()
	isPressed    bool
}

func NewButton(x, y, width, height int, onclick func()) *Button {
	bounds := Rect{
		X:      float64(x),
		Y:      float64(y),
		Width:  float64(width),
		Height: float64(height),
	}

	idleColor := color.RGBA{250, 218, 155, 255}
	pressedColor := color.RGBA{161, 123, 47, 255}
	textColor := color.RGBA{77, 54, 8, 255}
	idle := generateButtonImage(width, height, "Shuffle", idleColor, textColor)
	pressed := generateButtonImage(width, height, "Shuffle", pressedColor, textColor)

	return &Button{
		bounds:       bounds,
		idleImage:    idle,
		pressedImage: pressed,
		onclick:      onclick,
		isPressed:    false,
	}
}

func (b *Button) Update() {
	cx, cy := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.isPressed = b.bounds.Contains(float64(cx), float64(cy))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if b.isPressed && b.bounds.Contains(float64(cx), float64(cy)) {
			b.isPressed = false
			b.onclick()
		}
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.bounds.X, b.bounds.Y)
	if b.isPressed {
		screen.DrawImage(b.pressedImage, op)
	} else {
		screen.DrawImage(b.idleImage, op)
	}
}

func generateButtonImage(width, height int, buttonText string, buttonColor, textColor color.RGBA) *ebiten.Image {
	strokewidth := 7

	img := ebiten.NewImage(width, height)
	img.Fill(buttonColor)
	vector.StrokeRect(img, 0, 0, float32(width), float32(height), float32(strokewidth), textColor, false)

	textBounds, _ := font.BoundString(ButtonFont, buttonText)
	textWidth := int((textBounds.Max.X - textBounds.Min.X) / 64)
	textHeight := int((textBounds.Max.Y - textBounds.Min.Y) / 64)
	text.Draw(img, "Shuffle", ButtonFont, (width-textWidth)/2, (height+textHeight)/2, textColor)

	return img
}
