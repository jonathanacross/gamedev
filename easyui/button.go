package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Button represents a clickable UI button.
type Button struct {
	interactiveComponent
	Label    string
	onClick  func()
	renderer UiRenderer
}

// NewButton creates a new Button instance with the specified dimensions, label, and renderer.
func NewButton(x, y, width, height int, label string, renderer UiRenderer) *Button {
	idle := renderer.GenerateButtonImage(width, height, label, ButtonIdle)
	pressed := renderer.GenerateButtonImage(width, height, label, ButtonPressed)
	hover := renderer.GenerateButtonImage(width, height, label, ButtonHover)
	disabled := renderer.GenerateButtonImage(width, height, label, ButtonDisabled)

	// Create the Button first, then pass its pointer as 'self'
	b := &Button{
		Label:    label,
		onClick:  nil, // Click handler set separately
		renderer: renderer,
	}
	b.interactiveComponent = NewInteractiveComponent(x, y, width, height, idle, pressed, hover, disabled, b)

	return b
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

// Draw draws the button's current state image to the screen using its absolute position.
// It now calls the embedded interactiveComponent's Draw method.
func (b *Button) Draw(screen *ebiten.Image) {
	b.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (b *Button) HandlePress() {
	b.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
func (b *Button) HandleRelease() {
	b.interactiveComponent.HandleRelease()
}

// HandleClick calls the specific button's onClick handler.
func (b *Button) HandleClick() {
	if b.state != ButtonDisabled && b.onClick != nil {
		b.onClick()
	}
}

// Focus is a no-op for a Button.
func (b *Button) Focus() {}

// Unfocus is a no-op for a Button.
func (b *Button) Unfocus() {}
