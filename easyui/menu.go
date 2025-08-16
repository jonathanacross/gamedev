package main

import (
	"image"
	"log" // Import for logging

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil" // For checking clicks outside
)

// Menu represents a pop-up menu containing a list of menu items.
type Menu struct {
	component                         // Embeds the base component struct
	items       []*MenuItem           // The list of items in the menu
	isVisible   bool                  // Whether the menu is currently visible/open
	theme       BareBonesTheme        // Reference to the theme for drawing
	uiGenerator *BareBonesUiGenerator // To generate menu item images
	background  *ebiten.Image         // Image for the menu background
	parentUi    *Ui                   // Reference to the root UI to manage modal state
	justOpened  bool                  // True if the menu was just opened on the *previous* frame
}

// NewMenu creates a new Menu instance.
func NewMenu(x, y, width int, theme BareBonesTheme, uiGen *BareBonesUiGenerator, parentUi *Ui) *Menu {
	m := &Menu{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y}, // Max Y will be adjusted later
			},
		},
		items:       []*MenuItem{},
		isVisible:   false,
		theme:       theme,
		uiGenerator: uiGen,
		parentUi:    parentUi,
		justOpened:  false, // Initialize to false
	}
	log.Printf("NewMenu: Created menu at bounds %v", m.Bounds)
	return m
}

// AddItem adds a new MenuItem to the menu.
func (m *Menu) AddItem(label string, handler func()) *MenuItem {
	itemHeight := 30 // Fixed height for each menu item
	itemWidth := m.Bounds.Dx()

	// Calculate item bounds
	yOffset := m.Bounds.Min.Y + len(m.items)*itemHeight
	itemBounds := image.Rectangle{
		Min: image.Point{X: m.Bounds.Min.X, Y: yOffset},
		Max: image.Point{X: m.Bounds.Min.X + itemWidth, Y: yOffset + itemHeight},
	}

	// Generate images for the menu item states
	idleImg := m.uiGenerator.generateMenuItemImage(itemWidth, itemHeight, m.theme.MenuColor, m.theme.OnPrimaryColor, label)
	hoverImg := m.uiGenerator.generateMenuItemImage(itemWidth, itemHeight, m.theme.MenuItemHoverColor, m.theme.OnPrimaryColor, label)
	pressedImg := m.uiGenerator.generateMenuItemImage(itemWidth, itemHeight, m.theme.AccentColor, m.theme.OnPrimaryColor, label)

	item := NewMenuItem(itemBounds.Min.X, itemBounds.Min.Y, itemBounds.Dx(), itemBounds.Dy(), label, idleImg, hoverImg, pressedImg)
	item.SetClickHandler(func() {
		handler() // Call the user-defined handler
		m.Hide()  // Hide menu after selection
	})
	m.items = append(m.items, item)
	m.AddChild(item) // Add as a child component so Update/Draw propagate

	// Update the menu's overall bounds to accommodate new item
	m.Bounds.Max.Y = m.Bounds.Min.Y + len(m.items)*itemHeight
	// Regenerate background with the updated height
	m.background = m.uiGenerator.generateMenuImage(m.Bounds.Dx(), m.Bounds.Dy(), m.theme.MenuColor)
	log.Printf("Menu.AddItem: Added item '%s'. Menu now has %d items. Menu bounds updated to %v. Background regenerated.", label, len(m.items), m.Bounds)

	return item
}

