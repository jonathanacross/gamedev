package main

import (
	"math"
)

type FrogState int

const (
	Idle = iota
	Jumping
	Surprised
	Dying
)

type Frog struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animations  map[FrogState]*Animation
	state       FrogState
	jumpStartX  float64
	jumpTargetX float64
}

func NewFrog() *Frog {
	spriteSheet := NewSpriteSheet(32, 32, 5, 4)
	animations := map[FrogState]*Animation{
		Idle:      NewAnimation([]int{0, 1}, 15, true),
		Jumping:   NewAnimation([]int{10, 11, 12, 13, 14}, 10, false),
		Surprised: NewAnimation([]int{5, 6}, 10, false),
		Dying:     NewAnimation([]int{15, 16, 17, 18, 19}, 10, false),
	}

	frog := &Frog{
		BaseSprite: BaseSprite{
			Location: Location{X: 0, Y: 0},
			image:    FrogSpriteSheet,
			srcRect:  spriteSheet.Rect(0),
			hitbox:   Rect{left: 12, top: 12, right: 24, bottom: 24},
		},
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       Idle,
	}

	return frog
}

// Update the frog's animation and position.
func (f *Frog) Update() {
	animation := f.animations[f.state]
	animation.Update()
	f.srcRect = f.spriteSheet.Rect(animation.Frame())

	// Only update position if in a jumping state
	if f.state == Jumping {
		jumpAnimation := f.animations[Jumping]
		totalAnimationFrames := float64(len(jumpAnimation.frames) * jumpAnimation.speed)
		currentFrameCount := float64(jumpAnimation.frameIndex*jumpAnimation.speed) + float64(jumpAnimation.speed-jumpAnimation.frameCounter)
		progress := currentFrameCount / totalAnimationFrames

		jumpDistance := f.jumpTargetX - f.jumpStartX
		f.X = f.jumpStartX + jumpDistance*progress
		f.Y = float64(PlatformY-FrogOffsetY) - 15*math.Sin(math.Pi*progress)
	}
}

// Jump initiates a jump to a target X coordinate.
func (f *Frog) Jump(targetX float64) {
	if f.state != Idle {
		return // Can only jump from an idle state
	}
	f.state = Jumping
	f.animations[Jumping].Reset()
	f.jumpStartX = f.X
	f.jumpTargetX = targetX
}

// IsJumping checks if the frog is in the middle of a jump.
func (f *Frog) IsJumping() bool {
	return f.state == Jumping
}

// IsJumpFinished checks if the jump animation has completed.
func (f *Frog) IsJumpFinished() bool {
	return f.animations[Jumping].IsFinished()
}

// Land is called when the jump is complete to reset the state and position.
func (f *Frog) Land() {
	f.state = Idle
	f.Y = float64(PlatformY - FrogOffsetY)
	f.X = f.jumpTargetX
}

// Hit changes the frog's state to Dying.
func (f *Frog) Hit() {
	f.state = Dying
	f.animations[Dying].Reset()
}

// IsDyingFinished returns true if the dying animation has completed.
func (f *Frog) IsDyingFinished() bool {
	return f.animations[Dying].IsFinished()
}

// Rock and Platform structs are unchanged
type Platform struct {
	BaseSprite
}

func NewPlatform(x, y float64, end bool) *Platform {
	if end {
		return &Platform{
			BaseSprite: BaseSprite{
				Location: Location{X: x, Y: y},
				image:    EndPlatformSprite,
				srcRect:  EndPlatformSprite.Bounds(),
				hitbox:   NewRect(EndPlatformSprite.Bounds()),
			},
		}
	} else {
		return &Platform{
			BaseSprite: BaseSprite{
				Location: Location{X: x, Y: y},
				image:    PlatformSprite,
				srcRect:  PlatformSprite.Bounds(),
				hitbox:   NewRect(PlatformSprite.Bounds()),
			},
		}
	}
}

// Boot struct to represent the falling obstacle
type Boot struct {
	BaseSprite
	velocityY float64
	minY      float64
	maxY      float64
}

