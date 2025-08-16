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
	modalComponent   Component // The currently active modal component (e.g., a menu)
	pressedComponent Component // The component that was most recently pressed down
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
		modalComponent:   nil,
		pressedComponent: nil,
	}
}

// Update iterates through all child components and calls their Update methods.
// It also handles centralized mouse input detection for true click behavior and modal management.
func (u *Ui) Update() {
	// 1. Update all currently active components (modal first, then others).
	if u.modalComponent != nil {
		u.modalComponent.Update()
	} else {
		for _, child := range u.children {
			child.Update()
		}
	}

	cx, cy := ebiten.CursorPosition() // Get current cursor position

	// 2. Handle Mouse Button Press (ButtonDown)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// If a modal is active, only consider components within the modal for pressing.
		if u.modalComponent != nil {
			// Check modal's children first (iterating in reverse for Z-order, if applicable)
			for i := len(u.modalComponent.GetChildren()) - 1; i >= 0; i-- {
				child := u.modalComponent.GetChildren()[i]
				if ContainsPoint(child.GetBounds(), cx, cy) {
					u.pressedComponent = child
					u.pressedComponent.HandlePress()
					return // A component was pressed, so we're done with press handling for this frame.
				}
			}
			// If no child was pressed, check if the modal itself was pressed
			if ContainsPoint(u.modalComponent.GetBounds(), cx, cy) {
				u.pressedComponent = u.modalComponent
				u.pressedComponent.HandlePress()
				return // Modal was pressed.
			}
		} else {
			// No modal active, check regular children for press (in reverse Z-order)
			for i := len(u.children) - 1; i >= 0; i-- {
				child := u.children[i]
				if ContainsPoint(child.GetBounds(), cx, cy) {
					u.pressedComponent = child
					u.pressedComponent.HandlePress()
					return // A component was pressed.
				}
			}
		}
		// If no component was pressed, u.pressedComponent remains nil.
	}

	// 3. Handle Mouse Button Release (ButtonUp)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// Capture modal state *before* any component's HandleClick might change it.
		modalBeforeClickProcessing := u.modalComponent

		// Phase 3.1: Always process the release for the component that was initially pressed.
		if u.pressedComponent != nil {
			u.pressedComponent.HandleRelease() // Allow the component to reset its visual state

			// If the release occurred over the same component, it's a "true click".
			if ContainsPoint(u.pressedComponent.GetBounds(), cx, cy) {
				u.pressedComponent.HandleClick() // This might call u.SetModal()
			}
			// Clear pressedComponent after its release/click has been processed.
			u.pressedComponent = nil
		}

		// Phase 3.2: Handle modal state changes and outside clicks.
		// This logic runs *after* any potential component click handler has executed.

		// If a modal was opened by the click event in THIS frame
		// (i.e., there was NO modal before processing, but there IS a modal now).
		if modalBeforeClickProcessing == nil && u.modalComponent != nil {
			// A modal was just opened by a component's HandleClick.
			// Crucially, we MUST return immediately to prevent it from being dismissed
			// by the "outside click" logic in this *same* frame.
			return
		}

		// If a modal is currently active (and it wasn't just opened by this click event),
		// check if the click occurred outside it to dismiss it.
		if u.modalComponent != nil {
			clickWasInsideModalContent := false
			// Check if click was within the modal's primary bounds
			if ContainsPoint(u.modalComponent.GetBounds(), cx, cy) {
				clickWasInsideModalContent = true // Assume inside initially

				// Now, check if it was on any of the modal's children
				clickedOnChild := false
				for _, child := range u.modalComponent.GetChildren() {
					if ContainsPoint(child.GetBounds(), cx, cy) {
						clickedOnChild = true
						break
					}
				}

				// If it was inside the modal's bounds but NOT on a child, it's on the modal's background.
				if !clickedOnChild {
					log.Println("Ui.Update: Mouse released on modal background, calling modal's HandleClick (to close).")
					u.modalComponent.HandleClick() // This will typically hide the menu
					return                         // Modal background click handled, done for this frame.
				}
			}

			// If the click was entirely outside the modal's bounds (not even on its background), clear the modal.
			if !clickWasInsideModalContent {
				log.Println("Ui.Update: Mouse released entirely outside modal's bounds, clearing modal.")
				u.ClearModal()
				return // Modal cleared by outside click, done for this frame.
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
