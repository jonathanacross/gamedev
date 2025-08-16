package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ButtonState represents the current visual state of the button.
type ButtonState int

const (
	ButtonIdle ButtonState = iota
	ButtonPressed
	ButtonHover
	ButtonDisabled
)

// Button represents a clickable UI button.
type Button struct {
	component

	idle     *ebiten.Image
	pressed  *ebiten.Image
	hover    *ebiten.Image
	disabled *ebiten.Image

	state   ButtonState
	onClick func()
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// Update handles the button's logic, primarily hover state.
// It ensures the pressed state only persists if the mouse is over the button.
func (b *Button) Update() {
	if b.state == ButtonDisabled {
		return
	}

	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(b.Bounds, cx, cy)

	if b.state == ButtonPressed {
		// If currently pressed, check if mouse moved *off* the button.
		// If so, switch to ButtonIdle (visual feedback: no longer highlighted as pressed).
		if !cursorInBounds {
			b.state = ButtonIdle // Reset to idle if mouse moves away while pressed
		}
		// If it's ButtonPressed and cursor is still in bounds, keep it ButtonPressed.
		return // Do not apply normal hover logic while actively pressed.
	}

	// Standard hover logic for Idle/Hover states
	if cursorInBounds {
		b.state = ButtonHover
	} else {
		b.state = ButtonIdle
	}
}

// Draw draws the button's current state image to the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Bounds.Min.X), float64(b.Bounds.Min.Y))

	switch b.state {
	case ButtonIdle:
		screen.DrawImage(b.idle, op)
	case ButtonPressed:
		screen.DrawImage(b.pressed, op)
	case ButtonHover:
		screen.DrawImage(b.hover, op)
	case ButtonDisabled:
		screen.DrawImage(b.disabled, op)
	}
}

// HandlePress is called when the mouse button is pressed down on this button.
func (b *Button) HandlePress() {
	if b.state != ButtonDisabled {
		b.state = ButtonPressed // Set to pressed state
	}
}

// HandleRelease is called when the mouse button is released, if this button was the pressed component.
// It ensures the button state transitions out of ButtonPressed to either Hover or Idle.
func (b *Button) HandleRelease() {
	if b.state == ButtonDisabled {
		return // Do nothing if disabled
	}
	cx, cy := ebiten.CursorPosition()
	if ContainsPoint(b.Bounds, cx, cy) {
		b.state = ButtonHover // Mouse released over button
	} else {
		b.state = ButtonIdle // Mouse released away from button
	}
}

// HandleClick is called when the mouse button is released on this button,
// and it was also pressed on this button (true click).
func (b *Button) HandleClick() {
	if b.state != ButtonDisabled && b.onClick != nil {
		b.onClick()
	}
	// State transition (from ButtonPressed to Hover/Idle) is handled by HandleRelease.
}
