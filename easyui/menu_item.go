package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// MenuItem represents a single selectable item within a menu.
type MenuItem struct {
	component
	Label        string
	onClick      func()
	state        ButtonState
	idleImage    *ebiten.Image
	hoverImage   *ebiten.Image
	pressedImage *ebiten.Image
}

// NewMenuItem creates a new MenuItem instance.
func NewMenuItem(x, y, width, height int, label string, idle, hover, pressed *ebiten.Image) *MenuItem {
	return &MenuItem{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		Label:        label,
		state:        ButtonIdle,
		idleImage:    idle,
		hoverImage:   hover,
		pressedImage: pressed,
	}
}

// SetClickHandler sets the function to be executed when this menu item is clicked.
func (m *MenuItem) SetClickHandler(handler func()) {
	m.onClick = handler
}

// Update handles the interaction logic for the menu item (hover).
// It ensures the pressed state only persists if the mouse is over the menu item.
func (m *MenuItem) Update() {
	if m.state == ButtonDisabled {
		return
	}

	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(m.Bounds, cx, cy)

	if m.state == ButtonPressed {
		// If currently pressed, check if mouse moved *off* the menu item.
		if !cursorInBounds {
			m.state = ButtonIdle // Reset to idle if mouse moves away while pressed
		}
		// If it's ButtonPressed and cursor is still in bounds, keep it ButtonPressed.
		return // Do not apply normal hover logic while actively pressed.
	}

	// Standard hover logic for Idle/Hover states
	if cursorInBounds {
		m.state = ButtonHover
	} else {
		m.state = ButtonIdle
	}
}

// Draw draws the menu item, including its background color based on state and its label.
func (m *MenuItem) Draw(screen *ebiten.Image) {
	if m.idleImage == nil || m.hoverImage == nil || m.pressedImage == nil {
		log.Printf("MenuItem '%s': WARNING: One or more state images are nil! Idle: %t, Hover: %t, Pressed: %t", m.Label, m.idleImage != nil, m.hoverImage != nil, m.pressedImage != nil)
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.Bounds.Min.X), float64(m.Bounds.Min.Y))

	var imgToDraw *ebiten.Image

	switch m.state {
	case ButtonIdle:
		imgToDraw = m.idleImage
	case ButtonHover:
		imgToDraw = m.hoverImage
	case ButtonPressed:
		imgToDraw = m.pressedImage
	default:
		imgToDraw = m.idleImage
	}

	screen.DrawImage(imgToDraw, op)
}

// HandlePress is called when the mouse button is pressed down on this menu item.
func (m *MenuItem) HandlePress() {
	m.state = ButtonPressed
}

// HandleRelease is called when the mouse button is released, if this menu item was the pressed component.
// It ensures the menu item state transitions out of ButtonPressed to either Hover or Idle.
func (m *MenuItem) HandleRelease() {
	cx, cy := ebiten.CursorPosition()
	if ContainsPoint(m.Bounds, cx, cy) {
		m.state = ButtonHover // Mouse released over menu item
	} else {
		m.state = ButtonIdle // Mouse released away from menu item
	}
}

// HandleClick calls the menu item's onClick handler.
func (m *MenuItem) HandleClick() {
	if m.onClick != nil {
		log.Printf("MenuItem '%s': Click handler triggered.", m.Label)
		m.onClick()
	}
	// State transition (from ButtonPressed to Hover/Idle) is handled by HandleRelease.
}
