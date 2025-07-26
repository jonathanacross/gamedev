package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	sprite *ebiten.Image
	width  float64
	height float64

	position Vector
	speed    float64
}

func NewBullet(position Vector) *Bullet {
	speed := 8.0
	sprite := BulletSprite

	width := float64(sprite.Bounds().Dx())
	height := float64(sprite.Bounds().Dy())

	return &Bullet{
		sprite:   sprite,
		width:    width,
		height:   height,
		position: position,
		speed:    speed,
	}
}

func (b *Bullet) Update() error {
	b.position.Y -= b.speed
	return nil
}

func (p *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X-p.width/2, p.position.Y-p.height/2)
	screen.DrawImage(p.sprite, op)
}

func (m *Bullet) HasFallenOffscreen() bool {
	return m.position.Y < 0
}

func (p *Bullet) HitRect() Rect {
	return Rect{
		Left:   p.position.X - p.width/2,
		Top:    p.position.Y - p.height/2,
		Right:  p.position.X + p.width/2,
		Bottom: p.position.Y + p.height/2,
	}
}
