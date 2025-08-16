package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	isPressedInside bool
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// Update handles the button's logic, including mouse interactions to change its state.
func (b *Button) Update() {
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := b.ContainsPoint(cx, cy)

	if b.state == ButtonDisabled {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorInBounds {
			b.state = ButtonPressed
		}
		b.isPressedInside = cursorInBounds
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if b.isPressedInside {
			if cursorInBounds {
				b.state = ButtonPressed
			} else {
				b.state = ButtonIdle
			}
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if b.isPressedInside && cursorInBounds {
			b.state = ButtonIdle
			if b.onClick != nil {
				b.onClick()
			}
		} else if b.isPressedInside {
			b.state = ButtonIdle
		}
		b.isPressedInside = false
	} else {
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
