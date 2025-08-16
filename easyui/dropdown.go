package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// DropDown represents a clickable component that reveals a menu when clicked.
type DropDown struct {
	component                                    // Embeds the base component struct
	Label                  string                // The text displayed on the drop-down button itself (e.g., "Select an animal")
	SelectedOption         string                // The currently selected option's label
	menu                   *Menu                 // The associated menu component
	state                  ButtonState           // To handle hover and pressed visual states (similar to a button)
	theme                  BareBonesTheme        // Reference to the theme for drawing
	uiGenerator            *BareBonesUiGenerator // To generate dropdown button images
	isPressedInside        bool                  // True if the mouse button was initially pressed while cursor was inside the button's bounds
	toggleHandledThisFrame bool                  // New: true if the menu was just toggled (shown/hidden) this frame
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
		Label:                  initialLabel,
		SelectedOption:         initialLabel, // Initially, the label is the selected option
		menu:                   menu,
		state:                  ButtonIdle,
		theme:                  theme,
		uiGenerator:            uiGen,
		isPressedInside:        false,
		toggleHandledThisFrame: false, // Initialize to false
	}
	return dd
}

// Update handles the interaction for the drop-down button.
func (d *DropDown) Update() {
	// Reset the toggleHandledThisFrame flag at the start of each update cycle
	defer func() {
		d.toggleHandledThisFrame = false
	}()

	// Get current cursor position
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := d.ContainsPoint(cx, cy)

	// Step 1: Handle mouse button just pressed
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorInBounds {
			d.state = ButtonPressed
			d.isPressedInside = true // Mark that the press originated inside
		} else {
			d.isPressedInside = false // Mark that the press originated outside
			// If click started outside, hide the menu if it's open.
			// This provides a way to close the menu by clicking anywhere else.
			if d.menu.isVisible {
				d.menu.Hide()
				d.toggleHandledThisFrame = true // Mark that a toggle happened
			}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Step 2: Handle mouse button currently held down
		if d.isPressedInside { // Only if the press originated inside the dropdown
			if cursorInBounds {
				d.state = ButtonPressed // Keep showing pressed if still over the dropdown
			} else {
				d.state = ButtonIdle // Revert to idle if mouse moves off while held
			}
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// Step 3: Handle mouse button just released
		// Only process if a toggle hasn't already been handled this frame (e.g., from an initial outside click closing the menu)
		if !d.toggleHandledThisFrame {
			if d.isPressedInside && cursorInBounds {
				// If released inside AND press originated inside, toggle menu visibility
				if d.menu.isVisible {
					d.menu.Hide()
				} else {
					// Position the menu below the dropdown
					d.menu.SetPosition(d.Bounds.Min.X, d.Bounds.Max.Y)
					d.menu.Show()
				}
				d.state = ButtonIdle            // Reset to idle after click
				d.toggleHandledThisFrame = true // Mark that a toggle happened
			} else {
				// If released outside, or released after an outside press, revert to idle/hover based on cursor position
				if cursorInBounds {
					d.state = ButtonHover
				} else {
					d.state = ButtonIdle
				}
			}
		}
		d.isPressedInside = false // Always reset this flag on release
	} else {
		// Step 4: Handle hover state when no mouse button is pressed
		if cursorInBounds {
			d.state = ButtonHover // Show hover if cursor is over and not pressed
		} else {
			d.state = ButtonIdle // Otherwise, return to idle
		}
	}
	// The menu's own Update method is handled by the Ui's modal logic.
	// No need to call d.menu.Update() here directly as it would be double-called.
}

// Draw draws the drop-down button and its current label.
func (d *DropDown) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(d.Bounds.Min.X), float64(d.Bounds.Min.Y))

	// Regenerate the button image every frame with the current 'SelectedOption' text.
	// We use generateDropdownImage for the desired styling, including the arrow.
	currentLabelImage := d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.PrimaryColor, d.theme.OnPrimaryColor, d.SelectedOption)
	screen.DrawImage(currentLabelImage, op) // Draw the image that reflects the selected option
}
