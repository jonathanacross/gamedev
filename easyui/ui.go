package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
		modalComponent: nil,
	}
}

func (u *Ui) Update() {
	// First, run the Update methods for all components.
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
		clickedHandled := false

		// Check the modal component first for a click on any of its children.
		if u.modalComponent != nil {
			// Iterate through the modal component's children to find the clicked one.
			// This is crucial for dropdown menus and other modal components with children.
			for i := len(u.modalComponent.GetChildren()) - 1; i >= 0; i-- {
				child := u.modalComponent.GetChildren()[i]
				if ContainsPoint(child.GetBounds(), cx, cy) {
					child.HandleClick()
					clickedHandled = true
					break
				}
			}

			// If the click wasn't on a child, check the modal component itself.
			if !clickedHandled && ContainsPoint(u.modalComponent.GetBounds(), cx, cy) {
				u.modalComponent.HandleClick()
				clickedHandled = true
			}
		}

		// If no modal component was clicked (or active), check the regular children.
		if !clickedHandled {
			for i := len(u.children) - 1; i >= 0; i-- {
				child := u.children[i]
				if ContainsPoint(child.GetBounds(), cx, cy) {
					child.HandleClick()
					break
				}
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
