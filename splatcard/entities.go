package main

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
}

func NewFrog() *Frog {
	spriteSheet := NewSpriteSheet(32, 32, 5, 4)

	frog := &Frog{
		BaseSprite: BaseSprite{
			Location: Location{X: 0, Y: 0},
			image:    FrogSpriteSheet,
			srcRect:  spriteSheet.Rect(0),
			hitbox:   Rect{left: 12, top: 12, right: 24, bottom: 24},
		},
		spriteSheet: spriteSheet,
		animations:  make(map[FrogState]*Animation),
		state:       Idle,
	}

	return frog
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
