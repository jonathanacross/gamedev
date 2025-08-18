package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// RadioButton represents a clickable UI component with a boolean (checked/unchecked) state,
// typically used in groups where only one can be selected at a time.
type RadioButton struct {
	interactiveComponent
	Label          string
	Checked        bool
	OnCheckChanged func(bool)
	renderer       UiRenderer
}

// NewRadioButton creates a new RadioButton instance.
func NewRadioButton(x, y, width, height int, label string, initialChecked bool, renderer UiRenderer) *RadioButton {
	// Generate initial images based on the checked state
	initialIdleImg := renderer.GenerateRadioButtonImage(width, height, label, ButtonIdle, initialChecked)
	initialPressedImg := renderer.GenerateRadioButtonImage(width, height, label, ButtonPressed, initialChecked)
	initialHoverImg := renderer.GenerateRadioButtonImage(width, height, label, ButtonHover, initialChecked)
	initialDisabledImg := renderer.GenerateRadioButtonImage(width, height, label, ButtonDisabled, initialChecked)

	rb := &RadioButton{
		Label:    label,
		Checked:  initialChecked,
		renderer: renderer,
	}
	rb.interactiveComponent = NewInteractiveComponent(x, y, width, height,
		initialIdleImg, initialPressedImg, initialHoverImg, initialDisabledImg, rb)

	return rb
}

// SetChecked updates the radio button's checked state and regenerates its images.
// This method should primarily be called by a RadioButtonGroup to ensure mutual exclusivity.
func (rb *RadioButton) SetChecked(checked bool) {
	if rb.Checked == checked {
		return // No change needed
	}
	rb.Checked = checked
	rb.updateCurrentStateImages() // Update the images based on the new checked state
	if rb.OnCheckChanged != nil {
		rb.OnCheckChanged(rb.Checked)
	}
}

func (rb *RadioButton) IsChecked() bool {
	return rb.Checked
}

// updateCurrentStateImages regenerates the images for the current checked state.
func (rb *RadioButton) updateCurrentStateImages() {
	// Remember the current interactive state before regenerating images
	currentInteractiveState := rb.state

	rb.idleImg = rb.renderer.GenerateRadioButtonImage(rb.Bounds.Dx(), rb.Bounds.Dy(), rb.Label, ButtonIdle, rb.Checked)
	rb.pressedImg = rb.renderer.GenerateRadioButtonImage(rb.Bounds.Dx(), rb.Bounds.Dy(), rb.Label, ButtonPressed, rb.Checked)
	rb.hoverImg = rb.renderer.GenerateRadioButtonImage(rb.Bounds.Dx(), rb.Bounds.Dy(), rb.Label, ButtonHover, rb.Checked)
	rb.disabledImg = rb.renderer.GenerateRadioButtonImage(rb.Bounds.Dx(), rb.Bounds.Dy(), rb.Label, ButtonDisabled, rb.Checked)

	rb.state = currentInteractiveState // Restore interactive state
}

// Update calls the embedded interactiveComponent's Update method.
func (rb *RadioButton) Update() {
	rb.interactiveComponent.Update()
}

// Draw draws the radio button component using the image for its current visual state.
func (rb *RadioButton) Draw(screen *ebiten.Image) {
	rb.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (rb *RadioButton) HandlePress() {
	rb.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
func (rb *RadioButton) HandleRelease() {
	rb.interactiveComponent.HandleRelease()
}

// HandleClick toggles the radio button's state and calls the OnCheckChanged handler.
func (rb *RadioButton) HandleClick() {
	if rb.state != ButtonDisabled && !rb.Checked { // Only allow checking if not disabled and not already checked
		rb.SetChecked(true)
		log.Printf("RadioButton '%s' checked.", rb.Label)

		// Delegate the click to the parent ButtonGroup to handle mutual exclusivity.
		if parent := rb.GetParent(); parent != nil {
			if bg, ok := parent.(*ButtonGroup); ok {
				if bg.SelectionMode == SingleSelection {
					bg.HandleChildClick(rb.self)
				}
			}
		}
	}
}
