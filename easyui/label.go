package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Label represents a static text display component.
type Label struct {
	component                         // Embed the base component for position and children
	Text        string                // The text to display
	uiGenerator *BareBonesUiGenerator // Reference to the generator for image regeneration
	idleImg     *ebiten.Image         // The pre-rendered image of the label text
}

// NewLabel creates a new Label instance.
// This constructor is typically called by the BareBonesUiGenerator.
func NewLabel(x, y, width, height int, text string, uiGen *BareBonesUiGenerator, idleImage *ebiten.Image) *Label {
	l := &Label{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		Text:        text,
		uiGenerator: uiGen,
		idleImg:     idleImage,
	}
	return l
}

// SetText updates the label's text and regenerates its image.
func (l *Label) SetText(newText string) {
	l.Text = newText
	// Regenerate the label's image with the new text
	l.idleImg = l.uiGenerator.generateLabelImage(
		l.Bounds.Dx(), l.Bounds.Dy(),
		l.uiGenerator.theme.OnPrimaryColor, // Use OnPrimaryColor for label text
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
