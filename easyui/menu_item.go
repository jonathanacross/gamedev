package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// MenuItem represents a single selectable item within a menu.
type MenuItem struct {
	interactiveComponent // Embed the new interactive component
	Label                string
	onClick              func()
	uiGenerator          *BareBonesUiGenerator // Reference to the generator
	// Removed: theme        BareBonesTheme        // No longer needed, access via uiGenerator.theme
}

// NewMenuItem function definition is now in ui_generator.go

// SetClickHandler sets the function to be executed when this menu item is clicked.
func (m *MenuItem) SetClickHandler(handler func()) {
	m.onClick = handler
}

// SetLabel updates the menu item's text and regenerates its state images.
func (m *MenuItem) SetLabel(newLabel string) {
	m.Label = newLabel
	// Regenerate all state images with the new text, accessing theme via uiGenerator
	m.idleImg = m.uiGenerator.generateMenuItemImage(m.Bounds.Dx(), m.Bounds.Dy(), m.uiGenerator.theme.MenuColor, m.uiGenerator.theme.OnPrimaryColor, m.Label)
	m.hoverImg = m.uiGenerator.generateMenuItemImage(m.Bounds.Dx(), m.Bounds.Dy(), m.uiGenerator.theme.MenuItemHoverColor, m.uiGenerator.theme.OnPrimaryColor, m.Label)
	m.pressedImg = m.uiGenerator.generateMenuItemImage(m.Bounds.Dx(), m.Bounds.Dy(), m.uiGenerator.theme.AccentColor, m.uiGenerator.theme.OnPrimaryColor, m.Label)
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
