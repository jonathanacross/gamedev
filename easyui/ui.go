package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Ui represents the root UI container, managing a collection of components.
type Ui struct {
	component
	modalComponent Component
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
		modalComponent: nil,
	}
}

// Update iterates through all child components and calls their Update methods.
// It also handles centralized click detection.
func (u *Ui) Update() {
	// First, update the relevant components.
	if u.modalComponent != nil {
		u.modalComponent.Update()
	} else {
		for _, child := range u.children {
			child.Update()
		}
	}

	// Centralized click handling
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()

		// If a modal component exists, handle clicks exclusively for it and its children.
		if u.modalComponent != nil {
			var handledClick bool

			// Check for clicks on the modal's children in reverse order.
			for i := len(u.modalComponent.GetChildren()) - 1; i >= 0; i-- {
				child := u.modalComponent.GetChildren()[i]
				if ContainsPoint(child.GetBounds(), cx, cy) {
					child.HandleClick()
					handledClick = true
					break
				}
			}

			// If no child was clicked, check if the click was on the modal itself.
			if !handledClick && ContainsPoint(u.modalComponent.GetBounds(), cx, cy) {
				u.modalComponent.HandleClick()
				handledClick = true
			}

			// If the click was not on the modal or its children, clear the modal.
			// The key here is to `return` immediately after clearing the modal.
			if !handledClick {
				u.ClearModal()
				return
			}
		}

		// If no modal component is active, handle clicks on regular children.
		for i := len(u.children) - 1; i >= 0; i-- {
			child := u.children[i]
			if ContainsPoint(child.GetBounds(), cx, cy) {
				child.HandleClick()
				return // A click was handled. Exit the update loop.
			}
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
	log.Printf("Ui.SetModal: Modal component set to type %T.", c)
}

// ClearModal removes the current modal component, returning input focus to the regular UI.
func (u *Ui) ClearModal() {
	log.Printf("Ui.ClearModal: Modal component cleared.")
	u.modalComponent = nil
}
