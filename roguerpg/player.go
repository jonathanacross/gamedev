package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walking
	Dying
)

type PlayerDirection int

const (
	Left PlayerDirection = iota
	Right
	Up
	Down
)

type Player struct {
	BaseSprite
	images      map[PlayerState]*ebiten.Image
	spriteSheet *SpriteSheet
	animations  map[PlayerState]map[PlayerDirection]*Animation
	state       PlayerState
	direction   PlayerDirection
}

func NewPlayer() *Player {
	animations := map[PlayerState]map[PlayerDirection]*Animation{
		Idle: {
			Left:  NewAnimation([]int{6, 7, 8, 9, 10, 11}, 15, true),
			Right: NewAnimation([]int{30, 31, 32, 33, 34, 35}, 15, true),
			Up:    NewAnimation([]int{18, 19, 20, 21, 22, 23}, 15, true),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5}, 15, true),
		},
		Walking: {
			Left:  NewAnimation([]int{6, 7, 8, 9, 10, 11}, 15, true),
			Right: NewAnimation([]int{30, 31, 32, 33, 34, 35}, 15, true),
			Up:    NewAnimation([]int{18, 19, 20, 21, 22, 23}, 15, true),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5}, 15, true),
		},
		Dying: {
			Left:  NewAnimation([]int{6, 7, 8, 9, 10, 11}, 15, false),
			Right: NewAnimation([]int{30, 31, 32, 33, 34, 35}, 15, false),
			Up:    NewAnimation([]int{18, 19, 20, 21, 22, 23}, 15, false),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5}, 15, false),
		},
	}

	charImages := map[PlayerState]*ebiten.Image{
		Idle:    PlayerIdleSpritesImage,
		Walking: PlayerWalkSpritesImage,
		Dying:   PlayerDeathSpritesImage,
	}

	spriteSheet := NewSpriteSheet(48, 48, 6, 6)

	return &Player{
		BaseSprite: BaseSprite{
			Location: Location{
				X: 100,
				Y: 50,
			},
			srcRect: spriteSheet.Rect(0),
			hitbox: Rect{
				Left:   0,
				Top:    0,
				Right:  8,
				Bottom: 16,
			},
		},
		images:      charImages,
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       Idle,
		direction:   Down,
	}
}

func (c *Player) Update(state PlayerState) {
	c.state = state

	animationSet, exists := c.animations[state]
	if !exists {
		return
	}
	animation, exists := animationSet[c.direction]
	if !exists {
		return
	}

	animation.Update()

	c.image = c.images[state]
	c.srcRect = c.spriteSheet.Rect(animation.Frame())
}
