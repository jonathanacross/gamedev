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

// Entity interface defines common behaviors for all game objects.
type Entity interface {
	Update()
	Draw(camera *Camera, screen *ebiten.Image)
	HitRect() image.Rectangle
	GetX() float64
	GetY() float64
}

// BaseSprite provides common fields and methods for any visible game entity.
// It handles drawing a single sprite or the current frame of an animation.
type BaseSprite struct {
	Location
	spriteSheet *SpriteSheet
	srcRect     image.Rectangle // The specific rectangle on the sprite sheet to draw
}

// NewBaseSprite is a helper to initialize common fields for a static sprite.
func NewBaseSprite(x, y float64, ss *SpriteSheet, srcRect image.Rectangle) BaseSprite {
	return BaseSprite{
		Location:    Location{X: x, Y: y},
		spriteSheet: ss,
		srcRect:     srcRect,
	}
}

// Draw renders the BaseSprite on the screen, adjusted by the camera's position.
func (bs *BaseSprite) Draw(camera *Camera, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X-camera.X, bs.Y-camera.Y)
	currImage := bs.spriteSheet.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

// HitRect returns the collision rectangle for the BaseSprite.
func (bs *BaseSprite) HitRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(bs.X), Y: int(bs.Y)},
		Max: image.Point{
			X: int(bs.X + float64(bs.spriteSheet.tileWidth)),
			Y: int(bs.Y + float64(bs.spriteSheet.tileHeight))},
	}
}

// GetX returns the X coordinate of the BaseSprite.
func (bs *BaseSprite) GetX() float64 { return bs.X }

// GetY returns the Y coordinate of the BaseSprite.
func (bs *BaseSprite) GetY() float64 { return bs.Y }

// Update is a no-op for BaseSprite as it's static.
func (bs *BaseSprite) Update() {}

// AnimatedSprite embeds BaseSprite and adds animation capabilities.
type AnimatedSprite struct {
	BaseSprite
	animation *Animation
}

// NewAnimatedSprite is a helper to initialize fields for an animated sprite.
func NewAnimatedSprite(x, y float64, ss *SpriteSheet, anim *Animation) AnimatedSprite {
	// Initialize srcRect with the first frame of the animation
	initialSrcRect := ss.Rect(anim.Frame())
	return AnimatedSprite{
		BaseSprite: NewBaseSprite(x, y, ss, initialSrcRect),
		animation:  anim,
	}
}

// Update advances the animation and updates the source rectangle.
func (as *AnimatedSprite) Update() {
	as.animation.Update()
	as.srcRect = as.spriteSheet.Rect(as.animation.Frame())
}

// PhysicsEntity used for physics-driven entities.
type PhysicsEntity struct {
	AnimatedSprite
	Velocity
}

// NewPhysicsEntity is a helper to initialize a PhysicsEntity.
func NewPhysicsEntity(x, y, dx, dy float64, ss *SpriteSheet, anim *Animation) PhysicsEntity {
	return PhysicsEntity{
		AnimatedSprite: NewAnimatedSprite(x, y, ss, anim),
		Velocity:       Velocity{Dx: dx, Dy: dy},
	}
}

// Background handles the scrolling background.
type Background struct {
	Location
	image *ebiten.Image
	width int
	speed float64
}

// NewBackground creates a new Background instance.
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

// Update scrolls the background.
func (b *Background) Update() {
	b.X -= b.speed
	if b.X < -float64(b.width) {
		b.X = 0
	}
}

// Draw renders the background, tiling it to create a continuous scroll.
func (b *Background) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.image, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(b.X+float64(b.width), b.Y)
	screen.DrawImage(b.image, op2)
}

// Tile represents a static environmental tile.
type Tile struct {
	BaseSprite
}

// Player represents the main player character.
type Player struct {
	PhysicsEntity
	maxHealth       int
	health          int
	invincible      bool
	invincibleFrame int
	invincibleTimer *Timer
}

