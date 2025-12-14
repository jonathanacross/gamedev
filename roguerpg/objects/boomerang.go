package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
)

type Boomerang struct {
	core.BaseSprite
	spriteSheet *core.SpriteSheet
	animation   *core.Animation

	leaving   bool
	direction core.Vector
	velocity  float64
	accel     float64
	level     int
	finished  bool
}

func NewBoomerang(location core.Location, direction core.Vector, level int) *Boomerang {
	spriteSheet := core.NewSpriteSheet(16, 16, 8, 3)
	var animation *core.Animation = nil
	var initialVelocity float64
	var targetDist float64

	switch level {
	case 1:
		animation = core.NewAnimation([]int{0, 1, 2, 3, 4, 5, 6, 7}, 5, true)
		initialVelocity = 2.0
		targetDist = 4.0 * core.TileSize
	case 2:
		animation = core.NewAnimation([]int{8, 9, 10, 11, 12, 13, 14, 15}, 4, true)
		initialVelocity = 3.0
		targetDist = 6.0 * core.TileSize
	case 3:
		animation = core.NewAnimation([]int{16, 17, 18, 19, 20, 21, 22, 23}, 3, true)
		initialVelocity = 4.0
		targetDist = 8.0 * core.TileSize
	}
	accel := initialVelocity * initialVelocity / (2 * targetDist)

	return &Boomerang{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.BoomerangSpritesImage,
			SrcRect: spriteSheet.Rect(0),
			DrawOffset: core.Location{
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

func (b *Boomerang) Update(level core.Level, player core.Player) core.UpdateResult {
	b.animation.Update()
	b.SrcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.leaving {
		b.velocity -= b.accel
		v := b.direction.Normalize().Scale(b.velocity)
		b.Loc.X += v.X
		b.Loc.Y += v.Y

		if b.velocity <= 0 {
			b.leaving = false
		}
	} else {
		b.velocity += b.accel
		bToPlayer := core.Vector(player.Location()).Minus(core.Vector(b.Loc))
		b.direction = bToPlayer.Normalize()
		v := b.direction.Normalize().Scale(b.velocity)
		b.Loc.X += v.X
		b.Loc.Y += v.Y

		if bToPlayer.Length() <= 6 {
			b.finished = true
			action := core.Action{
				Type: core.ActionReturnBoomerang,
			}
			return core.UpdateResult{Actions: []core.Action{action}}
		}
	}

	var actions []core.Action
	stunDuration := 120 // 2 seconds stun
	hitBox := core.Rect{
		Left:   b.Loc.X - 5,
		Top:    b.Loc.Y - 5,
		Right:  b.Loc.X + 5,
		Bottom: b.Loc.Y + 5,
	}
	// Create a stun source at the boomerang's location
	ds := core.NewStunSource(core.TagPlayer, hitBox, stunDuration)
	ds.OnHit = func() {
		b.leaving = false
	}
	actions = append(actions, core.Action{
		Type:         core.ActionCreateDamageSource,
		DamageSource: ds,
	})

	return core.UpdateResult{Actions: actions}
}

func (b *Boomerang) CanRemove() bool {
	return b.finished
}
