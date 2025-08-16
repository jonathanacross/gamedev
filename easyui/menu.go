package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	justOpened  bool                  // True if the menu was just opened this frame. Prevents immediate self-close.
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
		justOpened:  false,
	}
	return m
}

// AddItem adds a new MenuItem to the menu.
func (m *Menu) AddItem(label string, handler func()) *MenuItem {
	itemHeight := 30
	itemWidth := m.Bounds.Dx()

	yOffset := m.Bounds.Min.Y + len(m.items)*itemHeight
	itemBounds := image.Rectangle{
		Min: image.Point{X: m.Bounds.Min.X, Y: yOffset},
		Max: image.Point{X: m.Bounds.Min.X + itemWidth, Y: yOffset + itemHeight},
	}

	idleImg := m.uiGenerator.generateMenuItemImage(itemWidth, itemHeight, m.theme.MenuColor, m.theme.OnPrimaryColor, label)
	hoverImg := m.uiGenerator.generateMenuItemImage(itemWidth, itemHeight, m.theme.MenuItemHoverColor, m.theme.OnPrimaryColor, label)
	pressedImg := m.uiGenerator.generateMenuItemImage(itemWidth, itemHeight, m.theme.AccentColor, m.theme.OnPrimaryColor, label)

	item := NewMenuItem(itemBounds.Min.X, itemBounds.Min.Y, itemBounds.Dx(), itemBounds.Dy(), label, idleImg, hoverImg, pressedImg)
	item.SetClickHandler(func() {
		handler() // Call the user-defined handler
		m.Hide()  // Now, the menu item's handler should also hide the menu.
	})
	m.items = append(m.items, item)
	m.AddChild(item)

	m.Bounds.Max.Y = m.Bounds.Min.Y + len(m.items)*itemHeight
	m.background = m.uiGenerator.generateMenuImage(m.Bounds.Dx(), m.Bounds.Dy(), m.theme.MenuColor)

	return item
}

// Update handles updating all child menu items and detecting clicks outside to close the menu.
func (m *Menu) Update() {
	if !m.isVisible {
		return
	}

	for _, item := range m.items {
		item.Update()
	}

	if m.justOpened {
		m.justOpened = false
		return
	}

	// This logic remains, as the menu still needs to close if the user clicks outside.
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		if !ContainsPoint(m.Bounds, cx, cy) {
			log.Printf("Menu.Update: Click outside menu detected. Hiding menu.")
			m.Hide()
			return
		}
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

	m.background = m.uiGenerator.generateMenuImage(m.Bounds.Dx(), m.Bounds.Dy(), m.theme.MenuColor)

	for _, item := range m.items {
		item.Bounds.Min.X += diffX
		item.Bounds.Min.Y += diffY
		item.Bounds.Max.X += diffX
		item.Bounds.Max.Y += diffY
	}
}

// HandleClick allows the menu to conform to the Component interface.
func (m *Menu) HandleClick() {
	// A menu itself doesn't have a click action, but its items do.
	// The centralized logic in ui.go will now call the HandleClick of the specific item.
}
