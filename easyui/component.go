package main

import (
	"image"
	// Added log for drawing warnings
	"github.com/hajimehoshi/ebiten/v2"
	// Needed for gg.NewContextForRGBA
	// Needed for drawing shapes
)

// ButtonState represents the current visual state of the button.
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
}

// component is the base struct that other UI widgets will embed.
type component struct {
	Bounds   image.Rectangle // The rectangular area occupied by the component
	children []Component     // Child components nested within this component
}

// GetBounds returns the rectangular bounds of the component.
func (c component) GetBounds() image.Rectangle {
	return c.Bounds
}

// AddChild adds a child component to this component's list of children.
func (c *component) AddChild(child Component) {
	c.children = append(c.children, child)
}

// GetChildren returns the child components.
func (c *component) GetChildren() []Component {
	return c.children
}

// ContainsPoint checks if a given (x, y) coordinate is within the component's bounds.
func ContainsPoint(rect image.Rectangle, x, y int) bool {
	return x >= rect.Min.X && x < rect.Max.X &&
		y >= rect.Min.Y && y < rect.Max.Y
}

// Dummy implementations for the base component to satisfy the interface.
func (c *component) HandlePress()   {}
func (c *component) HandleRelease() {}
func (c *component) HandleClick()   {}

// -----------------------------------------------------------------------------
// NEW: InteractiveComponent - encapsulates common state & input logic for buttons/dropdowns/menuitems
// -----------------------------------------------------------------------------

// interactiveComponent is a reusable struct for UI elements that have interactive states.
type interactiveComponent struct {
	component // Embeds the base component

	state ButtonState

	// References to the images for different states
	idleImg     *ebiten.Image
	pressedImg  *ebiten.Image
	hoverImg    *ebiten.Image
	disabledImg *ebiten.Image // Added for consistency, though ButtonDisabled is less dynamic
}

// NewInteractiveComponent creates a new interactiveComponent instance.
// It initializes the common component logic. The `drawFunc` is a placeholder
// for how the specific UI element generates its appearance for each state.
func NewInteractiveComponent(x, y, width, height int, idle, pressed, hover, disabled *ebiten.Image) interactiveComponent {
	return interactiveComponent{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		state:       ButtonIdle,
		idleImg:     idle,
		pressedImg:  pressed,
		hoverImg:    hover,
		disabledImg: disabled,
	}
}

// Update handles common hover state logic for interactive components.
// It ensures the pressed state only persists if the mouse is over the component.
func (ic *interactiveComponent) Update() {
	if ic.state == ButtonDisabled {
		return
	}

	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(ic.Bounds, cx, cy)

	if ic.state == ButtonPressed {
		// If currently pressed, check if mouse moved *off* the component.
		// If so, switch to ButtonIdle (visual feedback: no longer highlighted as pressed).
		if !cursorInBounds {
			ic.state = ButtonIdle // Reset to idle if mouse moves away while pressed
		}
		// If it's ButtonPressed and cursor is still in bounds, keep it ButtonPressed.
		return // Do not apply normal hover logic while actively pressed.
	}

	// Standard hover logic for Idle/Hover states
	if cursorInBounds {
		ic.state = ButtonHover
	} else {
		ic.state = ButtonIdle
	}
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
	if ContainsPoint(ic.Bounds, cx, cy) {
		ic.state = ButtonHover // Mouse released over component
	} else {
		ic.state = ButtonIdle // Mouse released away from component
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
		return ic.idleImg // Fallback
	}
}
