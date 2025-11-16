package main

type BlobEnemy struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation
	Vx          float64
	Vy          float64
}

func NewBlobEnemy() *BlobEnemy {
	animation := NewAnimation([]int{0, 1, 2}, 20, true)
	spriteSheet := NewSpriteSheet(16, 16, 3, 1)
	hitbox := Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	return &BlobEnemy{
		BaseSprite: BaseSprite{
			Location: Location{
				X: 0,
				Y: 0,
			},
			drawOffset: Location{
				X: 8,
				Y: 8,
			},
			srcRect:    spriteSheet.Rect(0),
			image:      BlobSpritesImage,
			hitbox:     hitbox,
			debugImage: createDebugRectImage(hitbox),
		},
		spriteSheet: spriteSheet,
		animation:   animation,
	}
}

func (c *BlobEnemy) Update(level *Level) {
	c.animation.Update()
	c.srcRect = c.spriteSheet.Rect(c.animation.Frame())

	// c.X += c.Vx
	// c.HandleTileCollisions(level, AxisX)
	// c.Y += c.Vy
	// c.HandleTileCollisions(level, AxisY)
}
