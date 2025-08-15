package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonState int

const (
	ButtonIdle ButtonState = iota
	ButtonPressed
	ButtonHover
	ButtonDisabled
)

type Button struct {
	component

	idle     *ebiten.Image
	pressed  *ebiten.Image
	hover    *ebiten.Image
	disabled *ebiten.Image

	state ButtonState
}

func (b *Button) Update() {
	// TODO: implement
	// cx, cy := ebiten.CursorPosition()
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	b.isPressed = b.Bounds.Contains(float64(cx), float64(cy))
	// }
	// if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
	// 	if b.isPressed && b.Bounds.Contains(float64(cx), float64(cy)) {
	// 		b.isPressed = false
	// 		b.onclick()
	// 	}
	// }
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Bounds.Min.X), float64(b.Bounds.Min.Y))
	switch b.state {
	case ButtonIdle:
		screen.DrawImage(b.idle, op)
	case ButtonPressed:
		screen.DrawImage(b.pressed, op)
	case ButtonHover:
		screen.DrawImage(b.hover, op)
	case ButtonDisabled:
		screen.DrawImage(b.disabled, op)
	}
}
