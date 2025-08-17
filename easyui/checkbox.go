package main

import (
	"log"
	// For cursor blinking
	"github.com/hajimehoshi/ebiten/v2"
	// Required for color.RGBA
	// Required for image.Rectangle, image.Point
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
	interactiveComponent                       // Embed the reusable interactive component logic
	Label                string                // Text label next to the checkbox
	Checked              bool                  // Current checked state
	OnCheckChanged       func(bool)            // Handler for when the checked state changes
	uiGenerator          *BareBonesUiGenerator // Reference to the generator for image regeneration
	stateImages          CheckboxStateImages   // All state-specific images for checked/unchecked
}

// NewCheckbox creates a new Checkbox instance.
// It is now a standalone function.
func NewCheckbox(x, y, width, height int, label string, initialChecked bool, uiGen *BareBonesUiGenerator) *Checkbox {
	// Colors that stay constant for the overall component background and label text
	componentBgColor := uiGen.theme.BackgroundColor
	labelColor := uiGen.theme.OnPrimaryColor
	checkmarkColor := uiGen.theme.OnPrimaryColor // Checkmark color remains constant

	// Generate images for unchecked states
	uncheckedIdle := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.PrimaryColor, checkmarkColor, labelColor, false, label)
	uncheckedPressed := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.AccentColor, checkmarkColor, labelColor, false, label) // Box outline changes to AccentColor
	uncheckedHover := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.AccentColor, checkmarkColor, labelColor, false, label)   // Box outline changes to AccentColor
	uncheckedDisabled := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.PrimaryColor, checkmarkColor, labelColor, false, label)

	// Generate images for checked states
	checkedIdle := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.PrimaryColor, checkmarkColor, labelColor, true, label)
	checkedPressed := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.AccentColor, checkmarkColor, labelColor, true, label) // Box outline changes to AccentColor
	checkedHover := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.AccentColor, checkmarkColor, labelColor, true, label)   // Box outline changes to AccentColor
	checkedDisabled := uiGen.generateCheckboxImage(width, height, componentBgColor, uiGen.theme.PrimaryColor, checkmarkColor, labelColor, true, label)

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
		interactiveComponent: NewInteractiveComponent(x, y, width, height,
			initialIdleImg, initialPressedImg, initialHoverImg, initialDisabledImg),
		Label:       label,
		Checked:     initialChecked,
		uiGenerator: uiGen,
		stateImages: stateImages,
	}
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
	if c.Checked {
		c.idleImg = c.stateImages.CheckedIdle
		c.pressedImg = c.stateImages.CheckedPressed
		c.hoverImg = c.stateImages.CheckedHover
		c.disabledImg = c.stateImages.CheckedDisabled
	} else {
		c.idleImg = c.stateImages.UncheckedIdle
		c.pressedImg = c.stateImages.UncheckedPressed
		c.hoverImg = c.stateImages.UncheckedHover
		c.disabledImg = c.stateImages.UncheckedDisabled
	}
	// Also ensure the interactive component's state is updated, e.g., if it was pressed.
	// This ensures it correctly reflects the new image.
	cx, cy := ebiten.CursorPosition()
	if c.state == ButtonPressed && !ContainsPoint(c.Bounds, cx, cy) {
		c.state = ButtonIdle // If it was pressed but cursor moved off, revert to idle.
	} else if c.state != ButtonPressed && ContainsPoint(c.Bounds, cx, cy) {
		c.state = ButtonHover // If not pressed but cursor is over, revert to hover.
	} else {
		c.state = ButtonIdle // Otherwise, revert to idle.
	}
}

// Update calls the embedded interactiveComponent's Update method.
func (c *Checkbox) Update() {
	c.interactiveComponent.Update()
}

// Draw draws the checkbox component using the image for its current visual state.
func (c *Checkbox) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.Bounds.Min.X), float64(c.Bounds.Min.Y))
	screen.DrawImage(c.GetCurrentStateImage(), op)
}

// HandlePress calls the embedded interactiveComponent's HandlePress method.
func (c *Checkbox) HandlePress() {
	c.interactiveComponent.HandlePress()
}

// HandleRelease calls the embedded interactiveComponent's HandleRelease method.
func (c *Checkbox) HandleRelease() {
	c.interactiveComponent.HandleRelease()
}

// HandleClick toggles the checkbox's state and calls the OnCheckChanged handler.
func (c *Checkbox) HandleClick() {
	if c.state != ButtonDisabled { // Only toggle if not disabled
		c.SetChecked(!c.Checked) // Toggle the checked state
		log.Printf("Checkbox '%s' clicked. New state: %t", c.Label, c.Checked)
	}
}
