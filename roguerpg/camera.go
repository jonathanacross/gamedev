package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	Location Location
	Width    float64
	Height   float64
	BufferX  float64
	BufferY  float64
}

func NewCamera(width, height float64) *Camera {
	return &Camera{
		Width:   width,
		Height:  height,
		BufferX: width / 4,
		BufferY: height / 4,
	}
}

func (c *Camera) CenterOn(loc Location) {
	// TODO: Smooth camera movement
	c.Location.X = loc.X
	c.Location.Y = loc.Y
}

func (c *Camera) GetViewRect() Rect {
	return Rect{
		Left:   c.Location.X - c.Width/2,
		Top:    c.Location.Y - c.Height/2,
		Right:  c.Location.X + c.Width/2,
		Bottom: c.Location.Y + c.Height/2,
	}
}

func (c *Camera) WorldToScreen() ebiten.GeoM {
	offsetX := c.Location.X - c.Width/2
	offsetY := c.Location.Y - c.Height/2

	// Create the matrix and apply the inverse translation (subtract the camera's offset)
	m := ebiten.GeoM{}
	m.Translate(-offsetX, -offsetY)

	return m
}
