package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MenuItem represents a single selectable item within a menu.
type MenuItem struct {
	component             // Embeds the base component struct for bounds and common properties
	Label     string      // The text displayed for this menu item
	onClick   func()      // Function to call when this menu item is selected
	state     ButtonState // To handle hover and pressed visual states (similar to a button)

	// Images for different states of the menu item
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

// Update handles the interaction logic for the menu item (hover, click).
func (m *MenuItem) Update() {
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := m.ContainsPoint(cx, cy)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorInBounds {
			m.state = ButtonPressed
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if m.state == ButtonPressed && cursorInBounds {
			if m.onClick != nil {
				m.onClick()
			}
			m.state = ButtonHover // Go to hover if released inside
		} else {
			m.state = ButtonIdle // Go to idle if released outside or not pressed
		}
	} else {
		if cursorInBounds {
			m.state = ButtonHover
		} else {
			m.state = ButtonIdle
		}
	}
}

// Draw draws the menu item, including its background color based on state and its label.
func (m *MenuItem) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.Bounds.Min.X), float64(m.Bounds.Min.Y))

	switch m.state {
	case ButtonIdle:
		screen.DrawImage(m.idleImage, op)
	case ButtonHover:
		screen.DrawImage(m.hoverImage, op)
	case ButtonPressed:
		screen.DrawImage(m.pressedImage, op)
	}
}
