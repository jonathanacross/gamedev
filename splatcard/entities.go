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
}

func NewFrog() *Frog {
	spriteSheet := NewSpriteSheet(32, 32, 5, 4)
	animations := map[FrogState]*Animation{
		// Idle animation now loops
		Idle: NewAnimation(0, 1, 15, true),
		// Jumping and Surprised animations do not loop
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

func (f *Frog) Update(g *Game) {
	animation := f.animations[f.state]
	animation.Update()
	f.srcRect = f.spriteSheet.Rect(animation.Frame())

	if f.state == Jumping {
		jumpAnimation := f.animations[Jumping]

		// This logic is now in `main.go`, but the position update stays here
		if jumpAnimation.IsFinished() {
			// This part is handled by main.go, but we must return here
			// to avoid re-calculating position on the final frame
			return
		}

		totalAnimationFrames := float64(jumpAnimation.last-jumpAnimation.first+1) * float64(jumpAnimation.speed)
		currentFrameCount := float64(jumpAnimation.frame-jumpAnimation.first)*float64(jumpAnimation.speed) + float64(jumpAnimation.speed-jumpAnimation.frameCounter)
		progress := currentFrameCount / totalAnimationFrames

		jumpDistance := g.jumpTargetX - f.jumpStartX
		f.X = f.jumpStartX + jumpDistance*progress
		f.Y = float64(PlatformY-32) - 15*math.Sin(math.Pi*progress)
	} else {
		f.Y = float64(PlatformY - 32)
	}
}

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

type Rock struct {
	BaseSprite
}
