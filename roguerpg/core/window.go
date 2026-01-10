package core

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Window struct {
	Rect  Rect
	Image *ebiten.Image
}

func newWindowImage(rect Rect, backgroundColor color.Color, frameColor color.Color) *ebiten.Image {
	cornerRadius := 2.0
	lineWidth := 2.0

	dc := gg.NewContext(int(rect.Width()), int(rect.Height()))

	// Background
	dc.SetColor(backgroundColor)
	DrawBeveledRectangle(dc, lineWidth/2, lineWidth/2, rect.Width()-lineWidth, rect.Height()-lineWidth, cornerRadius)
	dc.FillPreserve()

	// Border
	dc.SetColor(frameColor)
	dc.SetLineWidth(lineWidth)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

func NewWindow(rect Rect) *Window {
	backgroundColor := color.RGBA{R: 73, G: 65, B: 130, A: 255}
	frameColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	return &Window{
		Rect:  rect,
		Image: newWindowImage(rect, backgroundColor, frameColor),
	}
}

func (w *Window) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.Rect.Left, w.Rect.Top)
	screen.DrawImage(w.Image, op)
}
