package main

import (
	"log"

	"image" // Required for image.Rectangle, image.Point

	"github.com/hajimehoshi/ebiten/v2"
)

// Menu represents a pop-up menu containing a list of menu items.
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
	m := &Menu{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + defaultInitialMenuHeight}, // Set initial height here
			},
		},
		items:      []*MenuItem{},
		isVisible:  false,
		theme:      theme,
		renderer:   renderer, // Store the renderer
		parentUi:   parentUi,
		justOpened: false,
	}
	// Initial background image generation with a non-zero height
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy())
	return m
}

// AddItem adds a new MenuItem to the menu.
// It now uses the standalone NewMenuItem function and the stored renderer.
func (m *Menu) AddItem(label string, handler func()) *MenuItem {
	itemHeight := 30 // Still hardcoded here, but can be pulled from theme later.
	itemWidth := m.Bounds.Dx()

	yOffset := m.Bounds.Min.Y + len(m.items)*itemHeight

	// Use the standalone NewMenuItem function and pass the stored renderer
	item := NewMenuItem(m.Bounds.Min.X, yOffset, itemWidth, itemHeight, label, m.renderer)
	item.SetClickHandler(func() {
		handler() // Call the user-defined handler
		m.Hide()  // Now, the menu item's handler should also hide the menu.
	})
	m.items = append(m.items, item)
	m.AddChild(item)

	// Update the menu's overall height based on added items
	m.Bounds.Max.Y = m.Bounds.Min.Y + len(m.items)*itemHeight
	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy()) // Regenerate background on size change

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

	if m.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(m.Bounds.Min.X), float64(m.Bounds.Min.Y))
		screen.DrawImage(m.background, op)
	} else {
		log.Printf("Menu.Draw: WARNING: Menu background is nil! Cannot draw background for menu at bounds %v", m.Bounds)
	}

	for _, item := range m.items {
		item.Draw(screen)
	}
}

// Show makes the menu visible and sets it as the modal component in the UI.
func (m *Menu) Show() {
	m.isVisible = true
	m.justOpened = true
	log.Printf("Menu.Show: Menu set to visible and justOpened=true. Current Bounds: %v", m.Bounds)
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
	diffX := x - m.Bounds.Min.X
	diffY := y - m.Bounds.Min.Y

	m.Bounds.Min.X = x
	m.Bounds.Min.Y = y
	m.Bounds.Max.X += diffX
	m.Bounds.Max.Y += diffY

	m.background = m.renderer.GenerateMenuImage(m.Bounds.Dx(), m.Bounds.Dy()) // Regenerate background on reposition

	for _, item := range m.items {
		item.Bounds.Min.X += diffX
		item.Bounds.Min.Y += diffY
		item.Bounds.Max.X += diffX
		item.Bounds.Max.Y += diffY
	}
}

// HandlePress is a dummy method for the base struct.
func (m *Menu) HandlePress() {}

// HandleClick is now primarily used by the Ui to handle clicks on the menu's background
// when it's a modal.
func (m *Menu) HandleClick() {
	// If a click lands on the menu's background (not an item), it should close the menu.
	m.Hide()
}
