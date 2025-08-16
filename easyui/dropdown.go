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
func (d *DropDown) Update() {
	cx, cy := ebiten.CursorPosition()
	cursorInBounds := ContainsPoint(d.Bounds, cx, cy)

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

	currentLabelImage := d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.theme.PrimaryColor, d.theme.OnPrimaryColor, d.SelectedOption)
	screen.DrawImage(currentLabelImage, op)
}

// HandleClick calls the button's onClick handler.
func (d *DropDown) HandleClick() {
	if d.menu.isVisible {
		d.menu.Hide()
	} else {
		d.menu.SetPosition(d.Bounds.Min.X, d.Bounds.Max.Y)
		d.menu.Show()
	}
}
