package main

import (
	"image"
)

type Camera struct {
	Location // top left of the view
	Width    float64
	Height   float64
	Buffer   float64
}

func NewCamera() *Camera {
	return &Camera{
		Location: Location{
			X: 0,
			Y: 0,
		},
		Width:  ScreenWidth,
		Height: ScreenHeight,
	}
}

func (c *Camera) GetViewRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(c.X),
			Y: int(c.Y),
		},
		Max: image.Point{
			X: int(c.X + c.Width),
			Y: int(c.Y + c.Height),
		},
	}
}

func (c *Camera) Center(loc Location) {
	c.X = loc.X
	c.Y = loc.Y
}
