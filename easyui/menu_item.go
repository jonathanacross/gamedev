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
func (m *MenuItem) Update() {
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(m.Bounds, cx, cy)

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

// HandleClick calls the menu item's onClick handler.
func (m *MenuItem) HandleClick() {
	if m.onClick != nil {
		log.Printf("MenuItem '%s': Click handler triggered.", m.Label)
		m.onClick()
	}
}
