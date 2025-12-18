package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
	"time"
)

type Star struct {
	core.BaseSprite
	spriteSheet *core.SpriteSheet
	animation   *core.Animation

	target       core.Character
	direction    core.Vector
	speed        float64
	finished     bool
	idleTimer    *core.Timer
	damageSource core.DamageSourceConfig
}

func NewStar(location core.Location, direction core.Vector, target core.Character) *Star {
	spriteSheet := core.NewSpriteSheet(16, 16, 4, 3)
	var animation *core.Animation = nil

	//animation = NewAnimation([]int{0, 1, 2, 3}, 5, true)
	animation = core.NewAnimation([]int{0, 4, 1, 5, 2, 6, 3, 7}, 5, true)

	hitBox := core.Rect{Left: -7, Top: -7, Right: 7, Bottom: 7}
	damageSource := core.DamageSourceConfig{
		HitBox: hitBox,
		Damage: 1,
	}

	return &Star{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.StarSpritesImage,
			SrcRect: spriteSheet.Rect(0),
			DrawOffset: core.Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet:  spriteSheet,
		animation:    animation,
		target:       target, // Checks nil in original code
		direction:    direction.Normalize(),
		speed:        1.0,
		finished:     false,
		idleTimer:    core.NewTimer(1500 * time.Millisecond),
		damageSource: damageSource,
	}
}

func (b *Star) Update(level core.Level, player core.Player) core.UpdateResult {
	b.animation.Update()
	b.idleTimer.Update()
	b.SrcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.target != nil && !b.target.IsDead() {
		// If we have a valid target, track it
		turnRate := 0.05
		dCurr := b.direction
		dTarget := core.Vector(b.target.Location()).Minus(core.Vector(b.Loc)).Normalize()
		dNew := (dCurr.Scale(1 - turnRate)).Plus(dTarget.Scale(turnRate))
		b.direction = dNew.Normalize()
		v := b.direction.Scale(b.speed)
		b.Loc.X += v.X
		b.Loc.Y += v.Y

		// Check collision
		targetDir := core.Vector(b.target.Location()).Minus(core.Vector(b.Loc))
		if targetDir.Length() <= 10 {
			b.finished = true
		}
	} else {
		// No target or target dead
		v := b.direction.Scale(b.speed)
		b.Loc.X += v.X
		b.Loc.Y += v.Y
	}
	if b.idleTimer.IsReady() {
		b.finished = true
	}

	if ds := b.getActiveDamageSource(); ds != nil {
		return core.UpdateResult{Actions: []core.Action{{
			Type:         core.ActionCreateDamageSource,
			Location:     b.Loc,
			DamageSource: ds}}}
	}
	return core.UpdateResult{}
}

func (b *Star) CanRemove() bool {
	return b.finished
}

func (b *Star) getActiveDamageSource() *core.DamageSource {
	worldHitbox := b.damageSource.HitBox.Offset(b.Loc.X, b.Loc.Y)
	// TODO: make this magical damage source and handle in monster hits
	ds := core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypePhysical, b.damageSource.Damage)
	ds.OnHit = func() {
		b.finished = true
	}
	return ds
}
