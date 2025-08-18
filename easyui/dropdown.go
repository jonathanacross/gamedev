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
func NewDropDown(x, y, width, height int, initialLabel string, menu *Menu, renderer UiRenderer) *DropDown {
	// Generate specific dropdown images using the renderer's methods
	idleImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonIdle)
	hoverImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonHover)
	pressedImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonPressed)
	disabledImg := renderer.GenerateDropdownImage(width, height, initialLabel, ButtonDisabled)

	// Create the DropDown first, then pass its pointer as 'self'
	d := &DropDown{
		Label:          initialLabel,
		SelectedOption: initialLabel,
		menu:           menu,
		renderer:       renderer, // Store the renderer
	}
	d.interactiveComponent = NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg, d) // Pass 'd' as self

	return d
}

// SetSelectedOption updates the displayed option and regenerates the dropdown button's images.
func (d *DropDown) SetSelectedOption(newOption string) {
	d.SelectedOption = newOption
	// Regenerate all state images with the new text using the renderer
	newIdleImg := d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonIdle)
	newHoverImg := d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonHover)
	newPressedImg := d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonPressed)
	newDisabledImg := d.renderer.GenerateDropdownImage(d.Bounds.Dx(), d.Bounds.Dy(), d.SelectedOption, ButtonDisabled)

	// Explicitly update the images held by the embedded interactiveComponent
	d.interactiveComponent.idleImg = newIdleImg
	d.interactiveComponent.hoverImg = newHoverImg
	d.interactiveComponent.pressedImg = newPressedImg
	d.interactiveComponent.disabledImg = newDisabledImg
}

// Update calls the embedded interactiveComponent's Update method.
func (d *DropDown) Update() {
	d.interactiveComponent.Update()
}

// Draw draws the dropdown button using the image from its current state.
func (d *DropDown) Draw(screen *ebiten.Image) {
	d.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (d *DropDown) HandlePress() {
	d.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
func (d *DropDown) HandleRelease() {
	d.interactiveComponent.HandleRelease()
}

// HandleClick manages the dropdown menu visibility.
func (d *DropDown) HandleClick() {
	if d.state == ButtonDisabled {
		return
	}

	if d.menu.parentUi != nil && d.menu.parentUi.modalComponent == Component(d.menu) {
		d.menu.Hide() // If our menu is currently the modal, close it.
	} else if d.menu.parentUi != nil && d.menu.parentUi.modalComponent == nil {
		absX, absY := d.GetAbsolutePosition()
		d.menu.SetPosition(absX, absY+d.Bounds.Dy()) // Menu appears directly below the dropdown
		d.menu.Show()
	}
}

// Focus is a no-op for a DropDown.
func (d *DropDown) Focus() {}

// Unfocus is a no-op for a DropDown.
func (d *DropDown) Unfocus() {}
