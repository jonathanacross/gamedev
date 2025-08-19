package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// UiRenderer defines the interface for rendering all UI components.
// Implementations will map component states to specific visual outputs (e.g., colors, textures).
type UiRenderer interface {
	// GenerateButtonImage creates an image for a button in a specific state.
	GenerateButtonImage(width, height int, text string, icon image.Image, state ButtonState, isChecked bool) *ebiten.Image

	// GenerateDropdownImage creates an image for a dropdown button in a specific state.
	GenerateDropdownImage(width, height int, text string, state ButtonState) *ebiten.Image

	// GenerateMenuItemImage creates an image for a menu item in a specific state.
	GenerateMenuItemImage(width, height int, text string, state ButtonState) *ebiten.Image

	// GenerateMenuImage creates an image for the menu's background.
	GenerateMenuImage(width, height int) *ebiten.Image

	// GenerateCheckboxImage creates an image for a checkbox in a specific state and checked status.
	// `componentState` refers to ButtonState (Idle, Hover, Pressed, Disabled)
	// `isChecked` refers to the checkbox's boolean checked/unchecked status
	GenerateCheckboxImage(width, height int, label string, componentState ButtonState, isChecked bool) *ebiten.Image

	// GenerateRadioButtonImage creates an image for a checkbox in a specific state and checked status.
	// `componentState` refers to ButtonState (Idle, Hover, Pressed, Disabled)
	// `isChecked` refers to the checkbox's boolean checked/unchecked status
	GenerateRadioButtonImage(width, height int, label string, componentState ButtonState, isChecked bool) *ebiten.Image

	// GenerateTextFieldImage creates an image for a text field in a specific state, with its current text and cursor.
	// `showCursor` is a boolean indicating whether the cursor should be drawn.
	GenerateTextFieldImage(width, height int, text string, componentState ButtonState, isFocused bool, cursorPos int, showCursor bool) *ebiten.Image

	// GenerateLabelImage creates an image for a static text label.
	GenerateLabelImage(width, height int, text string) *ebiten.Image

	// GenerateContainerImage creates an image for a container's background.
	GenerateContainerImage(width, height int) *ebiten.Image
}
