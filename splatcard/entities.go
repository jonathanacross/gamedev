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
	spriteSheet := NewSpriteSheet(48, 48, 9, 5)

	frog := &Frog{
		BaseSprite: BaseSprite{
			Location: Location{X: 0, Y: 0},
			image:    FrogSpriteSheet,
			srcRect:  spriteSheet.Rect(0),
			hitbox:   Rect{left: 16, top: 16, right: 32, bottom: 32},
		},
		spriteSheet: spriteSheet,
		animations:  make(map[FrogState]*Animation),
		state:       Idle,
	}

	return frog
}

type Rock struct {
	BaseSprite
}
