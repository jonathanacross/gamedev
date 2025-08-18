package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// CheckboxState holds the different image sets for a checkbox's visual states.
// This is used internally by the Checkbox to manage its appearance.
type CheckboxStateImages struct {
	// Images when unchecked
	UncheckedIdle     *ebiten.Image
	UncheckedPressed  *ebiten.Image
	UncheckedHover    *ebiten.Image
	UncheckedDisabled *ebiten.Image

	// Images when checked
	CheckedIdle     *ebiten.Image
	CheckedPressed  *ebiten.Image
	CheckedHover    *ebiten.Image
	CheckedDisabled *ebiten.Image
}

// Checkbox represents a clickable UI component with a boolean (checked/unchecked) state.
type Checkbox struct {
	interactiveComponent
	Label          string
	Checked        bool
	OnCheckChanged func(bool)
	renderer       UiRenderer
	stateImages    CheckboxStateImages
}

// NewCheckbox creates a new Checkbox instance.
func NewCheckbox(x, y, width, height int, label string, initialChecked bool, renderer UiRenderer) *Checkbox {
	// Generate images for unchecked states
	uncheckedIdle := renderer.GenerateCheckboxImage(width, height, label, ButtonIdle, false)
	uncheckedPressed := renderer.GenerateCheckboxImage(width, height, label, ButtonPressed, false)
	uncheckedHover := renderer.GenerateCheckboxImage(width, height, label, ButtonHover, false)
	uncheckedDisabled := renderer.GenerateCheckboxImage(width, height, label, ButtonDisabled, false)

	// Generate images for checked states
	checkedIdle := renderer.GenerateCheckboxImage(width, height, label, ButtonIdle, true)
	checkedPressed := renderer.GenerateCheckboxImage(width, height, label, ButtonPressed, true)
	checkedHover := renderer.GenerateCheckboxImage(width, height, label, ButtonHover, true)
	checkedDisabled := renderer.GenerateCheckboxImage(width, height, label, ButtonDisabled, true)

	stateImages := CheckboxStateImages{
		UncheckedIdle:     uncheckedIdle,
		UncheckedPressed:  uncheckedPressed,
		UncheckedHover:    uncheckedHover,
		UncheckedDisabled: uncheckedDisabled,
		CheckedIdle:       checkedIdle,
		CheckedPressed:    checkedPressed,
		CheckedHover:      checkedHover,
		CheckedDisabled:   checkedDisabled,
	}

	// Initialize with the correct initial image based on `initialChecked`
	initialIdleImg := stateImages.UncheckedIdle
	initialPressedImg := stateImages.UncheckedPressed
	initialHoverImg := stateImages.UncheckedHover
	initialDisabledImg := stateImages.UncheckedDisabled

	if initialChecked {
		initialIdleImg = stateImages.CheckedIdle
		initialPressedImg = stateImages.CheckedPressed
		initialHoverImg = stateImages.CheckedHover
		initialDisabledImg = stateImages.CheckedDisabled
	}

	cb := &Checkbox{
		Label:       label,
		Checked:     initialChecked,
		renderer:    renderer, // Store the renderer
		stateImages: stateImages,
	}
	cb.interactiveComponent = NewInteractiveComponent(x, y, width, height,
		initialIdleImg, initialPressedImg, initialHoverImg, initialDisabledImg, cb) // Pass 'cb' as self
	return cb
}

// SetChecked updates the checkbox's checked state and regenerates its images.
func (c *Checkbox) SetChecked(checked bool) {
	if c.Checked == checked {
		return // No change needed
	}
	c.Checked = checked
	c.updateCurrentStateImages() // Update the images based on the new checked state
	if c.OnCheckChanged != nil {
		c.OnCheckChanged(c.Checked)
	}
}

// updateCurrentStateImages sets the correct image references in the embedded interactiveComponent
// based on the current `Checked` state.
func (c *Checkbox) updateCurrentStateImages() {
	// Re-generate images based on the new checked state.
	// The component's current interactive state (idle, hover, pressed) needs to be preserved
	currentInteractiveState := c.state
	if c.Checked {
		c.idleImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonIdle, true)
		c.pressedImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonPressed, true)
		c.hoverImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonHover, true)
		c.disabledImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonDisabled, true)
	} else {
		c.idleImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonIdle, false)
		c.pressedImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonPressed, false)
		c.hoverImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonHover, false)
		c.disabledImg = c.renderer.GenerateCheckboxImage(c.Bounds.Dx(), c.Bounds.Dy(), c.Label, ButtonDisabled, false)
	}
	c.state = currentInteractiveState // Restore interactive state
}

// Update calls the embedded interactiveComponent's Update method.
func (c *Checkbox) Update() {
	c.interactiveComponent.Update()
}

// Draw draws the checkbox component using the image for its current visual state.
// It now calls the embedded interactiveComponent's Draw method.
func (c *Checkbox) Draw(screen *ebiten.Image) {
	c.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (c *Checkbox) HandlePress() {
	c.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
func (c *Checkbox) HandleRelease() {
	c.interactiveComponent.HandleRelease()
}

// HandleClick toggles the checkbox's state and calls the OnCheckChanged handler.
func (c *Checkbox) HandleClick() {
	if c.state != ButtonDisabled {
		c.SetChecked(!c.Checked)
		log.Printf("Checkbox '%s' clicked. New state: %t", c.Label, c.Checked)
	}
}

// Focus is a no-op for a Checkbox.
func (c *Checkbox) Focus() {}

// Unfocus is a no-op for a Checkbox.
func (c *Checkbox) Unfocus() {}
