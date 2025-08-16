package main

import (
	"image"

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
}

// NewMenu creates a new Menu instance.
func NewMenu(x, y, width int, theme BareBonesTheme, uiGen *BareBonesUiGenerator, parentUi *Ui) *Menu {
	// Initial height is 0, it will be calculated when items are added
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
	}
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

	// Update the menu's overall bounds
	m.Bounds.Max.Y = m.Bounds.Min.Y + len(m.items)*itemHeight
	m.background = m.uiGenerator.generateMenuImage(m.Bounds.Dx(), m.Bounds.Dy(), m.theme.MenuColor) // Regenerate background

	return item
}

// Update handles updating all child menu items and detecting clicks outside to close the menu.
func (m *Menu) Update() {
	if !m.isVisible {
		return // Only update if visible
	}

	// Check for clicks outside the menu to close it
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		if !m.ContainsPoint(cx, cy) {
			m.Hide()
			return // Menu is hidden, no need to update items
		}
	}

	for _, item := range m.items {
		item.Update()
	}
}

// Draw draws the menu background and all its items.
func (m *Menu) Draw(screen *ebiten.Image) {
	if !m.isVisible {
		return
	}

	// Draw the menu's background
	if m.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(m.Bounds.Min.X), float64(m.Bounds.Min.Y))
		screen.DrawImage(m.background, op)
	}

	// Draw all child menu items
	for _, item := range m.items {
		item.Draw(screen)
	}
}

// Show makes the menu visible and sets it as the modal component in the UI.
func (m *Menu) Show() {
	m.isVisible = true
	if m.parentUi != nil {
		m.parentUi.SetModal(m) // Set self as modal
	}
}

// Hide makes the menu invisible and clears it as the modal component in the UI.
func (m *Menu) Hide() {
	m.isVisible = false
	if m.parentUi != nil {
		m.parentUi.ClearModal() // Clear self as modal
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

	for _, item := range m.items {
		item.Bounds.Min.X += diffX
		item.Bounds.Min.Y += diffY
		item.Bounds.Max.X += diffX
		item.Bounds.Max.Y += diffY
	}
}
