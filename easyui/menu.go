package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Menu represents a pop-up menu containing a list of menu items.
type Menu struct {
	component
	items      []*MenuItem
	isVisible  bool
	theme      ShapeTheme // TODO: is this needed?
	renderer   UiRenderer
	background *ebiten.Image
	parentUi   *Ui
	justOpened bool // True if the menu was just opened this frame. Prevents immediate self-close.
}

// NewMenu creates a new Menu instance.
func NewMenu(x, y, width int, theme ShapeTheme, renderer UiRenderer, parentUi *Ui) *Menu {
	// Set an initial default height to prevent panic when generating initial background image.
	// This can be a small value, as it will be expanded by AddItem.
	const defaultInitialMenuHeight = 30

	// Create the Menu first, then pass its pointer as 'self'
	m := &Menu{
		theme:    theme,
		renderer: renderer,
		parentUi: parentUi,
		items:    []*MenuItem{},
	}
	m.component = NewComponent(x, y, width, defaultInitialMenuHeight, m)

	m.isVisible = false
	m.justOpened = false

	// Initial background image generation
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())
	log.Printf("Menu.NewMenu: Initial Bounds: %v (Dx: %d, Dy: %d)", m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic
	return m
}

// AddItem adds a new MenuItem to the menu.
func (m *Menu) AddItem(label string, handler func()) *MenuItem {
	itemHeight := 30
	itemWidth := m.Bounds.Dx()

	// Menu items' positions are relative to the menu itself.
	// Calculate relative Y offset for the new item.
	relativeYOffset := len(m.items) * itemHeight

	// Use the standalone NewMenuItem function and pass the stored renderer
	item := NewMenuItem(0, relativeYOffset, itemWidth, itemHeight, label, m.renderer)
	item.SetClickHandler(func() {
		log.Printf("Menu: Anonymous handler for item '%s' called (from AddItem).", label)
		handler()
		m.Hide()
	})
	m.items = append(m.items, item)
	m.AddChild(item)

	// Update the menu's overall height based on added items
	// Preserve the current X dimensions, only update Y based on added items
	m.Bounds.Max.Y = m.Bounds.Min.Y + len(m.items)*itemHeight
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())                                                      // Regenerate background on size change
	log.Printf("Menu.AddItem: After adding '%s', Menu Bounds: %v (Dx: %d, Dy: %d)", label, m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic

	return item
}

// Update handles updating all child menu items.
func (m *Menu) Update() {
	if !m.isVisible {
		return
	}

	for _, item := range m.items {
		item.Update()
	}

	if m.justOpened {
		m.justOpened = false
	}
}

// Draw draws the menu background and all its items.
func (m *Menu) Draw(screen *ebiten.Image) {
	if !m.isVisible {
		return
	}

	absX, absY := m.GetAbsolutePosition()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(m.background, op)

	for _, item := range m.items {
		item.Draw(screen)
	}
}

// Show makes the menu visible and sets it as the modal component in the UI.
func (m *Menu) Show() {
	m.isVisible = true
	m.justOpened = true
	log.Printf("Menu.Show: Menu set to visible and justOpened=true. Current Bounds: %v (Dx: %d, Dy: %d)", m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic
	if m.parentUi != nil {
		m.parentUi.SetModal(m)
	}
}

// Hide makes the menu invisible and clears it as the modal component in the UI.
func (m *Menu) Hide() {
	m.isVisible = false
	if m.parentUi != nil {
		m.parentUi.ClearModal()
	}
}

// SetPosition allows dynamic repositioning of the menu and its items.
func (m *Menu) SetPosition(x, y int) {
	// Capture current dimensions before modifying Min.X/Y
	currentWidth := m.Bounds.Dx()
	currentHeight := m.Bounds.Dy()

	m.Bounds.Min.X = x
	m.Bounds.Min.Y = y
	m.Bounds.Max.X = x + currentWidth
	m.Bounds.Max.Y = y + currentHeight

	// Regenerate background image if size changed (though not typical with SetPosition)
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())
	log.Printf("Menu.SetPosition: Menu position set to (%d, %d). New Bounds: %v (Dx: %d, Dy: %d)", x, y, m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic
}

// HandlePress is a no-op for the menu background.
func (m *Menu) HandlePress() {}

// HandleRelease is a no-op for the menu background.
func (m *Menu) HandleRelease() {}

// HandleClick is used by the Ui to handle clicks on the menu's background
// when it's a modal, closing the menu.
func (m *Menu) HandleClick() {
	log.Printf("Menu: HandleClick called for menu (modal background). isVisible: %t", m.isVisible)
	m.Hide()
}
