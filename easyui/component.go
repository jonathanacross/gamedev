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

// ContainsPoint checks if a given (x, y) coordinate is within the component's bounds.
func (c *component) ContainsPoint(x, y int) bool {
	return x >= c.Bounds.Min.X && x < c.Bounds.Max.X &&
		y >= c.Bounds.Min.Y && y < c.Bounds.Max.Y
}
