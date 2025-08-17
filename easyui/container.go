package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	// Required for image.Rectangle, image.Point
)

// Container represents a non-interactive grouping component that can hold child UI elements.
// It serves as a visual and logical container for other components.
// It now fully implements the Component interface.
type Container struct {
	component                   // Embed the base component struct to manage bounds and children
	renderer      UiRenderer    // Reference to the UI renderer for drawing its background
	backgroundImg *ebiten.Image // The pre-rendered image for the container's background
}

// NewContainer creates a new Container instance.
// It initializes the container with its bounds and a reference to the UiRenderer.
func NewContainer(x, y, width, height int, renderer UiRenderer) *Container {
	// Create the Container first, then pass its pointer as 'self'
	c := &Container{
		renderer: renderer,
	}
	c.component = NewComponent(x, y, width, height, c) // Pass 'c' as self

	// Generate the initial background image for the container
	c.backgroundImg = c.renderer.GenerateContainerImage(width, height)
	return c
}

// Update iterates through and updates all child components within the container.
// This method fully implements Component.Update().
func (c *Container) Update() {
	for _, child := range c.children {
		child.Update()
	}
}

// Draw draws the container's background image first, and then draws all its child components.
// All drawing is relative to the container's absolute position.
// This method fully implements Component.Draw().
func (c *Container) Draw(screen *ebiten.Image) {
	absX, absY := c.GetAbsolutePosition() // Get absolute position of the container itself

	// Draw the container's background
	if c.backgroundImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(absX), float64(absY)) // Draw background at its absolute position
		screen.DrawImage(c.backgroundImg, op)
	} else {
		log.Printf("Container at %v: WARNING: backgroundImg is nil! Cannot draw background.", c.Bounds)
	}

	// Draw all child components
	for _, child := range c.children {
		child.Draw(screen) // Child items will draw themselves using their own absolute positions
	}
}

// SetSize updates the size of the container and regenerates its background image.
func (c *Container) SetSize(width, height int) {
	c.Bounds.Max.X = c.Bounds.Min.X + width
	c.Bounds.Max.Y = c.Bounds.Min.Y + height
	c.backgroundImg = c.renderer.GenerateContainerImage(width, height)
}

// HandlePress is a no-op for a static container.
// This method fully implements Component.HandlePress().
func (c *Container) HandlePress() {}

// HandleRelease is a no-op for a static container.
// This method fully implements Component.HandleRelease().
func (c *Container) HandleRelease() {}

// HandleClick is a no-op for a static container.
// This method fully implements Component.HandleClick().
func (c *Container) HandleClick() {}
