package core

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonStatus int

const (
	ButtonDisabled ButtonStatus = iota
	ButtonEnabled
	ButtonHovered
)

type Button struct {
	Rect          Rect
	Status        ButtonStatus
	Image         *ebiten.Image
	HoverImage    *ebiten.Image
	DisabledImage *ebiten.Image
}

// This is better than rounded rectangles for small corner radii;
// there is some artifacting with gg.DrawRoundedRectangle.
func DrawBeveledRectangle(dc *gg.Context, x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.NewSubPath()
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.LineTo(x3, y1)
	dc.LineTo(x3, y2)
	dc.LineTo(x2, y3)
	dc.LineTo(x1, y3)
	dc.LineTo(x0, y2)
	dc.LineTo(x0, y1)
	dc.ClosePath()
}

func newButtonImage(rect Rect, backgroundColor color.Color, frameColor color.Color) *ebiten.Image {
	cornerRadius := 2.0
	lineWidth := 1.0

	dc := gg.NewContext(int(rect.Width()), int(rect.Height()))

	// Background
	// dc.SetColor(backgroundColor)
	DrawBeveledRectangle(dc, lineWidth/2, lineWidth/2, rect.Width()-lineWidth, rect.Height()-lineWidth, cornerRadius)
	// dc.FillPreserve()

	// Border
	dc.SetColor(frameColor)
	dc.SetLineWidth(lineWidth)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// TODO: simplify drawing to only have one image, if that's all that is needed
func NewButton(rect Rect) *Button {
	transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	frameColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	return &Button{
		Rect:          rect,
		Status:        ButtonEnabled,
		Image:         newButtonImage(rect, transparent, transparent),
		HoverImage:    newButtonImage(rect, transparent, frameColor),
		DisabledImage: newButtonImage(rect, transparent, transparent),
	}
}

func (w *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.Rect.Left, w.Rect.Top)
	switch w.Status {
	case ButtonDisabled:
		screen.DrawImage(w.DisabledImage, op)
	case ButtonHovered:
		screen.DrawImage(w.HoverImage, op)
	default:
		screen.DrawImage(w.Image, op)
	}
}
