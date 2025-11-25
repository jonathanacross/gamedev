package main

// Implements GameObject
type Bomb struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation
	exploded    bool
}

func NewBomb(location Location) *Bomb {
	animation := NewAnimation([]int{
		0, 1, 0, 1, 2, 3, 2, 3,
		0, 1, 0, 1, 2, 3, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 3, 0, 3,
	}, 8, false)
	spriteSheet := NewSpriteSheet(16, 16, 4, 1)

	return &Bomb{
		BaseSprite: BaseSprite{
			Location: location,
			image:    BombSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		exploded:    false,
	}
}

func (b *Bomb) Update(level *Level) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.animation.IsFinished() {
		b.exploded = true
		return UpdateResult{Actions: []Action{{Type: ActionExplosion, Location: Location{X: b.X, Y: b.Y}}}}
	}
	return UpdateResult{}
}

func (b *Bomb) CanRemove() bool {
	return b.exploded
}

type BombExplosion struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation
	finished    bool
}

func NewBombExplosion(location Location) *BombExplosion {
	animation := NewAnimation([]int{0, 1, 2, 3, 4, 5}, 5, false)
	spriteSheet := NewSpriteSheet(32, 32, 6, 1)

	return &BombExplosion{
		BaseSprite: BaseSprite{
			Location: location,
			image:    BombExplosionSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
				X: 16,
				Y: 16,
			},
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		finished:    false,
	}
}

func (b *BombExplosion) Update(level *Level) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())
	if b.animation.IsFinished() {
		b.finished = true
	}
	return UpdateResult{}
}

func (b *BombExplosion) CanRemove() bool {
	return b.finished
}
