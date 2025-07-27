package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Location struct {
	X float64
	Y float64
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
	image   *ebiten.Image
	isSolid bool
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
	spriteSheet SpriteSheet
	animation   *Animation
}

func NewPlayer() *Player {
	return &Player{
		Location: Location{
			X: 50,
			Y: 50,
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
	p.animation.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)
	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}
