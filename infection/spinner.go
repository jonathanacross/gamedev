package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Spinner struct {
	image   *ebiten.Image
	x       int
	y       int
	theta   float64
	visible bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		image:   SpinnerImage,
		x:       50,
		y:       ScreenHeight - 100,
		theta:   0,
		visible: false,
	}
}

func (s *Spinner) IsVisible() bool { return s.visible }

func (s *Spinner) SetVisible(visible bool) {
	s.visible = visible
}

func (s *Spinner) Update() {
	s.theta += 0.01
}

func (s *Spinner) Draw(screen *ebiten.Image) {
	if !s.visible {
		return
	}
	width := float64(s.image.Bounds().Dx())
	height := float64(s.image.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-width/2, -height/2)
	op.GeoM.Rotate(s.theta)
	op.GeoM.Translate(float64(s.x), float64(s.y))
	screen.DrawImage(s.image, op)
}
