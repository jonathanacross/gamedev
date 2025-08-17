package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Ui represents the root UI container, managing a collection of components.
// It now fully implements the Component interface.
type Ui struct {
	component                  // Embed the base component struct (its bounds are window-relative)
	modalComponent   Component // The currently active modal component (e.g., a menu)
	pressedComponent Component // The component that was most recently pressed down
}

// NewUi creates a new Ui instance with the specified dimensions.
func NewUi(x, y, width, height int) *Ui {
	// Create the Ui first, then pass its pointer as 'self'
	u := &Ui{
		modalComponent:   nil,
		pressedComponent: nil,
	}
	u.component = NewComponent(x, y, width, height, u) // Pass 'u' as self; root UI has no parent, so it remains nil internally
	return u
}

// Update iterates through all child components and calls their Update methods.
// It also handles centralized mouse input detection for true click behavior and modal management.
// This method fully implements Component.Update().
func (u *Ui) Update() {
	// 1. Update all currently active components (modal first, then others).
	// This allows components to update their internal state (e.g., text field cursor blink, hover effect).
	if u.modalComponent != nil {
		u.modalComponent.Update()
	} else {
		for _, child := range u.children {
			child.Update()
		}
	}

	cx, cy := ebiten.CursorPosition() // Get current absolute cursor position

	// 2. Handle Mouse Button Press (ButtonDown)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Prioritize modal component for presses
		if u.modalComponent != nil {
			// Check if the press occurred on the modal component or its children
			if ContainsPoint(u.modalComponent, cx, cy) {
				u.pressedComponent = u.modalComponent
				u.modalComponent.HandlePress() // Notify modal it was pressed
				// Also check modal's children for internal press if modal allows it
				for _, child := range u.modalComponent.GetChildren() {
					if ContainsPoint(child, cx, cy) {
						u.pressedComponent = child // Update pressedComponent to the specific child
						child.HandlePress()
						break // Only one child can be pressed
					}
				}
			} else {
				// Clicked outside the modal, treat as a background click for the modal
				u.pressedComponent = nil // No specific component within the modal was pressed
			}
		} else {
			// No modal, check regular children in reverse order (top-most first)
			for i := len(u.children) - 1; i >= 0; i-- {
				child := u.children[i]
				if ContainsPoint(child, cx, cy) {
					u.pressedComponent = child
					u.pressedComponent.HandlePress() // Notify child it was pressed
					break                            // Only one component can be pressed
				}
			}
		}
	}

	// 3. Handle Mouse Button Release (ButtonUp)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// If a component was pressed down this sequence, check if it was released over it
		if u.pressedComponent != nil {
			if ContainsPoint(u.pressedComponent, cx, cy) {
				// This is a "click" - released over the same component that was pressed
				u.pressedComponent.HandleClick() // Trigger the click action
			}
			u.pressedComponent.HandleRelease() // Always call release to reset its visual state
			u.pressedComponent = nil           // Clear pressed component after handling release
		} else {
			// No component was pressed (e.g., clicked on background), but released.
			// This specifically handles closing a modal when clicking outside its bounds.
			if u.modalComponent != nil {
				// If a modal exists and no specific component was pressed initially, and the release
				// occurred outside the modal's bounds, then close the modal.
				if !ContainsPoint(u.modalComponent, cx, cy) {
					log.Println("Ui.Update: Mouse released on modal background, calling modal's HandleClick (to close).")
					u.modalComponent.HandleClick() // This will typically hide the menu
				}
			}
		}
	}
}

// Draw iterates through all child components and calls their Draw methods,
// passing the screen to draw on. It draws the modal component last to ensure it's on top.
// This method fully implements Component.Draw().
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

// HandlePress is a no-op for the root UI.
func (u *Ui) HandlePress() {}

// HandleRelease is a no-op for the root UI.
func (u *Ui) HandleRelease() {}

// HandleClick is a no-op for the root UI.
func (u *Ui) HandleClick() {}
