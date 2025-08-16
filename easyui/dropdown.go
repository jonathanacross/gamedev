package main

import (
	"fmt"
	"image" // Needed for gg placeholder

	// Needed for gg placeholder

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// DropDown represents a clickable component that reveals a menu when clicked.
type DropDown struct {
	component                      // Embeds the base component struct
	Label           string         // The text displayed on the drop-down button itself (e.g., "Select an animal")
	SelectedOption  string         // The currently selected option's label
	menu            *Menu          // The associated menu component
	state           ButtonState    // To handle hover and pressed visual states (similar to a button)
	theme           BareBonesTheme // Reference to the theme for drawing
	idleImage       *ebiten.Image  // These are generated initially based on Label, but not directly used for rendering
	hoverImage      *ebiten.Image  // because the button image is regenerated dynamically based on SelectedOption.
	pressedImage    *ebiten.Image
	uiGenerator     *BareBonesUiGenerator // To generate dropdown button images
	isPressedInside bool                  // true if the mouse button was initially pressed while cursor was inside the button's bounds
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
		Label:           initialLabel,
		SelectedOption:  initialLabel, // Initially, the label is the selected option
		menu:            menu,
		state:           ButtonIdle,
		theme:           theme,
		uiGenerator:     uiGen,
		isPressedInside: false,
	}

	// Generate images for the dropdown button itself based on the initial label.
	// Note: These images are not directly used in the Draw method as it regenerates
	// the image with the SelectedOption each frame. This might be a point for optimization
	// if performance becomes an issue for very frequent redraws, by only regenerating
	// when SelectedOption changes.
	dd.idleImage = uiGen.generateDropdownImage(width, height, theme.PrimaryColor, theme.OnPrimaryColor, dd.Label)
	dd.hoverImage = uiGen.generateDropdownImage(width, height, theme.PrimaryColor, theme.OnPrimaryColor, dd.Label)  // Could be a distinct hover color
	dd.pressedImage = uiGen.generateDropdownImage(width, height, theme.AccentColor, theme.OnPrimaryColor, dd.Label) // Accent color for pressed

	return dd
}

// Update handles the interaction for the drop-down button.
func (d *DropDown) Update() {
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := d.ContainsPoint(cx, cy)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Println("mouse pressed")
		if cursorInBounds {
			d.state = ButtonPressed
			if d.menu.isVisible {
				fmt.Println("hiding menu")
				d.menu.Hide()
			} else {
				fmt.Println("showing menu")
				// Position the menu below the dropdown
				d.menu.SetPosition(d.Bounds.Min.X, d.Bounds.Max.Y)
				d.menu.Show()
			}
		} else {
			fmt.Printf("state = %v\n", d.state)
			fmt.Printf("cursorinbounds = %v\n", cursorInBounds)

			fmt.Println("state != pressed or cursor not in bounds")
			d.state = ButtonIdle
		}
		d.isPressedInside = false
	} else {
		if cursorInBounds && !d.isPressedInside {
			d.state = ButtonHover
		} else {
			d.state = ButtonIdle
		}
	}
	// The menu's own Update method is handled by the Ui's modal logic.
	// No need to call d.menu.Update() here directly as it would be double-called.
}

// Draw draws the drop-down button and its current label.
func (d *DropDown) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(d.Bounds.Min.X), float64(d.Bounds.Min.Y))

	// The logic here regenerates the button image every frame with the
	// current 'SelectedOption' text. This is why the 'idleImage', 'hoverImage',
	// and 'pressedImage' fields are currently not directly used for rendering
	// the dynamic text.
	currentLabelImage := d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.PrimaryColor, d.theme.OnPrimaryColor, d.SelectedOption)
	screen.DrawImage(currentLabelImage, op) // Draw the image that reflects the selected option
}
