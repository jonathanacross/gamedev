package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// DropDown represents a clickable component that reveals a menu when clicked.
type DropDown struct {
	component                            // Embeds the base component struct
	Label          string                // The text displayed on the drop-down button itself (e.g., "Select an animal")
	SelectedOption string                // The currently selected option's label
	menu           *Menu                 // The associated menu component
	state          ButtonState           // To handle hover and pressed visual states (similar to a button)
	theme          BareBonesTheme        // Reference to the theme for drawing
	uiGenerator    *BareBonesUiGenerator // To generate dropdown button images
}

// NewDropDown creates a new DropDown instance, associating it with a menu.
func NewDropDown(x, y, width, height int, initialLabel string, menu *Menu, theme BareBonesTheme, uiGen *BareBonesUiGenerator) *DropDown {
	dd := &DropDown{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		Label:          initialLabel,
		SelectedOption: initialLabel,
		menu:           menu,
		state:          ButtonIdle,
		theme:          theme,
		uiGenerator:    uiGen,
	}
	return dd
}

// Update handles the interaction for the drop-down button.
// It ensures the pressed state only persists if the mouse is over the button.
func (d *DropDown) Update() {
	if d.state == ButtonDisabled {
		return
	}

	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(d.Bounds, cx, cy)

	if d.state == ButtonPressed {
		// If currently pressed, check if mouse moved *off* the dropdown.
		// If so, switch to ButtonIdle (visual feedback: no longer highlighted as pressed).
		if !cursorInBounds {
			d.state = ButtonIdle // Reset to idle if mouse moves away while pressed
		}
		// If it's ButtonPressed and cursor is still in bounds, keep it ButtonPressed.
		return // Do not apply normal hover logic while actively pressed.
	}

	// Standard hover logic for Idle/Hover states
	if cursorInBounds {
		d.state = ButtonHover
	} else {
		d.state = ButtonIdle
	}
}

// Draw draws the drop-down button and its current label.
func (d *DropDown) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(d.Bounds.Min.X), float64(d.Bounds.Min.Y))

	// Determine image to draw based on state
	var imgToDraw *ebiten.Image
	switch d.state {
	case ButtonIdle:
		imgToDraw = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.PrimaryColor, d.theme.OnPrimaryColor, d.SelectedOption)
	case ButtonPressed:
		imgToDraw = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.AccentColor, d.theme.OnPrimaryColor, d.SelectedOption)
	case ButtonHover:
		imgToDraw = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.AccentColor, d.theme.OnPrimaryColor, d.SelectedOption)
	default:
		imgToDraw = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.PrimaryColor, d.theme.OnPrimaryColor, d.SelectedOption)
	}

	screen.DrawImage(imgToDraw, op)
}

// HandlePress is called when the mouse button is pressed down on this dropdown.
func (d *DropDown) HandlePress() {
	d.state = ButtonPressed
}

// HandleRelease is called when the mouse button is released, if this dropdown was the pressed component.
// It ensures the dropdown state transitions out of ButtonPressed to either Hover or Idle.
func (d *DropDown) HandleRelease() {
	cx, cy := ebiten.CursorPosition()
	if ContainsPoint(d.Bounds, cx, cy) {
		d.state = ButtonHover // Mouse released over dropdown
	} else {
		d.state = ButtonIdle // Mouse released away from dropdown
	}
}

// HandleClick is called when the mouse button is released on this dropdown,
// and it was also pressed on this dropdown.
func (d *DropDown) HandleClick() {
	// Toggle menu visibility based on current modal state.
	// We check parentUi.modalComponent to determine if *our* menu is the current modal.
	if d.menu.parentUi != nil && d.menu.parentUi.modalComponent == d.menu {
		d.menu.Hide() // If our menu is currently the modal, close it.
	} else if d.menu.parentUi != nil && d.menu.parentUi.modalComponent == nil {
		// If no modal is currently active, show our menu.
		d.menu.SetPosition(d.Bounds.Min.X, d.Bounds.Max.Y)
		d.menu.Show()
	}
	// State transition (from ButtonPressed to Hover/Idle) is handled by HandleRelease.
}
