package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Player struct {
	position Vector
	width    float64
	height   float64
	sprite   *ebiten.Image
}

func NewPlayer() *Player {
	sprite := PlayerSprite

	width := float64(sprite.Bounds().Dx())
	height := float64(sprite.Bounds().Dy())

	pos := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight * 0.8,
	}

	return &Player{
		position: pos,
		width:    width,
		height:   height,
		sprite:   sprite,
	}
}

func (p *Player) HitRect() Rect {
	return Rect{
		Left:   p.position.X - p.width/2,
		Top:    p.position.Y - p.height/2,
		Right:  p.position.X + p.width/2,
		Bottom: p.position.Y + p.height/2,
	}
}

func (p *Player) Update() error {
	speed := 5.0

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.position.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.position.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.position.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.position.X += speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.addBullet()
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X-p.width/2, p.position.Y-p.height/2)
	screen.DrawImage(p.sprite, op)
}

func (p *Player) addBullet() {
	bulletStart := Vector{
		X: p.position.X,
		Y: p.position.Y - p.height/2,
	}
	gameInstance.bullets = append(gameInstance.bullets, NewBullet(bulletStart))
}
