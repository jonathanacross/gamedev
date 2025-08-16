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

// Update handles the button's logic, including mouse interactions to change its state.
func (b *Button) Update() {
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(b.Bounds, cx, cy)

	if b.state == ButtonDisabled {
		return
	}

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

// HandleClick calls the button's onClick handler.
func (b *Button) HandleClick() {
	if b.onClick != nil {
		b.onClick()
	}
}
