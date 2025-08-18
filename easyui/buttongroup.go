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
		c.SetPosition(0, 0)
	}

	bg.children = append(bg.children, c)
}

// Update calls Update on all child components.
func (bg *ButtonGroup) Update() {
	for _, child := range bg.children {
		child.Update()
	}
}

// Draw draws the button group and its children.
func (bg *ButtonGroup) Draw(screen *ebiten.Image) {
	for _, child := range bg.children {
		child.Draw(screen)
	}
}

// HandleChildClick is called by a child component to delegate the click logic up to the group.
func (bg *ButtonGroup) HandleChildClick(clickedItem Component) {
	log.Printf("ButtonGroup: Handling click from child of type %T.", clickedItem)

	// Only proceed if the selection mode is SingleSelection.
	if bg.SelectionMode != SingleSelection {
		return
	}

	// Uncheck everything except for the item that was just clicked.
	for _, child := range bg.children {
		if checkable, ok := child.(checkable); ok {
			if child != clickedItem {
				checkable.SetChecked(false)
			}
		}
	}
}
