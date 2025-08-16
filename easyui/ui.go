package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ui represents the root UI container, managing a collection of components.
type Ui struct {
	component                // Embeds the base component struct
	modalComponent Component // A component that currently has modal focus (e.g., an open menu)
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
		modalComponent: nil, // Initially no modal component
	}
}

// Update iterates through all child components and calls their Update methods.
// It prioritizes updating the modal component if one exists, giving it exclusive input focus.
func (u *Ui) Update() {
	if u.modalComponent != nil {
		// If a modal component exists, only update it.
		// This ensures that input is directed exclusively to the modal UI.
		u.modalComponent.Update()
	} else {
		// If no modal component, update all regular child components.
		for _, child := range u.children {
			child.Update()
		}
	}
}

// Draw iterates through all child components and calls their Draw methods,
// passing the screen to draw on. It draws the modal component last to ensure it's on top.
func (u *Ui) Draw(screen *ebiten.Image) {
	// Draw all regular child components first.
	for _, child := range u.children {
		child.Draw(screen)
	}
	// If a modal component exists, draw it last so it appears on top of other UI elements.
	if u.modalComponent != nil {
		u.modalComponent.Draw(screen)
	}
}

// SetModal sets a component as the current modal, giving it exclusive input focus and drawing priority.
func (u *Ui) SetModal(c Component) {
	u.modalComponent = c
}

// ClearModal removes the current modal component, returning input focus to the regular UI.
func (u *Ui) ClearModal() {
	u.modalComponent = nil
}
