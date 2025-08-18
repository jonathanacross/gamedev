package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Ui represents the root UI container, managing a collection of components.
type Ui struct {
	component
	modalComponent   Component
	pressedComponent Component
	focusedComponent Component
}

// NewUi creates a new Ui instance with the specified dimensions.
func NewUi(x, y, width, height int) *Ui {
	u := &Ui{
		modalComponent:   nil,
		pressedComponent: nil,
		focusedComponent: nil,
	}
	u.component = NewComponent(x, y, width, height, u)
	return u
}

// Update iterates through all child components and calls their Update methods.
// It also handles centralized mouse input detection for true click behavior and modal management.
func (u *Ui) Update() {
	// Update all currently active components (modal first, then others).
	// This allows components to update their internal state (e.g., text field cursor blink, hover effect).
	if u.modalComponent != nil {
		u.modalComponent.Update()
	} else {
		for _, child := range u.children {
			child.Update()
		}
	}

	cx, cy := ebiten.CursorPosition()

	// Handle Mouse Button Press (ButtonDown)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		currentModal := u.modalComponent
		if currentModal != nil {
			// Check if the press occurred on the modal component or its children
			if ContainsPoint(currentModal, cx, cy) {
				u.pressedComponent = currentModal
				currentModal.HandlePress()

				// Check modal's children to find the most specific pressed component
				// Iterate in reverse to prioritize children drawn on top
				modalChildren := currentModal.GetChildren()
				for i := len(modalChildren) - 1; i >= 0; i-- {
					child := modalChildren[i]

					if ContainsPoint(child, cx, cy) {
						u.pressedComponent = child
						child.HandlePress()
						break // Only one child can be pressed
					}
				}
			} else {
				// Clicked outside the modal, treat as a background click for the modal
				u.pressedComponent = nil
			}
		} else {
			// No modal, check regular children in reverse order (top-most first)
			for i := len(u.children) - 1; i >= 0; i-- {
				child := u.children[i]
				if ContainsPoint(child, cx, cy) {
					u.pressedComponent = child
					u.pressedComponent.HandlePress()
					break
				}
			}
		}
	}

	// Handle Mouse Button Release (ButtonUp)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if u.pressedComponent != nil {
			// Store pressedComponent in a local variable BEFORE calling HandleClick
			// because HandleClick might clear u.pressedComponent (e.g., via ClearModal).
			originalPressedComponent := u.pressedComponent

			if ContainsPoint(originalPressedComponent, cx, cy) {
				// This is a "click" - released over the same component that was pressed
				originalPressedComponent.HandleClick()
			}

			// Always call release on the original component, regardless if HandleClick cleared u.pressedComponent.
			// The original component should still be valid.
			originalPressedComponent.HandleRelease()
			u.pressedComponent = nil // Clear pressed component after handling release
		} else {
			// No component was pressed (e.g., clicked on background), but released.
			// This specifically handles closing a modal when clicking outside its bounds.
			currentModal := u.modalComponent
			if currentModal != nil {
				if !ContainsPoint(currentModal, cx, cy) {
					log.Println("Ui.Update: Mouse released on modal background (outside modal's bounds), calling modal's HandleClick (to close).") // Diagnostic
					currentModal.HandleClick()                                                                                                     // This will typically hide the menu
				}
			}
		}
	}
}

// Draw iterates through all child components and calls their Draw methods.
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
	log.Printf("Ui.ClearModal: Modal component cleared (was type %T).", u.modalComponent)
	u.modalComponent = nil
	u.pressedComponent = nil // Also clear any lingering pressed state related to the modal
}

// SetFocusedComponent manages focus for interactive components.
// It will unfocus the previously focused component and focus the new one.
func (u *Ui) SetFocusedComponent(c Component) {
	if u.focusedComponent != nil {
		u.focusedComponent.Unfocus()
	}
	u.focusedComponent = c
	if u.focusedComponent != nil {
		u.focusedComponent.Focus()
	}
}

// Focus is a no-op for the root UI.
func (u *Ui) Focus() {}

// Unfocus is a no-op for the root UI.
func (u *Ui) Unfocus() {}

// HandlePress is a no-op for the root UI.
func (u *Ui) HandlePress() {}

// HandleRelease is a no-op for the root UI.
func (u *Ui) HandleRelease() {}

// HandleClick is a no-op for the root UI.
func (u *Ui) HandleClick() {}
