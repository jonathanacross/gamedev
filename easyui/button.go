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

	isPressedInside bool // New internal state: true if the mouse button was initially pressed while cursor was inside the button's bounds
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// Update handles the button's logic, including mouse interactions to change its state.
func (b *Button) Update() {
	// Get current cursor position
	cx, cy := ebiten.CursorPosition()

	// Use the ContainsPoint method from the embedded component to check if cursor is over the button
	cursorInBounds := b.ContainsPoint(cx, cy)

	// If the button is disabled, do not process any input and return immediately.
	if b.state == ButtonDisabled {
		return
	}

	// Step 1: Handle mouse button just pressed
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorInBounds {
			b.state = ButtonPressed
			b.isPressedInside = true // Mark that the press originated inside
		} else {
			b.isPressedInside = false // Mark that the press originated outside
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Step 2: Handle mouse button currently held down
		if b.isPressedInside { // Only if the press originated inside the button
			if cursorInBounds {
				b.state = ButtonPressed // Keep showing pressed if still over the button
			} else {
				b.state = ButtonIdle // Revert to idle if mouse moves off while held
			}
		}
		// If isPressedInside is false, it means the click started outside, so we don't change state based on drag.
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// Step 3: Handle mouse button just released
		if b.isPressedInside && cursorInBounds {
			// If released inside AND press originated inside, trigger click
			if b.onClick != nil {
				b.onClick()
			}
			b.state = ButtonIdle // Reset to idle after click
		} else {
			// If released outside, or released after an outside press, revert to idle/hover based on cursor position
			if cursorInBounds {
				b.state = ButtonHover
			} else {
				b.state = ButtonIdle
			}
		}
		b.isPressedInside = false // Always reset this flag on release
	} else {
		// Step 4: Handle hover state when no mouse button is pressed
		if cursorInBounds {
			b.state = ButtonHover // Show hover if cursor is over and not pressed
		} else {
			b.state = ButtonIdle // Otherwise, return to idle
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
