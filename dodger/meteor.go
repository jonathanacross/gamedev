package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Meteor struct {
	sprite *ebiten.Image
	width  float64
	height float64

	position Vector
	angle    float64
	speed    float64
	rotSpeed float64
}

func NewMeteor() *Meteor {
	speed := 2.0 + rand.Float64()*2.0
	rotSpeed := 0.2 * (-0.5 + rand.Float64())
	sprite := MeteorSprites[rand.IntN(len(MeteorSprites))]

	width := float64(sprite.Bounds().Dx())
	height := float64(sprite.Bounds().Dy())

	pos := Vector{
		X: rand.Float64()*ScreenWidth - width/2,
		Y: -height,
	}

	return &Meteor{
		sprite:   sprite,
		width:    width,
		height:   height,
		position: pos,
		speed:    speed,
		rotSpeed: rotSpeed,
	}
}

func (p *Meteor) Update() error {
	p.position.Y += p.speed
	p.angle += p.rotSpeed
	return nil
}

func (p *Meteor) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-p.width/2, -p.height/2)
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(p.sprite, op)
}

func (m *Meteor) HasFallenOffscreen() bool {
	return m.position.Y > ScreenHeight+m.height
}

func (p *Meteor) HitRect() Rect {
	return Rect{
		Left:   p.position.X - p.width/2,
		Top:    p.position.Y - p.height/2,
		Right:  p.position.X + p.width/2,
		Bottom: p.position.Y + p.height/2,
	}
}
