package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// DropDown represents a clickable component that reveals a menu when clicked.
type DropDown struct {
	interactiveComponent // Embed the new interactive component
	Label                string
	SelectedOption       string
	menu                 *Menu
	theme                BareBonesTheme        // Reference to the theme for drawing (still needed for image generation)
	uiGenerator          *BareBonesUiGenerator // To generate dropdown button images (still needed for image generation)
}

// NewDropDown creates a new DropDown instance, associating it with a menu.
func NewDropDown(x, y, width, height int, initialLabel string, menu *Menu, theme BareBonesTheme, uiGen *BareBonesUiGenerator) *DropDown {
	// Generate specific dropdown images using the uiGenerator
	idleImg := uiGen.generateDropdownImage(width, height, theme.PrimaryColor, theme.OnPrimaryColor, initialLabel)
	hoverImg := uiGen.generateDropdownImage(width, height, theme.AccentColor, theme.OnPrimaryColor, initialLabel)
	pressedImg := uiGen.generateDropdownImage(width, height, theme.AccentColor, theme.OnPrimaryColor, initialLabel)   // Often same as hover for dropdown pressed
	disabledImg := uiGen.generateDropdownImage(width, height, theme.PrimaryColor, theme.OnPrimaryColor, initialLabel) // Example: darker version

	return &DropDown{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg),
		Label:                initialLabel,
		SelectedOption:       initialLabel,
		menu:                 menu,
		theme:                theme,
		uiGenerator:          uiGen,
	}
}

// Update calls the embedded interactiveComponent's Update method.
func (d *DropDown) Update() {
	d.interactiveComponent.Update()
}

// Draw draws the dropdown button using the image from its current state.
func (d *DropDown) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(d.Bounds.Min.X), float64(d.Bounds.Min.Y))
	screen.DrawImage(d.GetCurrentStateImage(), op)
}

// HandlePress calls the embedded interactiveComponent's HandlePress method.
func (d *DropDown) HandlePress() {
	d.interactiveComponent.HandlePress()
}

// HandleRelease calls the embedded interactiveComponent's HandleRelease method.
func (d *DropDown) HandleRelease() {
	d.interactiveComponent.HandleRelease()
}

// HandleClick manages the dropdown menu visibility.
func (d *DropDown) HandleClick() {
	if d.state == ButtonDisabled { // Do not respond if disabled
		return
	}

	if d.menu.parentUi != nil && d.menu.parentUi.modalComponent == d.menu {
		d.menu.Hide() // If our menu is currently the modal, close it.
	} else if d.menu.parentUi != nil && d.menu.parentUi.modalComponent == nil {
		d.menu.SetPosition(d.Bounds.Min.X, d.Bounds.Max.Y)
		d.menu.Show() // Show the menu if no other modal is active.
	}
}
