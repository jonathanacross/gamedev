package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ui represents the root UI container, managing a collection of components.
type Ui struct {
	component // Embeds the base component struct
}

// NewUi creates a new Ui instance with the specified dimensions.
func NewUi(x, y, width, height int) *Ui {
	return &Ui{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
	}
}

// Update iterates through all child components and calls their Update methods.
func (u *Ui) Update() {
	for _, child := range u.children {
		child.Update()
	}
}

// Draw iterates through all child components and calls their Draw methods,
// passing the screen to draw on.
func (u *Ui) Draw(screen *ebiten.Image) {
	for _, child := range u.children {
		child.Draw(screen)
	}
}
