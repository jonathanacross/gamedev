package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Button represents a clickable UI button.
type Button struct {
	interactiveComponent        // Embed the new interactive component
	Label                string // Store the button's text
	onClick              func()
	renderer             UiRenderer // Changed to UiRenderer interface
}

// NewButton creates a new Button instance with the specified dimensions, label, and renderer.
func NewButton(x, y, width, height int, label string, renderer UiRenderer) *Button {
	idle := renderer.GenerateButtonImage(width, height, label, ButtonIdle)
	pressed := renderer.GenerateButtonImage(width, height, label, ButtonPressed)
	hover := renderer.GenerateButtonImage(width, height, label, ButtonHover)
	disabled := renderer.GenerateButtonImage(width, height, label, ButtonDisabled)

	return &Button{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idle, pressed, hover, disabled),
		Label:                label,
		onClick:              nil,      // Click handler set separately
		renderer:             renderer, // Store the renderer
	}
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// SetText updates the button's text and regenerates its state images.
func (b *Button) SetText(newText string) {
	b.Label = newText
	// Regenerate all state images with the new text using the renderer
	b.idleImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonIdle)
	b.pressedImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonPressed)
	b.hoverImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonHover)
	b.disabledImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonDisabled)
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
