package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	// Required for image.Rectangle, image.Point
)

// Menu represents a pop-up menu containing a list of menu items.
// It now fully implements the Component interface.
type Menu struct {
	component                 // Embeds the base component struct
	items      []*MenuItem    // The list of items in the menu
	isVisible  bool           // Whether the menu is currently visible/open
	theme      BareBonesTheme // Reference to the theme for drawing (still needed for menu-specific items not covered by renderer states)
	renderer   UiRenderer     // Changed to UiRenderer interface
	background *ebiten.Image  // Image for the menu background
	parentUi   *Ui            // Reference to the root UI to manage modal state
	justOpened bool           // True if the menu was just opened this frame. Prevents immediate self-close.
}

// NewMenu creates a new Menu instance. It is now a standalone function.
func NewMenu(x, y, width int, theme BareBonesTheme, renderer UiRenderer, parentUi *Ui) *Menu {
	// Set an initial default height to prevent panic when generating initial background image.
	// This can be a small value, as it will be expanded by AddItem.
	const defaultInitialMenuHeight = 30 // A reasonable default height for an empty menu

	// Create the Menu first, then pass its pointer as 'self'
	m := &Menu{
		theme:    theme,
		renderer: renderer, // Store the renderer
		parentUi: parentUi,
		items:    []*MenuItem{}, // Initialize slice
	}
	m.component = NewComponent(x, y, width, defaultInitialMenuHeight, m) // Pass 'm' as self

	m.isVisible = false
	m.justOpened = false

	// Initial background image generation
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())
	log.Printf("Menu.NewMenu: Initial Bounds: %v (Dx: %d, Dy: %d)", m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic
	return m
}

// AddItem adds a new MenuItem to the menu.
// It now uses the standalone NewMenuItem function and the stored renderer.
func (m *Menu) AddItem(label string, handler func()) *MenuItem {
	itemHeight := 30 // Still hardcoded here, but can be pulled from theme later.
	itemWidth := m.Bounds.Dx()

	// Menu items' positions are relative to the menu itself.
	// Calculate relative Y offset for the new item.
	relativeYOffset := len(m.items) * itemHeight

	// Use the standalone NewMenuItem function and pass the stored renderer
	item := NewMenuItem(0, relativeYOffset, itemWidth, itemHeight, label, m.renderer) // x,y are 0, relativeYOffset
	item.SetClickHandler(func() {
		log.Printf("Menu: Anonymous handler for item '%s' called (from AddItem).", label) // Diagnostic log
		handler()                                                                         // Call the user-defined handler (from demo.go)
		m.Hide()                                                                          // Hide the menu
	})
	m.items = append(m.items, item)
	m.AddChild(item) // Add child to the menu, setting menu as parent

	// Update the menu's overall height based on added items
	// IMPORTANT: Preserve the current X dimensions, only update Y based on added items
	m.Bounds.Max.Y = m.Bounds.Min.Y + len(m.items)*itemHeight
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())                                                      // Regenerate background on size change
	log.Printf("Menu.AddItem: After adding '%s', Menu Bounds: %v (Dx: %d, Dy: %d)", label, m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic

	return item
}

// Update handles updating all child menu items.
// This method fully implements Component.Update().
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
// This method fully implements Component.Draw().
func (m *Menu) Draw(screen *ebiten.Image) {
	if !m.isVisible {
		return
	}

	absX, absY := m.GetAbsolutePosition() // Get absolute position of the menu itself
	if m.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(absX), float64(absY)) // Draw menu background at its absolute position
		screen.DrawImage(m.background, op)
	} else {
		log.Printf("Menu.Draw: WARNING: Menu background is nil! Cannot draw background for menu at bounds %v", m.Bounds)
	}

	for _, item := range m.items {
		item.Draw(screen) // Child items will draw themselves using their own absolute positions
	}
}

// Show makes the menu visible and sets it as the modal component in the UI.
func (m *Menu) Show() {
	m.isVisible = true
	m.justOpened = true
	log.Printf("Menu.Show: Menu set to visible and justOpened=true. Current Bounds: %v (Dx: %d, Dy: %d)", m.Bounds, m.Bounds.Dx(), m.Bounds.Dy()) // Diagnostic
	if m.parentUi != nil {
		m.parentUi.SetModal(m) // This now correctly passes a Component
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
// This function sets the menu's own relative position, and the children's positions
// will automatically follow due to the GetAbsolutePosition logic.
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
// This method fully implements Component.HandlePress().
func (m *Menu) HandlePress() {}

// HandleRelease is a no-op for the menu background.
// This method fully implements Component.HandleRelease().
func (m *Menu) HandleRelease() {}

// HandleClick is now primarily used by the Ui to handle clicks on the menu's background
// when it's a modal, closing the menu.
// This method fully implements Component.HandleClick().
func (m *Menu) HandleClick() {
	log.Printf("Menu: HandleClick called for menu (modal background). isVisible: %t", m.isVisible) // Diagnostic log
	// If a click lands on the menu's background (not an item), it should close the menu.
	m.Hide()
}
