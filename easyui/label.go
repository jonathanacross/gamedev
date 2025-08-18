package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Label represents a static text display component.
type Label struct {
	component
	Text     string
	renderer UiRenderer
	idleImg  *ebiten.Image
}

// NewLabel creates a new Label instance.
func NewLabel(x, y, width, height int, text string, renderer UiRenderer) *Label {
	// Labels are static; only an idle image is needed.
	labelImage := renderer.GenerateLabelImage(width, height, text) // Text color from theme

	// Create the Label first, then pass its pointer as 'self'
	l := &Label{
		Text:     text,
		renderer: renderer,
		idleImg:  labelImage,
	}
	l.component = NewComponent(x, y, width, height, l) // Pass 'l' as self

	return l
}

// SetText updates the label's text and regenerates its image.
func (l *Label) SetText(newText string) {
	l.Text = newText
	// Regenerate the label's image with the new text.
	l.idleImg = l.renderer.GenerateLabelImage(
		l.Bounds.Dx(), l.Bounds.Dy(),
		l.Text,
	)
}

// Update for Label is a no-op as it has no interactive logic.
func (l *Label) Update() {}

// Draw draws the label's pre-rendered image to the screen using its absolute position.
func (l *Label) Draw(screen *ebiten.Image) {
	if l.idleImg == nil {
		log.Printf("Label '%s': WARNING: idleImg is nil! Cannot draw label.", l.Text)
		return
	}

	absX, absY := l.GetAbsolutePosition()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(l.idleImg, op)
}

// HandlePress is a no-op for a static Label.
func (l *Label) HandlePress() {}

// HandleRelease is a no-op for a static Label.
func (l *Label) HandleRelease() {}

// HandleClick is a no-op for a static Label.
func (l *Label) HandleClick() {}
