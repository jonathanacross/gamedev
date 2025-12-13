package main

type Star struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation

	targeting bool
	direction Vector
	velocity  float64
	accel     float64
	finished  bool
}

func NewStar(location Location, direction Vector) *Star {
	spriteSheet := NewSpriteSheet(32, 32, 4, 3)
	var animation *Animation = nil
	var initialVelocity float64
	var targetDist float64

	animation = NewAnimation([]int{0, 1, 2, 3}, 5, true)
	//animation = NewAnimation([]int{8, 9, 10, 11}, 5, true)
	//animation = NewAnimation([]int{0, 4, 1, 5, 2, 6, 3, 7}, 5, true)
	initialVelocity = 3.0
	targetDist = 4.0 * TileSize
	accel := initialVelocity * initialVelocity / (2 * targetDist)

	return &Star{
		BaseSprite: BaseSprite{
			Location: location,
			image:    StarSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
				X: 16,
				Y: 16,
			},
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		targeting:   false,
		direction:   direction,
		velocity:    initialVelocity,
		accel:       accel,
		finished:    false,
	}
}

func (b *Star) Update(level *Level, player *Player) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())

	b.velocity -= b.accel
	v := b.direction.Normalize().Scale(b.velocity)
	b.X += v.X
	b.Y += v.Y

	if b.velocity <= 0 {
		b.finished = true
	}
	// } else {
	// 	b.velocity += b.accel
	// 	bToPlayer := Vector(player.Location()).Minus(Vector(b.Location))
	// 	b.direction = bToPlayer.Normalize()
	// 	v := b.direction.Normalize().Scale(b.velocity)
	// 	b.X += v.X
	// 	b.Y += v.Y

	// 	if bToPlayer.Length() <= 6 {
	// 		b.finished = true
	// 		action := Action{
	// 			Type: ActionReturnStar,
	// 		}
	// 		return UpdateResult{Actions: []Action{action}}
	// 	}
	// }

	// var actions []Action
	// stunDuration := 120 // 2 seconds stun
	// hitBox := Rect{
	// 	Left:   b.X - 5,
	// 	Top:    b.Y - 5,
	// 	Right:  b.X + 5,
	// 	Bottom: b.Y + 5,
	// }
	// // Create a stun source at the boomerang's location
	// ds := NewStunSource(TagPlayer, hitBox, stunDuration)
	// actions = append(actions, Action{
	// 	Type:         ActionCreateDamageSource,
	// 	DamageSource: ds,
	// })

	return UpdateResult{}
}

func (b *Star) CanRemove() bool {
	return b.finished
}
