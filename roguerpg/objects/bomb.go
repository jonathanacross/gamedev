package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
)

// Implements GameObject
type Bomb struct {
	core.BaseSprite
	spriteSheet *core.SpriteSheet
	animation   *core.Animation
	level       int
	exploded    bool
}

func NewBomb(location core.Location, level int) *Bomb {
	animation := core.NewAnimation([]int{
		0, 1, 0, 1, 2, 3, 2, 3,
		0, 1, 0, 1, 2, 3, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 3, 0, 3,
	}, 8, false)
	spriteSheet := core.NewSpriteSheet(16, 16, 4, 1)

	return &Bomb{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.BombSpritesImage,
			SrcRect: spriteSheet.Rect(0),
			DrawOffset: core.Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		level:       level,
		exploded:    false,
	}
}

func (b *Bomb) Update(level core.Level, _ core.Player) core.UpdateResult {
	b.animation.Update()
	b.SrcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.animation.IsFinished() {
		b.exploded = true
		return core.UpdateResult{Actions: []core.Action{{
			Type:     core.ActionExplosion,
			Location: b.Loc,
			Level:    b.level,
		}}}
	}
	return core.UpdateResult{}
}

func (b *Bomb) CanRemove() bool {
	return b.exploded
}

type BombExplosion struct {
	core.BaseSprite
	spriteSheet    *core.SpriteSheet
	animation      *core.Animation
	finished       bool
	attackHitboxes map[int]core.DamageSourceConfig
}

func NewBombExplosion(location core.Location, level int) *BombExplosion {
	spriteSheet := core.NewSpriteSheet(48, 48, 6, 3)

	var explosionDamage int
	var animation *core.Animation
	var attackHitboxes map[int]core.DamageSourceConfig
	level = core.Clamp(level, 1, 3)
	switch level {
	case 1:
		explosionDamage = 2
		animation = core.NewAnimation([]int{0, 1, 2, 3, 4, 5}, 5, false)
		attackHitboxes = map[int]core.DamageSourceConfig{
			1: {HitBox: core.Rect{Left: -16, Top: -16, Right: 16, Bottom: 16}, Damage: explosionDamage},
		}
	case 2:
		explosionDamage = 4
		animation = core.NewAnimation([]int{6, 7, 8, 9, 10, 11}, 5, false)
		attackHitboxes = map[int]core.DamageSourceConfig{
			7: {HitBox: core.Rect{Left: -20, Top: -20, Right: 20, Bottom: 20}, Damage: explosionDamage},
		}
	case 3:
		explosionDamage = 8
		animation = core.NewAnimation([]int{12, 13, 14, 15, 16, 17}, 5, false)
		attackHitboxes = map[int]core.DamageSourceConfig{
			13: {HitBox: core.Rect{Left: -24, Top: -24, Right: 24, Bottom: 24}, Damage: explosionDamage},
		}
	}

	return &BombExplosion{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.BombExplosionSpritesImage,
			SrcRect: spriteSheet.Rect(animation.Frame()),
			DrawOffset: core.Location{
				X: 24,
				Y: 24,
			},
		},
		spriteSheet:    spriteSheet,
		animation:      animation,
		finished:       false,
		attackHitboxes: attackHitboxes,
	}
}

func (b *BombExplosion) Update(_ core.Level, _ core.Player) core.UpdateResult {
	b.animation.Update()
	b.SrcRect = b.spriteSheet.Rect(b.animation.Frame())
	if b.animation.IsFinished() {
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

func (b *BombExplosion) CanRemove() bool {
	return b.finished
}

func (b *BombExplosion) getActiveDamageSource() *core.DamageSource {
	animIndex := b.animation.Frame()

	if config, ok := b.attackHitboxes[animIndex]; ok {
		worldHitbox := config.HitBox.Offset(b.Loc.X, b.Loc.Y)
		return core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypeImpact, config.Damage)
	}

	return nil
}
