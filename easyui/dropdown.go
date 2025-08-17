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
	renderer             UiRenderer // Changed to UiRenderer interface
}

// NewDropDown creates a new DropDown instance, generating its state-specific images.
// It is now a standalone function.
func NewDropDown(x, y, width, height int, initialLabel string, menu *Menu, renderer UiRenderer) *DropDown {
	// Generate specific dropdown images using the renderer's methods
	idleImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonIdle)
	hoverImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonHover)
	pressedImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonPressed)
	disabledImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonDisabled)

	return &DropDown{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg),
		Label:                initialLabel,
		SelectedOption:       initialLabel,
		menu:                 menu,
		renderer:             renderer, // Store the renderer
	}
}

// SetSelectedOption updates the displayed option and regenerates the dropdown button's images.
func (d *DropDown) SetSelectedOption(newOption string) {
	d.SelectedOption = newOption
	// Regenerate all state images with the new text using the renderer
	d.idleImg = d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonIdle)
	d.hoverImg = d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonHover)
	d.pressedImg = d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonPressed)
	d.disabledImg = d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonDisabled)
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
