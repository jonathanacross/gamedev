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
		Idle:      NewAnimation(0, 1, 15, true),
		Jumping:   NewAnimation(10, 14, 10, false),
		Surprised: NewAnimation(5, 6, 10, false),
		Dying:     NewAnimation(15, 17, 10, false),
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
		totalAnimationFrames := float64(jumpAnimation.last-jumpAnimation.first+1) * float64(jumpAnimation.speed)
		currentFrameCount := float64(jumpAnimation.frame-jumpAnimation.first)*float64(jumpAnimation.speed) + float64(jumpAnimation.speed-jumpAnimation.frameCounter)
		progress := currentFrameCount / totalAnimationFrames

		jumpDistance := f.jumpTargetX - f.jumpStartX
		f.X = f.jumpStartX + jumpDistance*progress
		f.Y = float64(PlatformY-32) - 15*math.Sin(math.Pi*progress)
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
	f.Y = float64(PlatformY - 32)
	f.X = f.jumpTargetX
}

// Hit changes the frog's state to Dying.
func (f *Frog) Hit() {
	f.state = Dying
	f.animations[Dying].Reset()
}

// Rock and Platform structs are unchanged
type Platform struct {
	BaseSprite
}

func NewPlatform(x, y float64) *Platform {
	PlatformSprite.Bounds()
	return &Platform{
		BaseSprite: BaseSprite{
			Location: Location{X: x, Y: y},
			image:    PlatformSprite,
			srcRect:  PlatformSprite.Bounds(),
			hitbox:   NewRect(PlatformSprite.Bounds()),
		},
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
