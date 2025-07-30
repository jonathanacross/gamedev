package main

import (
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
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

func (t *Tile) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(t.X),
			Y: int(t.Y),
		},
		Max: image.Point{
			X: int(t.X + TileSize),
			Y: int(t.Y + TileSize),
		},
	}
}

type Item interface {
	GetX() float64
	GetY() float64
	Update()
	Draw(camera *Camera, screen *ebiten.Image)
	HitRect() image.Rectangle
	UseItem(g *Game)
}

type CoinItem struct {
	Location
	spriteSheet *SpriteSheet
	animation   *Animation
}

func (i *CoinItem) Update() {
	i.animation.Update()
}

func (p *CoinItem) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)
	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

func (e *CoinItem) UseItem(g *Game) {
	g.score++
}

func (e *CoinItem) GetX() float64 { return e.X }
func (e *CoinItem) GetY() float64 { return e.Y }
func (t *CoinItem) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(t.X),
			Y: int(t.Y),
		},
		Max: image.Point{
			X: int(t.X + TileSize),
			Y: int(t.Y + TileSize),
		},
	}
}

type HeartItem struct {
	Location
	spriteSheet *SpriteSheet
	animation   *Animation
}

func (i *HeartItem) Update() {
	i.animation.Update()
}

func (p *HeartItem) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)
	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

func (e *HeartItem) UseItem(g *Game) {
	g.player.AddHeath()
}

func (e *HeartItem) GetX() float64 { return e.X }
func (e *HeartItem) GetY() float64 { return e.Y }
func (t *HeartItem) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(t.X),
			Y: int(t.Y),
		},
		Max: image.Point{
			X: int(t.X + TileSize),
			Y: int(t.Y + TileSize),
		},
	}
}

// type Item struct {
// 	Location
// 	spriteSheet *SpriteSheet
// 	animation   *Animation
// }
//
// func (i *Item) Update() {
// 	i.animation.Update()
// }

// func (i *Item) Draw(camera *Camera, screen *ebiten.Image) {
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(i.X-camera.X, i.Y-camera.Y)
// 	subRect := i.spriteSheet.Rect(i.animation.Frame())
// 	currImage := i.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
// 	screen.DrawImage(currImage, op)
// }

// func (t *Item) HitRect() image.Rectangle {
// 	return image.Rectangle{
// 		Min: image.Point{
// 			X: int(t.X),
// 			Y: int(t.Y),
// 		},
// 		Max: image.Point{
// 			X: int(t.X + TileSize),
// 			Y: int(t.Y + TileSize),
// 		},
// 	}
// }

type Enemy interface {
	GetX() float64
	GetY() float64
	Update()
	Draw(camera *Camera, screen *ebiten.Image)
	HitRect() image.Rectangle
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

func (e *Octo) GetX() float64 { return e.X }
func (e *Octo) GetY() float64 { return e.Y }

func (t *Octo) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(t.X),
			Y: int(t.Y),
		},
		Max: image.Point{
			X: int(t.X + TileSize),
			Y: int(t.Y + TileSize),
		},
	}
}

type Bee struct {
	Location
	spriteSheet *SpriteSheet
	animation   *Animation
	speed       float64
}

func (e *Bee) Update() {
	e.X -= e.speed
	e.animation.Update()
}

func (e *Bee) GetX() float64 { return e.X }
func (e *Bee) GetY() float64 { return e.Y }

func (t *Bee) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(t.X),
			Y: int(t.Y),
		},
		Max: image.Point{
			X: int(t.X + TileSize),
			Y: int(t.Y + TileSize),
		},
	}
}

func (b *Bee) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X-camera.X, b.Y-camera.Y)
	subRect := b.spriteSheet.Rect(b.animation.Frame())
	currImage := b.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type Player struct {
	Location
	Velocity
	spriteSheet     SpriteSheet
	animation       *Animation
	maxHealth       int
	health          int
	invincible      bool
	invincibleFrame int
	invincibleTimer *Timer
}

func NewPlayer() *Player {
	return &Player{
		Location: Location{
			X: ScreenWidth / 2,
			Y: 3 * TileSize,
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
		animation:       NewAnimation(8, 15, 5),
		maxHealth:       PlayerMaxHealth,
		health:          PlayerMaxHealth,
		invincible:      false,
		invincibleFrame: 0,
		invincibleTimer: NewTimer(1500 * time.Millisecond),
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

	if p.invincible {
		p.invincibleFrame++
		p.invincibleTimer.Update()
		if p.invincibleTimer.IsReady() {
			p.invincible = false
		}
	}

	p.X += p.Dx
	p.Y += p.Dy

	p.animation.Update()
}

func (p *Player) Draw(camera *Camera, screen *ebiten.Image) {
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)

	cm := colorm.ColorM{}
	if p.invincible && ((p.invincibleFrame/5)%2 == 0) {
		// Change the color of the sprite to pure white
		cm.Translate(1.0, 1.0, 1.0, 0.0)
	}

	subRect := p.spriteSheet.Rect(p.animation.Frame())
	currImage := p.spriteSheet.image.SubImage(subRect).(*ebiten.Image)
	colorm.DrawImage(screen, currImage, cm, op)

	// Draw player life
	for i := range p.maxHealth {
		x := float64(2*ScoreOffset + i*HeartWidth)
		y := float64(ScoreOffset)
		filled := i < p.health
		DrawHeart(screen, x, y, filled)
	}
}

func (t *Player) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(t.X),
			Y: int(t.Y),
		},
		Max: image.Point{
			X: int(t.X + TileSize),
			Y: int(t.Y + TileSize),
		},
	}
}

func (p *Player) DoHit() {
	if p.invincible {
		return
	}

	p.health--

	p.invincible = true
	p.invincibleFrame = 0
	p.invincibleTimer.Reset()
}

func (p *Player) AddHeath() {
	p.health++
	if p.health > p.maxHealth {
		p.health = p.maxHealth
	}
}