// NewPlayer creates a new Player instance.
func NewPlayer() *Player {
	playerSS := NewSpriteSheet(PlayerImage, 16, 16, 8, 3)
	anim := NewAnimation(8, 15, 5)

	p := &Player{
		PhysicsEntity: NewPhysicsEntity(
			ScreenWidth/2, 3*TileSize, // Initial position
			PlayerSpeed, 0, // Initial velocity
			playerSS, anim, // Sprite sheet and animation
		),
		maxHealth:       PlayerMaxHealth,
		health:          PlayerMaxHealth,
		invincible:      false,
		invincibleFrame: 0,
		invincibleTimer: NewTimer(1500 * time.Millisecond),
	}
	// Initial update to set the correct srcRect for the first frame
	p.AnimatedSprite.Update()
	return p
}

// Update handles player input, physics, and invincibility.
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

	p.Dy += Gravity // Apply gravity

	if p.invincible {
		p.invincibleFrame++
		p.invincibleTimer.Update()
		if p.invincibleTimer.IsReady() {
			p.invincible = false
		}
	}

	p.X += p.Dx // Apply velocity to position
	p.Y += p.Dy

	p.AnimatedSprite.Update()
}

// Draw renders the player, applying invincibility flicker effect.
func (p *Player) Draw(camera *Camera, screen *ebiten.Image) {
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)

	cm := colorm.ColorM{}
	if p.invincible && ((p.invincibleFrame/5)%2 == 0) {
		// Change the color of the sprite to pure white for flicker effect
		cm.Translate(1.0, 1.0, 1.0, 0.0)
	}

	// Use srcRect from the embedded BaseSprite
	currImage := p.spriteSheet.image.SubImage(p.srcRect).(*ebiten.Image)
	colorm.DrawImage(screen, currImage, cm, op)

	// Draw player life hearts
	for i := range p.maxHealth {
		x := float64(2*ScoreOffset + i*HeartWidth)
		y := float64(ScoreOffset)
		filled := i < p.health
		DrawHeart(screen, x, y, filled)
	}
}

// DoHit reduces player health and activates invincibility.
func (p *Player) DoHit() {
	if p.invincible {
		return
	}

	p.health--

	p.invincible = true
	p.invincibleFrame = 0
	p.invincibleTimer.Reset()
}

// AddHeath increases player health, up to maxHealth.
func (p *Player) AddHeath() {
	p.health++
	if p.health > p.maxHealth {
		p.health = p.maxHealth
	}
}

// Item interface defines behaviors for collectible items.
type Item interface {
	Entity
	UseItem(g *Game)
}

// CoinItem represents a collectible coin.
type CoinItem struct {
	AnimatedSprite
}

// UseItem increments the game score when collected.
func (i *CoinItem) UseItem(g *Game) {
	g.score++
}

// HeartItem represents a collectible heart that restores player health.
type HeartItem struct {
	AnimatedSprite
}

// UseItem restores player health when collected.
func (i *HeartItem) UseItem(g *Game) {
	g.player.AddHeath()
}

// Enemy interface defines behaviors for enemies.
type Enemy interface {
	Entity
}

// Octo represents an octopus enemy that moves sinusoidally.
type Octo struct {
	AnimatedSprite
	minY  float64
	maxY  float64
	t     float64 // Time parameter for sinusoidal movement
	speed float64 // Speed of sinusoidal movement
}

// Update calculates the Octo's vertical movement and updates its animation.
func (e *Octo) Update() {
	mid := (e.minY + e.maxY) / 2
	amplitude := (e.maxY - e.minY) / 2
	e.t++
	e.Y = mid + amplitude*math.Sin(e.t*e.speed)
	e.AnimatedSprite.Update()
}

// Bee represents a bee enemy that moves horizontally.
type Bee struct {
	AnimatedSprite
	speed float64
}

// Update moves the Bee horizontally and updates its animation.
func (e *Bee) Update() {
	e.X -= e.speed
	e.AnimatedSprite.Update()
}
