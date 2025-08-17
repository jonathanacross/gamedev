package main

import (
	"log"

	"image" // Required for image.Rectangle, image.Point

	"github.com/hajimehoshi/ebiten/v2"
)

// Container represents a non-interactive grouping component that can hold child UI elements.
// It serves as a visual and logical container for other components.
type Container struct {
	component                   // Embed the base component struct to manage bounds and children
	renderer      UiRenderer    // Reference to the UI renderer for drawing its background
	backgroundImg *ebiten.Image // The pre-rendered image for the container's background
}

// NewContainer creates a new Container instance.
// It initializes the container with its bounds and a reference to the UiRenderer.
func NewContainer(x, y, width, height int, renderer UiRenderer) *Container {
	c := &Container{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		renderer: renderer,
	}
	// Generate the initial background image for the container
	c.backgroundImg = c.renderer.GenerateContainerImage(width, height)
	return c
}

// Update iterates through and updates all child components within the container.
func (c *Container) Update() {
	for _, child := range c.children {
		child.Update()
	}
}

// Draw draws the container's background image first, and then draws all its child components.
func (c *Container) Draw(screen *ebiten.Image) {
	// Draw the container's background
	if c.backgroundImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(c.Bounds.Min.X), float64(c.Bounds.Min.Y))
		screen.DrawImage(c.backgroundImg, op)
	} else {
		log.Printf("Container at %v: WARNING: backgroundImg is nil! Cannot draw background.", c.Bounds)
	}

	// Draw all child components
	for _, child := range c.children {
		child.Draw(screen)
	}
}

// SetSize updates the size of the container and regenerates its background image.
func (c *Container) SetSize(width, height int) {
	c.Bounds.Max.X = c.Bounds.Min.X + width
	c.Bounds.Max.Y = c.Bounds.Min.Y + height
	c.backgroundImg = c.renderer.GenerateContainerImage(width, height)
}

// HandlePress is a no-op for a static container.
func (c *Container) HandlePress() {}

// HandleRelease is a no-op for a static container.
func (c *Container) HandleRelease() {}

// HandleClick is a no-op for a static container.
func (c *Container) HandleClick() {}
