package main

import (
	// Import for image.Rectangle
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil" // Import for input utilities
)

// ButtonState represents the current visual state of the button.
type ButtonState int

const (
	ButtonIdle     ButtonState = iota // Default state
	ButtonPressed                     // When the button is being clicked
	ButtonHover                       // When the mouse cursor is over the button
	ButtonDisabled                    // When the button is not interactive
)

// Button represents a clickable UI button.
type Button struct {
	component // Embeds the base component struct

	idle     *ebiten.Image // Image for the idle state
	pressed  *ebiten.Image // Image for the pressed state
	hover    *ebiten.Image // Image for the hover state
	disabled *ebiten.Image // Image for the disabled state

	state   ButtonState // Current state of the button
	onClick func()      // Function to call when the button is clicked
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// Update handles the button's logic, including mouse interactions to change its state.
func (b *Button) Update() {
	// Get current cursor position
	cx, cy := ebiten.CursorPosition()

	// Use the new ContainsPoint method from the embedded component
	cursorInBounds := b.ContainsPoint(cx, cy)

	// Check for disabled state first
	if b.state == ButtonDisabled {
		return // Do not process input if disabled
	}

	// Handle button press
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorInBounds {
			b.state = ButtonPressed
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if b.state == ButtonPressed && cursorInBounds {
			b.state = ButtonIdle // Return to idle after release
			if b.onClick != nil {
				b.onClick() // Trigger click handler
			}
		} else if b.state == ButtonPressed {
			b.state = ButtonIdle // If released outside, go back to idle
		}
	}

	// Handle hover state
	if b.state != ButtonPressed { // Do not change to hover if currently pressed
		if cursorInBounds {
			b.state = ButtonHover
		} else {
			b.state = ButtonIdle
		}
	}
}

// Draw draws the button's current state image to the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Translate the image to the button's position
	op.GeoM.Translate(float64(b.Bounds.Min.X), float64(b.Bounds.Min.Y))

	// Draw the image corresponding to the current state
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
