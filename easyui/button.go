package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Button represents a clickable UI button.
type Button struct {
	interactiveComponent        // Embed the new interactive component
	onClick              func() // Specific to Button
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// Update calls the embedded interactiveComponent's Update method.
func (b *Button) Update() {
	b.interactiveComponent.Update()
}

// Draw draws the button's current state image to the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Bounds.Min.X), float64(b.Bounds.Min.Y))
	screen.DrawImage(b.GetCurrentStateImage(), op)
}

// HandlePress calls the embedded interactiveComponent's HandlePress method.
func (b *Button) HandlePress() {
	b.interactiveComponent.HandlePress()
}

// HandleRelease calls the embedded interactiveComponent's HandleRelease method.
func (b *Button) HandleRelease() {
	b.interactiveComponent.HandleRelease()
}

// HandleClick calls the specific button's onClick handler.
func (b *Button) HandleClick() {
	// Only trigger onClick if not disabled. State is already managed by HandleRelease.
	if b.state != ButtonDisabled && b.onClick != nil {
		b.onClick()
	}
}
