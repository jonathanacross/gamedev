package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Draw(screen *ebiten.Image)
	Update()
	GetBounds() image.Rectangle
}

type component struct {
	Bounds   image.Rectangle
	children []Component
}

func (c component) GetBounds() image.Rectangle {
	return c.Bounds
}

func (c component) AddChild(child Component) {
	c.children = append(c.children, child)
}
