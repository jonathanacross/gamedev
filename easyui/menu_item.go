package main

import (
	"log" // Keep log for potential warnings, but it might be less needed now

	"github.com/hajimehoshi/ebiten/v2"
)

// MenuItem represents a single selectable item within a menu.
type MenuItem struct {
	interactiveComponent // Embed the new interactive component
	Label                string
	onClick              func()
}

// NewMenuItem creates a new MenuItem instance.
func NewMenuItem(x, y, width, height int, label string, idle, hover, pressed *ebiten.Image) *MenuItem {
	// Assuming `disabled` image for menu item is the same as idle for now, adjust as needed.
	disabled := idle
	return &MenuItem{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idle, pressed, hover, disabled),
		Label:                label,
		onClick:              nil, // Set via SetClickHandler
	}
}

// SetClickHandler sets the function to be executed when this menu item is clicked.
func (m *MenuItem) SetClickHandler(handler func()) {
	m.onClick = handler
}

// Update calls the embedded interactiveComponent's Update method.
func (m *MenuItem) Update() {
	m.interactiveComponent.Update()
}

// Draw draws the menu item using the image from its current state.
func (m *MenuItem) Draw(screen *ebiten.Image) {
	if m.idleImg == nil || m.hoverImg == nil || m.pressedImg == nil {
		log.Printf("MenuItem '%s': WARNING: One or more state images are nil! Idle: %t, Hover: %t, Pressed: %t", m.Label, m.idleImg != nil, m.hoverImg != nil, m.pressedImg != nil)
		// Consider drawing a placeholder or just returning if images are critical.
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
