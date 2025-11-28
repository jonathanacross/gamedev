package main

type Boomerang struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation

	leaving   bool
	direction Vector
	velocity  float64
	accel     float64
	level     int
	finished  bool
}

func NewBoomerang(location Location, direction Vector, level int) *Boomerang {
	spriteSheet := NewSpriteSheet(16, 16, 8, 3)
	var animation *Animation = nil
	var initialVelocity float64
	var targetDist float64

	switch level {
	case 1:
		animation = NewAnimation([]int{0, 1, 2, 3, 4, 5, 6, 7}, 10, true)
		initialVelocity = 2.0
		targetDist = 4.0 * TileSize
	case 2:
		animation = NewAnimation([]int{8, 9, 10, 11, 12, 13, 14, 15}, 10, true)
		initialVelocity = 3.0
		targetDist = 6.0 * TileSize
	case 3:
		animation = NewAnimation([]int{16, 17, 18, 19, 20, 21, 22, 23}, 10, true)
		initialVelocity = 4.0
		targetDist = 8.0 * TileSize
	}
	accel := initialVelocity * initialVelocity / (2 * targetDist)

	return &Boomerang{
		BaseSprite: BaseSprite{
			Location: location,
			image:    BoomerangSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		leaving:     true,
		direction:   direction,
		velocity:    initialVelocity,
		accel:       accel,
		level:       level,
		finished:    false,
	}
}

func (b *Boomerang) Update(level *Level, player *Player) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.leaving {
		b.velocity -= b.accel
		v := b.direction.Normalize().Scale(b.velocity)
		b.X += v.X
		b.Y += v.Y

		if b.velocity <= 0 {
			b.leaving = false
		}
	} else {
		b.velocity += b.accel
		bToPlayer := Vector(player.Location()).Minus(Vector(b.Location))
		b.direction = bToPlayer.Normalize()
		v := b.direction.Normalize().Scale(b.velocity)
		b.X += v.X
		b.Y += v.Y

		if bToPlayer.Length() <= 6 {
			b.finished = true
			action := Action{
				Type:         ActionReturnBoomerang,
				Location:     Location(ZeroVector()),
				Direction:    ZeroVector(),
				DamageSource: nil,
			}
			return UpdateResult{Actions: []Action{action}}
		}
	}

	return UpdateResult{}
}

func (b *Boomerang) CanRemove() bool {
	return b.finished
}
