package main

import "time"

type Star struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation

	target       Character
	direction    Vector
	speed        float64
	finished     bool
	idleTimer    *Timer
	damageSource DamageSourceConfig
}

func NewStar(location Location, direction Vector, target Character) *Star {
	spriteSheet := NewSpriteSheet(16, 16, 4, 3)
	var animation *Animation = nil

	animation = NewAnimation([]int{0, 1, 2, 3}, 5, true)
	//animation = NewAnimation([]int{8, 9, 10, 11}, 5, true)
	//animation = NewAnimation([]int{0, 4, 1, 5, 2, 6, 3, 7}, 5, true)

	hitBox := Rect{Left: -7, Top: -7, Right: 7, Bottom: 7}
	damageSource := DamageSourceConfig{
		HitBox: hitBox,
		Damage: 1,
	}

	return &Star{
		BaseSprite: BaseSprite{
			Location: location,
			image:    StarSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet:  spriteSheet,
		animation:    animation,
		target:       target,
		direction:    direction.Normalize(),
		speed:        1.0,
		finished:     false,
		idleTimer:    NewTimer(1500 * time.Millisecond),
		damageSource: damageSource,
	}
}

func (b *Star) Update(level *Level, player *Player) UpdateResult {
	b.animation.Update()
	b.idleTimer.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.target != nil && !b.target.IsDead() {
		// If we have a valid target, track it
		turnRate := 0.05
		dCurr := b.direction
		dTarget := Vector(b.target.Location()).Minus(Vector(b.Location)).Normalize()
		dNew := (dCurr.Scale(1 - turnRate)).Plus(dTarget.Scale(turnRate))
		b.direction = dNew.Normalize()
		v := b.direction.Scale(b.speed)
		b.X += v.X
		b.Y += v.Y

		// Check collision
		targetDir := Vector(b.target.Location()).Minus(Vector(b.Location))
		if targetDir.Length() <= 10 {
			b.finished = true
		}
	} else {
		// No target or target dead
		v := b.direction.Scale(b.speed)
		b.X += v.X
		b.Y += v.Y
	}
	if b.idleTimer.IsReady() {
		b.finished = true
	}

	if ds := b.getActiveDamageSource(); ds != nil {
		return UpdateResult{Actions: []Action{{
			Type:         ActionCreateDamageSource,
			Location:     b.Location,
			DamageSource: ds}}}
	}
	return UpdateResult{}
}

func (b *Star) CanRemove() bool {
	return b.finished
}

func (b *Star) getActiveDamageSource() *DamageSource {
	worldHitbox := b.damageSource.HitBox.Offset(b.X, b.Y)
	return NewDamageSource(TagPlayer, worldHitbox, b.damageSource.Damage)
}
