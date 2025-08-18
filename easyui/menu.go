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
	renderer   UiRenderer
	background *ebiten.Image
	parentUi   *Ui
	justOpened bool // True if the menu was just opened this frame. Prevents immediate self-close.
}

// NewMenu creates a new Menu instance.
func NewMenu(x, y, width int, renderer UiRenderer, parentUi *Ui) *Menu {
	// Set an initial default height to prevent panic when generating initial background image.
	// This can be a small value, as it will be expanded by AddItem.
	const defaultInitialMenuHeight = 30

	// Create the Menu first, then pass its pointer as 'self'
	m := &Menu{
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

// AddItem creates a new MenuItem and adds it to the menu, resizing the menu to fit.
func (m *Menu) AddItem(label string, onClick func()) {
	// Create the new menu item, positioning it below the previous one.
	// x and y are relative to the menu's top-left corner.
	var itemY int
	if len(m.items) > 0 {
		lastItem := m.items[len(m.items)-1]
		itemY = lastItem.Bounds.Max.Y + 2 // Add a small padding
	} else {
		itemY = 2 // Small top padding for the first item
	}
	// Menu items should match the width of the menu
	width := m.Bounds.Dx() - 4 // Small padding
	item := NewMenuItem(2, itemY, width, 25, label, m.renderer)
	item.SetClickHandler(onClick)

	// Add the item to the menu's children and internal item list
	m.AddChild(item)
	m.items = append(m.items, item)

	// Adjust the menu's size to accommodate the new item.
	newHeight := item.Bounds.Max.Y + 2 // Add a small bottom padding
	m.Bounds.Max.Y = m.Bounds.Min.Y + newHeight
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())          // Regenerate background to fit
	log.Printf("Menu.AddItem: Added item '%s'. New menu height: %d", label, newHeight) // Diagnostic
}

// Update iterates through all menu items and updates them if the menu is visible.
func (m *Menu) Update() {
	if !m.isVisible {
		return
	}
	// Clear the justOpened flag after the first frame
	m.justOpened = false

	for _, item := range m.items {
		item.Update()
	}
}

// Draw draws the menu's background and then its items if the menu is visible.
func (m *Menu) Draw(screen *ebiten.Image) {
	if !m.isVisible {
		return
	}

	absX, absY := m.GetAbsolutePosition()

	// Draw the menu's background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(m.background, op)

	// Draw all child menu items
	for _, item := range m.children {
		item.Draw(screen)
	}
}

// HandleClick passes the click to the correct menu item if the click is within the menu.
func (m *Menu) HandleClick() {
	// If a click happens on the menu background, it's a no-op. The click is handled by the items.
	// If the click is outside the menu, the modal system in ui.go will call this
	// and we should close the menu.
	if m.isVisible && !m.justOpened {
		log.Println("Menu.HandleClick: Closing menu because a click was registered outside its bounds.") // Diagnostic
		m.Hide()
	}
}

// Show makes the menu visible and sets it as the modal component in the UI.
func (m *Menu) Show() {
	log.Printf("Menu.Show: Parent UI: %v", m.parentUi) // Diagnostic
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
