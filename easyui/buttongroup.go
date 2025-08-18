package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// LayoutType defines how components within a group are arranged.
type LayoutType int

const (
	LayoutVertical LayoutType = iota
	LayoutHorizontal
)

// SelectionMode defines the group's click behavior.
type SelectionMode int

const (
	SingleSelection SelectionMode = iota // Only one item can be selected at a time (e.g., radio buttons)
	Independent                          // Items can be toggled independently (e.g., checkboxes)
)

// checkable is an interface for components that have a checked/unchecked state,
// like Checkbox and RadioButton. This allows the ButtonGroup to work with them generically.
type checkable interface {
	Component
	SetChecked(bool)
	IsChecked() bool
}

// ButtonGroup represents a container for a group of interactive components,
// managing their layout and selection behavior.
type ButtonGroup struct {
	component
	LayoutType    LayoutType
	SelectionMode SelectionMode
	spacing       int
}

// NewButtonGroup creates a new ButtonGroup instance.
func NewButtonGroup(x, y, width, height int, layoutType LayoutType, selectionMode SelectionMode, spacing int) *ButtonGroup {
	bg := &ButtonGroup{
		LayoutType:    layoutType,
		SelectionMode: selectionMode,
		spacing:       spacing,
	}
	bg.component = NewComponent(x, y, width, height, bg)
	return bg
}

// AddChild adds a component to the group and automatically positions it.
func (bg *ButtonGroup) AddChild(c Component) {
	c.SetParent(bg.self)

	if len(bg.children) > 0 {
		lastChild := bg.children[len(bg.children)-1]
		childBounds := lastChild.GetBounds()

		switch bg.LayoutType {
		case LayoutVertical:
			newY := childBounds.Max.Y + bg.spacing
			c.SetPosition(childBounds.Min.X, newY)
		case LayoutHorizontal:
			newX := childBounds.Max.X + bg.spacing
			c.SetPosition(newX, childBounds.Min.Y)
		}
	} else {
		// For the first child, its position is relative to the group's bounds.
		c.SetPosition(0, 0)
	}

	bg.children = append(bg.children, c)
}

// Update iterates through the children and calls their Update methods.
func (bg *ButtonGroup) Update() {
	for _, child := range bg.children {
		child.Update()
	}
}

// Draw iterates through the children and calls their Draw methods.
func (bg *ButtonGroup) Draw(screen *ebiten.Image) {
	for _, child := range bg.children {
		child.Draw(screen)
	}
}

func (bg *ButtonGroup) FindChildAt(x, y int) Component {
	for _, child := range bg.children {
		if ContainsPoint(child, x, y) {
			return child
		}
	}

	// No component was found at the coordinates.
	return nil
}

// HandleClick handles the SingleSelection logic.
func (bg *ButtonGroup) HandleClick() {
	log.Println("ButtonGroup: handle click.")
	// Only proceed if the selection mode is SingleSelection.
	if bg.SelectionMode != SingleSelection {
		return
	}
	log.Println("ButtonGroup: SingleSelection click triggered.")

	cx, cy := ebiten.CursorPosition()
	clickedItem := bg.FindChildAt(cx, cy)

	log.Printf("ButtonGroup: clickedItem is %T.", clickedItem)

	// If no item was clicked, or if it's not a checkable component, return early.
	checkableItem, ok := clickedItem.(checkable)
	if !ok {
		log.Println("ButtonGroup: No checkable item was clicked.")
		return
	}

	// If the clicked item is already checked, there's nothing to do.
	if checkableItem.IsChecked() {
		log.Println("ButtonGroup: already checked.")
		return
	}

	// Otherwise, uncheck everything and check the new item
	log.Printf("ButtonGroup: SingleSelection click on %T.", clickedItem)
	for _, child := range bg.children {
		if checkable, ok := child.(checkable); ok {
			log.Printf("ButtonGroup: Unchecking checkable item %T.\n", checkable)
			checkable.SetChecked(false)
		} else {
			log.Printf("ButtonGroup: Not a checkable item %T.\n", child)
		}
	}
	checkableItem.SetChecked(true)
}
