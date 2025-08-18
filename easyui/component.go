package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// ButtonState represents the current visual state of a button-like component.
type ButtonState int

const (
	ButtonIdle ButtonState = iota
	ButtonPressed
	ButtonHover
	ButtonDisabled
)

// Component defines the basic interface for all UI widgets.
type Component interface {
	Draw(screen *ebiten.Image)
	Update()
	GetBounds() image.Rectangle
	HandlePress()
	HandleRelease()
	HandleClick()
	GetChildren() []Component
	SetParent(Component)
	GetAbsolutePosition() (int, int)
}

// component is the base struct that other UI widgets will embed.
// It includes a reference to its parent for hierarchical positioning, and a self-reference
// to the concrete Component that embeds it, for correct parent setting.
type component struct {
	Bounds   image.Rectangle // Coordinates are relative to the parent.
	children []Component
	parent   Component
	self     Component
}

// NewComponent creates a new base component instance.
// The 'self' parameter should be the concrete Component that embeds this base component.
func NewComponent(x, y, width, height int, self Component) component {
	return component{
		Bounds: image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + width, Y: y + height},
		},
		self: self,
	}
}

// GetBounds returns the rectangular bounds of the component (relative to its parent).
func (c component) GetBounds() image.Rectangle {
	return c.Bounds
}

// AddChild adds a child component to this component's list of children
// and sets the child's parent reference to the embedding Component (c.self).
func (c *component) AddChild(child Component) {
	c.children = append(c.children, child)
	child.SetParent(c.self)
}

// GetChildren returns the child components.
func (c *component) GetChildren() []Component {
	return c.children
}

// SetParent sets the parent of this component.
func (c *component) SetParent(parent Component) {
	c.parent = parent
}

// GetAbsolutePosition calculates and returns the component's absolute (window) X and Y coordinates.
func (c *component) GetAbsolutePosition() (int, int) {
	absX, absY := c.Bounds.Min.X, c.Bounds.Min.Y

	if c.parent != nil {
		parentAbsX, parentAbsY := c.parent.GetAbsolutePosition()
		absX += parentAbsX
		absY += parentAbsY
	}
	return absX, absY
}

// ContainsPoint checks if a given (x, y) coordinate (absolute window coordinates)
// is within the component's absolute bounds.
// This function expects a Component interface to correctly call GetAbsolutePosition and GetBounds.
func ContainsPoint(comp Component, absX, absY int) bool {
	compAbsX, compAbsY := comp.GetAbsolutePosition()

	bounds := comp.GetBounds()
	compAbsBounds := image.Rectangle{
		Min: image.Point{X: compAbsX, Y: compAbsY},
		Max: image.Point{X: compAbsX + bounds.Dx(), Y: compAbsY + bounds.Dy()},
	}
	return absX >= compAbsBounds.Min.X && absX < compAbsBounds.Max.X &&
		absY >= compAbsBounds.Min.Y && absY < compAbsBounds.Max.Y
}

// interactiveComponent is a base struct for components that respond to mouse interaction.
// It manages common visual states like idle, pressed, hover, and disabled.
type interactiveComponent struct {
	component
	state       ButtonState
	idleImg     *ebiten.Image
	pressedImg  *ebiten.Image
	hoverImg    *ebiten.Image
	disabledImg *ebiten.Image
}

// NewInteractiveComponent creates a new interactiveComponent.
// It requires the 'self' parameter which is the concrete Component embedding it.
func NewInteractiveComponent(x, y, width, height int, idle, pressed, hover, disabled *ebiten.Image, self Component) interactiveComponent {
	return interactiveComponent{
		component:   NewComponent(x, y, width, height, self),
		state:       ButtonIdle,
		idleImg:     idle,
		pressedImg:  pressed,
		hoverImg:    hover,
		disabledImg: disabled,
	}
}

// Update handles the interaction logic for the component (primarily hover state).
func (ic *interactiveComponent) Update() {
	if ic.state == ButtonDisabled {
		return
	}

	cx, cy := ebiten.CursorPosition() // in absolute coordinates
	cursorInBounds := ContainsPoint(ic.self, cx, cy)

	// Only manage hover state here. Do not change state if currently pressed down.
	if ic.state != ButtonPressed {
		if cursorInBounds {
			ic.state = ButtonHover
		} else {
			ic.state = ButtonIdle
		}
	}
}

// Draw handles the generic drawing for any interactiveComponent.
// Concrete components will often call this method.
func (ic *interactiveComponent) Draw(screen *ebiten.Image) {
	absX, absY := ic.GetAbsolutePosition() // Get absolute position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(ic.GetCurrentStateImage(), op)
}

// HandlePress sets the component to the pressed state.
func (ic *interactiveComponent) HandlePress() {
	if ic.state != ButtonDisabled {
		ic.state = ButtonPressed
	}
}

// HandleRelease resets the component's state after a mouse release.
func (ic *interactiveComponent) HandleRelease() {
	if ic.state == ButtonDisabled {
		return
	}
	cx, cy := ebiten.CursorPosition()
	// If the mouse is released over the component, set to Hover, otherwise Idle.
	if ContainsPoint(ic.self, cx, cy) {
		ic.state = ButtonHover
	} else {
		ic.state = ButtonIdle
	}
}

// GetCurrentStateImage returns the correct image for the component's current state.
func (ic *interactiveComponent) GetCurrentStateImage() *ebiten.Image {
	switch ic.state {
	case ButtonIdle:
		return ic.idleImg
	case ButtonPressed:
		return ic.pressedImg
	case ButtonHover:
		return ic.hoverImg
	case ButtonDisabled:
		return ic.disabledImg
	default:
		// Should not happen, but return idle image as a fallback
		log.Printf("interactiveComponent: Unknown ButtonState %d, returning idle image.", ic.state)
		return ic.idleImg
	}
}

// HandleClick is a dummy method for the base interactiveComponent.
// Concrete interactive components (Button, TextField, Checkbox, Dropdown) will override this
// to perform their specific action.
func (ic *interactiveComponent) HandleClick() {}
