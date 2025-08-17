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

	// Create the MenuItem first, then pass its pointer as 'self'
	m := &MenuItem{
		Label:    label,
		onClick:  nil,      // Click handler set separately
		renderer: renderer, // Store the renderer
	}
	m.interactiveComponent = NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg, m) // Pass 'm' as self

	return m
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
// It now calls the embedded interactiveComponent's Draw method.
func (m *MenuItem) Draw(screen *ebiten.Image) {
	if m.idleImg == nil || m.hoverImg == nil || m.pressedImg == nil {
		log.Printf("MenuItem '%s': WARNING: One or more state images are nil! Idle: %t, Hover: %t, Pressed: %t", m.Label, m.idleImg != nil, m.hoverImg != nil, m.pressedImg != nil)
	}
	m.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (m *MenuItem) HandlePress() {
	m.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
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