// NewBoot creates a new boot instance
func NewBoot(x, y float64) *Boot {
	return &Boot{
		BaseSprite: BaseSprite{
			Location: Location{X: x, Y: y},
			image:    BootSprite,
			srcRect:  BootSprite.Bounds(),
			hitbox:   NewRect(BootSprite.Bounds()),
		},
		velocityY: FallDownVelocity,
		minY:      FallingItemTopY,
		maxY:      PlatformY,
	}
}

// Update handles the boot's vertical movement
func (b *Boot) Update() {
	b.Y += b.velocityY

	// If the boot falls off-screen, reset its position and direction
	if b.Y > b.maxY {
		b.Y = b.maxY
		b.velocityY = FallUpVelocity
	} else if b.Y < b.minY {
		b.Y = b.minY
		b.velocityY = FallDownVelocity
	}
}

type Heron struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation
	velocityX   float64
	velocityY   float64
	minY        float64
	maxY        float64
	targetX     float64
	targetY     float64
}

// NewHeron creates a new heron instance
func NewHeron(targetX, targetY float64) *Heron {
	spriteSheet := NewSpriteSheet(48, 32, 1, 4)

	startPos := Location{X: ScreenWidth, Y: -100} // Start off-screen

	heron := &Heron{
		BaseSprite: BaseSprite{
			Location: startPos,
			image:    HeronSpriteSheet,
			srcRect:  spriteSheet.Rect(0),
			hitbox:   Rect{left: 0, top: 10, right: 15, bottom: 20},
		},
		spriteSheet: spriteSheet,
		animation:   NewAnimation([]int{0, 1, 2, 3}, 10, true),
		targetX:     targetX,
		targetY:     targetY,
	}

	// Calculate velocity to reach target
	speed := 2.0 // Adjust speed as needed
	dx := targetX - heron.X
	dy := targetY - heron.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	heron.velocityX = (dx / distance) * speed
	heron.velocityY = (dy / distance) * speed

	return heron
}

// Update handles the heron's movement and animation.
func (h *Heron) Update() {
	h.animation.Update()
	h.srcRect = h.spriteSheet.Rect(h.animation.Frame())

	if math.Abs(h.X-h.targetX) < math.Abs(h.velocityX) &&
		math.Abs(h.Y-h.targetY) < math.Abs(h.velocityY) {
		// Reached target, now fly away
		h.velocityY = -h.velocityY
	}
	h.X += h.velocityX
	h.Y += h.velocityY
}

func (h *Heron) IsOffscreen() bool {
	// A generous check to ensure the heron is completely gone
	return h.X < -50 || h.Y > ScreenHeight+50
}

type CrocodileState int

const (
	Floating = iota
	Biting
)

type Crocodile struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animations  map[CrocodileState]*Animation
	state       CrocodileState
}

func NewCrocodile() *Crocodile {
	spriteSheet := NewSpriteSheet(154, 42, 2, 8)
	animations := map[CrocodileState]*Animation{
		Biting:   NewAnimation([]int{1, 3, 5, 7, 9, 11, 13, 15}, 15, true),
		Floating: NewAnimation([]int{0, 2, 4, 6, 8, 10, 12, 14}, 5, false),
	}

	croc := &Crocodile{
		BaseSprite: BaseSprite{
			Location: Location{X: 0, Y: 0},
			image:    CrocodileSpriteSheet,
			srcRect:  spriteSheet.Rect(0),
			hitbox:   Rect{left: 20, top: 20, right: 70, bottom: 42},
		},
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       Floating,
	}

	return croc
}

func (c *Crocodile) Update() {
	animation := c.animations[c.state]
	animation.Update()
	c.srcRect = c.spriteSheet.Rect(animation.Frame())

	if c.state == Biting && c.animations[Biting].IsFinished() {
		c.state = Floating
	}

	c.X -= CrocodileSpeed
	if c.X < float64(-c.srcRect.Dx()) {
		c.X = ScreenWidth
	}
}

func (c *Crocodile) Bite() {
	c.state = Biting
	c.animations[Biting].Reset()
}

func (c *Crocodile) IsBiting() bool {
	return c.state == Biting
}
