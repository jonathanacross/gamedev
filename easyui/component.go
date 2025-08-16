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
	HandleClick()
	GetChildren() []Component // New method to get child components
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
// This is now a standalone function, not a method of a struct.
func ContainsPoint(rect image.Rectangle, x, y int) bool {
	return x >= rect.Min.X && x < rect.Max.X &&
		y >= rect.Min.Y && y < rect.Max.Y
}

// This is a dummy method to satisfy the Component interface for the base struct.
func (c *component) HandleClick() {}