// Update handles updating all child menu items and detecting clicks outside to close the menu.
func (m *Menu) Update() {
	log.Printf("Menu.Update: Called. isVisible: %t, justOpened: %t, Items count: %d", m.isVisible, m.justOpened, len(m.items))

	if !m.isVisible {
		return // Only update if visible
	}

	// Handle the grace period for a newly opened menu.
	// This frame, we *just* made it visible. We allow items to update
	// but prevent the menu from immediately closing due to the same click.
	if m.justOpened {
		log.Printf("Menu.Update: Menu was just opened, resetting justOpened flag and skipping external click check this frame.")
		m.justOpened = false // Reset for the next frame
		// Still update children to allow immediate hover effects for the first item
		for _, item := range m.items {
			log.Printf("Menu.Update: Calling Update for MenuItem '%s'", item.Label)
			item.Update()
		}
		return // Skip further input processing this frame for the menu itself
	}

	// Check for clicks outside the menu to close it
	// Only respond to a *new* mouse press that occurs outside the menu.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		if !m.ContainsPoint(cx, cy) {
			log.Printf("Menu.Update: Click outside menu detected at (%d,%d). Menu bounds: %v. Hiding menu.", cx, cy, m.Bounds)
			m.Hide()
			return // Menu is hidden, no need to update items
		}
	}

	log.Printf("Menu.Update: Updating %d menu items.", len(m.items))
	for _, item := range m.items {
		log.Printf("Menu.Update: Calling Update for MenuItem '%s'", item.Label)
		item.Update()
	}
}

// Draw draws the menu background and all its items.
func (m *Menu) Draw(screen *ebiten.Image) {
	log.Printf("Menu.Draw: Called. isVisible: %t, Items count: %d", m.isVisible, len(m.items))

	if !m.isVisible {
		return
	}

	log.Printf("Menu.Draw: Attempting to draw menu background at bounds %v", m.Bounds)

	// Draw the menu's background
	if m.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(m.Bounds.Min.X), float64(m.Bounds.Min.Y))
		screen.DrawImage(m.background, op)
		log.Printf("Menu.Draw: Background drawn. Background image bounds: %v", m.background.Bounds())
	} else {
		log.Printf("Menu.Draw: WARNING: Menu background is nil! Cannot draw background.")
	}

	log.Printf("Menu.Draw: Drawing %d menu items.", len(m.items))
	// Draw all child menu items
	for _, item := range m.items {
		log.Printf("Menu.Draw: Calling Draw for MenuItem '%s'. Item bounds: %v", item.Label, item.Bounds)
		item.Draw(screen)
	}
}

// Show makes the menu visible and sets it as the modal component in the UI.
func (m *Menu) Show() {
	m.isVisible = true
	m.justOpened = true // Set to true when shown
	log.Printf("Menu.Show: Menu set to visible and justOpened=true. Current Bounds: %v", m.Bounds)
	if m.parentUi != nil {
		m.parentUi.SetModal(m) // Set self as modal
		log.Printf("Menu.Show: Menu set as modal component in parent Ui.")
	}
}

// Hide makes the menu invisible and clears it as the modal component in the UI.
func (m *Menu) Hide() {
	m.isVisible = false
	// m.justOpened remains false; it's only set to true in Show()
	log.Printf("Menu.Hide: Menu set to invisible.")
	if m.parentUi != nil {
		m.parentUi.ClearModal() // Clear self as modal
		log.Printf("Menu.Hide: Modal component cleared in parent Ui.")
	}
}

// SetPosition allows dynamic repositioning of the menu and its items.
func (m *Menu) SetPosition(x, y int) {
	log.Printf("Menu.SetPosition: Attempting to set position from (%d,%d) to (%d,%d)", m.Bounds.Min.X, m.Bounds.Min.Y, x, y)
	diffX := x - m.Bounds.Min.X
	diffY := y - m.Bounds.Min.Y

	m.Bounds.Min.X = x
	m.Bounds.Min.Y = y
	m.Bounds.Max.X += diffX
	m.Bounds.Max.Y += diffY

	// Regenerate background to match the new size and ensure it's correctly placed
	m.background = m.uiGenerator.generateMenuImage(m.Bounds.Dx(), m.Bounds.Dy(), m.theme.MenuColor)
	log.Printf("Menu.SetPosition: Menu repositioned. New bounds: %v. Background regenerated.", m.Bounds)

	for _, item := range m.items {
		item.Bounds.Min.X += diffX
		item.Bounds.Min.Y += diffY
		item.Bounds.Max.X += diffX
		item.Bounds.Max.Y += diffY
		log.Printf("Menu.SetPosition: Item '%s' repositioned to %v", item.Label, item.Bounds)
	}
}
