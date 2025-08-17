package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	// Required for image.Rectangle, image.Point
)

// Label represents a static text display component.
// It now fully implements the Component interface.
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

	// Create the Label first, then pass its pointer as 'self'
	l := &Label{
		Text:     text,
		renderer: renderer, // Store the renderer
		idleImg:  labelImage,
	}
	l.component = NewComponent(x, y, width, height, l) // Pass 'l' as self

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
// This method fully implements Component.Update().
func (l *Label) Update() {
	// Labels are static, so no update logic is needed here.
}

// Draw draws the label's pre-rendered image to the screen using its absolute position.
// This method fully implements Component.Draw().
func (l *Label) Draw(screen *ebiten.Image) {
	if l.idleImg == nil {
		log.Printf("Label '%s': WARNING: idleImg is nil! Cannot draw label.", l.Text)
		return
	}

	absX, absY := l.GetAbsolutePosition() // Get absolute position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(absX), float64(absY)) // Translate by absolute position
	screen.DrawImage(l.idleImg, op)
}

// HandlePress is a no-op for a static Label.
// This method fully implements Component.HandlePress().
func (l *Label) HandlePress() {}

// HandleRelease is a no-op for a static Label.
// This method fully implements Component.HandleRelease().
func (l *Label) HandleRelease() {}

// HandleClick is a no-op for a static Label.
// This method fully implements Component.HandleClick().
func (l *Label) HandleClick() {}
