package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ui struct {
	component
}

func NewUi(x, y, width, height int) *Ui {
	return &Ui{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
	}
}

func (b *Ui) Update() {
	// go through all children and update
}

func (b *Ui) Draw(screen *ebiten.Image) {
	// draw all children (recursively)
}
