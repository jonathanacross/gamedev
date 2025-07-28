package main

import (
	"image"
	"math"

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
		speed: BackgroundScrollSpeed,
	}
}

func (b *Background) Update() {
	b.X -= b.speed
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

func (t *Tile) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.X-camera.X, t.Y-camera.Y)
	currImage := t.spriteSheet.image.SubImage(t.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type Item struct {
	Location
	spriteSheet *SpriteSheet
	animation   *Animation
}

func (i *Item) Update() {
	i.animation.Update()
}

func (i *Item) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(i.X-camera.X, i.Y-camera.Y)
	subRect := i.spriteSheet.Rect(i.animation.Frame())
	currImage := i.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type Enemy interface {
	Update()
	Draw(camera *Camera, screen *ebiten.Image)
}

type Octo struct {
	Location
	spriteSheet *SpriteSheet
	animation   *Animation
	minY        float64
	maxY        float64
	t           float64
	speed       float64
}

func (e *Octo) Update() {
	mid := (e.minY + e.maxY) / 2
	amplitude := (e.maxY - e.minY) / 2
	e.t++
	e.Y = mid + amplitude*math.Sin(e.t*e.speed)
	e.animation.Update()
}

func (p *Octo) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)
	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
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
			X: ScreenWidth / 2,
			Y: ScreenHeight / 2,
		},
		Velocity: Velocity{
			Dx: PlayerSpeed,
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Dy = JumpVelocity
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		p.Dx = PlayerSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		p.Dx = -PlayerSpeed
	}

	p.Dy += Gravity

	p.X += p.Dx
	p.Y += p.Dy

	p.animation.Update()
}

func (p *Player) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)
	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}
