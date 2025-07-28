package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Location struct {
	X float64
	Y float64
}

type Velocity struct {
	Dx float64
	Dy float64
}

type Background struct {
	Location
	image *ebiten.Image
	width int
	speed float64
}

func NewBackground() *Background {
	return &Background{
		Location: Location{
			X: 0,
			Y: 0,
		},
		image: BackgroundImage,
		width: BackgroundImage.Bounds().Dx(),
		speed: 2,
	}
}

func (b *Background) Update() {
	b.X -= 0.25
	if b.X < -float64(b.width) {
		b.X = 0
	}
}

func (b *Background) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.image, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(b.X+float64(b.width), b.Y)
	screen.DrawImage(b.image, op2)
}

type Tile struct {
	Location
	spriteSheet *SpriteSheet
	srcRect     image.Rectangle
}

func (t *Tile) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.X, t.Y)
	currImage := t.spriteSheet.image.SubImage(t.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type Item struct {
	Location
	image *ebiten.Image
}

type Enemy struct {
	Location
	spriteSheet *SpriteSheet
	animation   *Animation
}

type Player struct {
	Location
	Velocity
	spriteSheet SpriteSheet
	animation   *Animation
}

func NewPlayer() *Player {
	return &Player{
		Location: Location{
			X: 50,
			Y: 50,
		},
		Velocity: Velocity{
			Dx: 0,
			Dy: 0,
		},
		spriteSheet: SpriteSheet{
			image:         PlayerImage,
			tileWidth:     16,
			tileHeight:    16,
			widthInTiles:  8,
			heightInTiles: 3,
		},
		animation: NewAnimation(8, 15, 5),
	}
}

func (p *Player) Update() {
	const gravity = 0.1
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Dy = -2.5
	}

	p.Dy += gravity

	p.X += p.Dx
	p.Y += p.Dy

	p.animation.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)
	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}
