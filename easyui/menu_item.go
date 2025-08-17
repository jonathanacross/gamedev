package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	// Required for image.Rectangle, image.Point
)

// MenuItem represents a single selectable item within a menu.
type MenuItem struct {
	interactiveComponent // Embed the new interactive component
	Label                string
	onClick              func()
	renderer             UiRenderer // Changed to UiRenderer interface
}

// NewMenuItem creates a new MenuItem instance.
// It is now a standalone function.
func NewMenuItem(x, y, width, height int, label string, renderer UiRenderer) *MenuItem {
	idleImg := renderer.GenerateMenuItemImage(width, height, label, ButtonIdle)
	hoverImg := renderer.GenerateMenuItemImage(width, height, label, ButtonHover)
	pressedImg := renderer.GenerateMenuItemImage(width, height, label, ButtonPressed)
	disabledImg := idleImg // Default disabled to idle, can be customized

	return &MenuItem{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg),
		Label:                label,
		onClick:              nil,      // Click handler set separately
		renderer:             renderer, // Store the renderer
	}
}

// SetClickHandler sets the function to be executed when this menu item is clicked.
func (m *MenuItem) SetClickHandler(handler func()) {
	m.onClick = handler
}

// SetLabel updates the menu item's text and regenerates its state images.
func (m *MenuItem) SetLabel(newLabel string) {
	m.Label = newLabel
	// Regenerate all state images with the new text using the renderer
	m.idleImg = m.renderer.GenerateMenuItemImage(m.Bounds.Dx(), m.Bounds.Dy(), m.Label, ButtonIdle)
	m.hoverImg = m.renderer.GenerateMenuItemImage(m.Bounds.Dx(), m.Bounds.Dy(), m.Label, ButtonHover)
	m.pressedImg = m.renderer.GenerateMenuItemImage(m.Bounds.Dx(), m.Bounds.Dy(), m.Label, ButtonPressed)
	m.disabledImg = m.idleImg // Assuming disabled is same as idle for now.
}

// Update calls the embedded interactiveComponent's Update method.
func (m *MenuItem) Update() {
	m.interactiveComponent.Update()
}

// Draw draws the menu item using the image from its current state.
func (m *MenuItem) Draw(screen *ebiten.Image) {
	if m.idleImg == nil || m.hoverImg == nil || m.pressedImg == nil {
		log.Printf("MenuItem '%s': WARNING: One or more state images are nil! Idle: %t, Hover: %t, Pressed: %t", m.Label, m.idleImg != nil, m.hoverImg != nil, m.pressedImg != nil)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.Bounds.Min.X), float64(m.Bounds.Min.Y))
	screen.DrawImage(m.GetCurrentStateImage(), op)
}

// HandlePress calls the embedded interactiveComponent's HandlePress method.
func (m *MenuItem) HandlePress() {
	m.interactiveComponent.HandlePress()
}

// HandleRelease calls the embedded interactiveComponent's HandleRelease method.
func (m *MenuItem) HandleRelease() {
	m.interactiveComponent.HandleRelease()
}

// HandleClick calls the menu item's onClick handler.
func (m *MenuItem) HandleClick() {
	if m.state != ButtonDisabled && m.onClick != nil {
		log.Printf("MenuItem '%s': Click handler triggered.", m.Label)
		m.onClick()
	}
}
