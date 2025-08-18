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
}

// NewUi creates a new Ui instance with the specified dimensions.
func NewUi(x, y, width, height int) *Ui {
	u := &Ui{
		modalComponent:   nil,
		pressedComponent: nil,
	}
	u.component = NewComponent(x, y, width, height, u)
	return u
}

// Update iterates through all child components and calls their Update methods.
// It also handles centralized mouse input detection for true click behavior and modal management.
func (u *Ui) Update() {
	// Update all currently active components (modal first, then others).
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
				// Check for disabled state before handling press
				if stateful, ok := currentModal.(ComponentWithState); ok && stateful.GetState() == ButtonDisabled {
					u.pressedComponent = nil
					return
				}

				u.pressedComponent = currentModal
				currentModal.HandlePress()

				// Check modal's children to find the most specific pressed component
				// Iterate in reverse to prioritize children drawn on top
				modalChildren := currentModal.GetChildren()
				for i := len(modalChildren) - 1; i >= 0; i-- {
					child := modalChildren[i]
					if ContainsPoint(child, cx, cy) {
						if stateful, ok := child.(ComponentWithState); ok && stateful.GetState() == ButtonDisabled {
							u.pressedComponent = nil
							return
						}

						u.pressedComponent = child
						u.pressedComponent.HandlePress()
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
					if stateful, ok := child.(ComponentWithState); ok && stateful.GetState() == ButtonDisabled {
						u.pressedComponent = nil
						return
					}

					u.pressedComponent = child
					u.pressedComponent.HandlePress()
					break
				}
			}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Mouse is currently held down. Check if the pressed component's state needs to be updated.
		if u.pressedComponent != nil {
			// Use a type assertion with a comma-ok check to handle components that do not have a state (like Label).
			if statefulComponent, ok := u.pressedComponent.(ComponentWithState); ok {
				if statefulComponent.GetState() == ButtonDisabled {
					return
				}

				if ContainsPoint(statefulComponent, cx, cy) {
					statefulComponent.SetState(ButtonPressed)
				} else {
					statefulComponent.SetState(ButtonIdle)
				}
			}
		}
	}

	// Handle Mouse Button Release (ButtonUp)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if u.pressedComponent != nil {
			originalPressedComponent := u.pressedComponent

			if stateful, ok := originalPressedComponent.(ComponentWithState); ok && stateful.GetState() == ButtonDisabled {
				u.pressedComponent = nil
				return
			}

			if ContainsPoint(originalPressedComponent, cx, cy) {
				originalPressedComponent.HandleClick()
			}
			originalPressedComponent.HandleRelease()
			u.pressedComponent = nil
		} else {
			currentModal := u.modalComponent
			if currentModal != nil {
				if !ContainsPoint(currentModal, cx, cy) {
					log.Println("Ui.Update: Mouse released on modal background (outside modal's bounds), calling modal's HandleClick (to close).")
					currentModal.HandleClick()
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
	u.pressedComponent = nil
}

// HandlePress is a no-op for the root UI.
func (u *Ui) HandlePress() {}

// HandleRelease is a no-op for the root UI.
func (u *Ui) HandleRelease() {}

// HandleClick is a no-op for the root UI.
func (u *Ui) HandleClick() {}
