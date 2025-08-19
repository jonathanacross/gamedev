package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Button represents a clickable UI button.
type Button struct {
	interactiveComponent
	Label        string
	IsToggleable bool
	Checked      bool
	onClick      func()
	renderer     UiRenderer
}

// NewButton creates a new Button instance with the specified dimensions, label, toggle state, and renderer.
func NewButton(x, y, width, height int, label string, isToggleable bool, renderer UiRenderer) *Button {
	// Generate initial images based on the checked state
	initialIdleImg := renderer.GenerateButtonImage(width, height, label, ButtonIdle, false)
	initialPressedImg := renderer.GenerateButtonImage(width, height, label, ButtonPressed, false)
	initialHoverImg := renderer.GenerateButtonImage(width, height, label, ButtonHover, false)
	initialDisabledImg := renderer.GenerateButtonImage(width, height, label, ButtonDisabled, false)

	// Create the Button first, then pass its pointer as 'self'
	b := &Button{
		Label:        label,
		IsToggleable: isToggleable,
		Checked:      false,
		renderer:     renderer,
	}
	b.interactiveComponent = NewInteractiveComponent(x, y, width, height, initialIdleImg, initialPressedImg, initialHoverImg, initialDisabledImg, b)

	return b
}

// SetClickHandler sets the function to be executed when the button is clicked.
func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

// SetText updates the button's text and regenerates its state images.
func (b *Button) SetText(newText string) {
	b.Label = newText
	b.updateCurrentStateImages()
}

// SetChecked updates the button's checked state and regenerates its images.
// This method should primarily be called by a ButtonGroup to ensure mutual exclusivity.
func (b *Button) SetChecked(checked bool) {
	if b.Checked == checked {
		return // No change needed
	}
	b.Checked = checked
	b.updateCurrentStateImages()
}

// IsChecked returns the current checked state of the button.
func (b *Button) IsChecked() bool {
	return b.Checked
}

// updateCurrentStateImages regenerates the images for the current checked state.
func (b *Button) updateCurrentStateImages() {
	// Remember the current interactive state before regenerating images
	currentInteractiveState := b.state

	b.idleImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonIdle, b.Checked)
	b.pressedImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonPressed, b.Checked)
	b.hoverImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonHover, b.Checked)
	b.disabledImg = b.renderer.GenerateButtonImage(b.Bounds.Dx(), b.Bounds.Dy(), b.Label, ButtonDisabled, b.Checked)

	b.state = currentInteractiveState // Restore interactive state
}

// Update calls the embedded interactiveComponent's Update method.
func (b *Button) Update() {
	b.interactiveComponent.Update()
}

// Draw draws the button's current state image to the screen using its absolute position.
func (b *Button) Draw(screen *ebiten.Image) {
	b.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (b *Button) HandlePress() {
	b.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
func (b *Button) HandleRelease() {
	b.interactiveComponent.HandleRelease()
}

// HandleClick handles the button's click behavior, including toggling if IsToggleable is true.
func (b *Button) HandleClick() {
	log.Printf("Button '%s' clicked.", b.Label)

	// Invert the if statement to reduce indentation.
	if b.state == ButtonDisabled {
		return
	}

	// Only toggle the state if the button is configured as toggleable.
	if b.IsToggleable {
		b.SetChecked(!b.Checked)

		log.Printf("Button '%s' toggled. New state: %t\n", b.Label, b.Checked)
		// Delegate the click to the parent ButtonGroup to handle mutual exclusivity.
		if parent := b.GetParent(); parent != nil {
			log.Printf("  Delegating to parent\n")
			if bg, ok := parent.(*ButtonGroup); ok && bg.SelectionMode == SingleSelection {
				bg.HandleChildClick(b.self)
			}
		}
	}

	if b.onClick != nil {
		b.onClick()
	}
}
