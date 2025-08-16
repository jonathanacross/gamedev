package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Component defines the basic interface for all UI widgets.
type Component interface {
	Draw(screen *ebiten.Image)
	Update()
	GetBounds() image.Rectangle
	HandlePress()   // Called when the mouse button is pressed down on this component.
	HandleRelease() // New: Called when the mouse button is released, if this component was pressed.
	HandleClick()   // Called if mouse was pressed AND released on this component.
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

// HandlePress is a dummy method for the base struct.
func (c *component) HandlePress() {}

// HandleRelease is a dummy method for the base struct.
func (c *component) HandleRelease() {}

// HandleClick is a dummy method for the base struct.
func (c *component) HandleClick() {}
