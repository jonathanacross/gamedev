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
	// Removed: theme          BareBonesTheme        // No longer needed, access via uiGenerator.theme
	uiGenerator *BareBonesUiGenerator // To generate dropdown button images (still needed for image generation)
}

// NewDropDown function definition is now in ui_generator.go

// SetSelectedOption updates the displayed option and regenerates the dropdown button's images.
func (d *DropDown) SetSelectedOption(newOption string) {
	d.SelectedOption = newOption
	// Regenerate all state images with the new text, accessing theme via uiGenerator
	d.idleImg = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.uiGenerator.theme.PrimaryColor, d.uiGenerator.theme.OnPrimaryColor, d.SelectedOption)
	d.hoverImg = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.uiGenerator.theme.AccentColor, d.uiGenerator.theme.OnPrimaryColor, d.SelectedOption)
	d.pressedImg = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.uiGenerator.theme.AccentColor, d.uiGenerator.theme.OnPrimaryColor, d.SelectedOption)
	d.disabledImg = d.uiGenerator.generateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.uiGenerator.theme.PrimaryColor, d.uiGenerator.theme.OnPrimaryColor, d.SelectedOption)
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
