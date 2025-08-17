package main

import (
	"log"

	"image" // Required for image.Rectangle, image.Point

	"github.com/hajimehoshi/ebiten/v2"
)

// Label represents a static text display component.
type Label struct {
	component               // Embed the base component for position and children
	Text      string        // The text to display
	renderer  UiRenderer    // Changed to UiRenderer interface
	idleImg   *ebiten.Image // The pre-rendered image of the label text
}

// NewLabel creates a new Label instance.
// It is now a standalone function.
func NewLabel(x, y, width, height int, text string, renderer UiRenderer) *Label {
	// Labels are static, so only an idle image is needed.
	labelImage := renderer.GenerateLabelImage(width, height, text) // Text color from theme
	l := &Label{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		Text:     text,
		renderer: renderer, // Store the renderer
		idleImg:  labelImage,
	}
	return l
}

// SetText updates the label's text and regenerates its image.
func (l *Label) SetText(newText string) {
	l.Text = newText
	// Regenerate the label's image with the new text using the renderer
	l.idleImg = l.renderer.GenerateLabelImage(
		l.Bounds.Dx(), l.Bounds.Dy(),
		l.Text,
	)
}

// Update for Label is a no-op as it has no interactive logic.
func (l *Label) Update() {
	// Labels are static, so no update logic is needed here.
}

// Draw draws the label's pre-rendered image to the screen.
func (l *Label) Draw(screen *ebiten.Image) {
	if l.idleImg == nil {
		log.Printf("Label '%s': WARNING: idleImg is nil! Cannot draw label.", l.Text)
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(l.Bounds.Min.X), float64(l.Bounds.Min.Y))
	screen.DrawImage(l.idleImg, op)
}

// HandlePress, HandleRelease, HandleClick are no-ops for a static Label.
func (l *Label) HandlePress()   {}
func (l *Label) HandleRelease() {}
func (l *Label) HandleClick()   {}
