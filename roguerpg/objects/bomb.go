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
	exploded    bool
}

func NewBomb(location core.Location) *Bomb {
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

func NewBombExplosion(location core.Location) *BombExplosion {
	animation := core.NewAnimation([]int{0, 1, 2, 3, 4, 5}, 5, false)
	spriteSheet := core.NewSpriteSheet(48, 48, 6, 1)

	attackHitboxes := map[int]core.DamageSourceConfig{
		1: {HitBox: core.Rect{Left: -24, Top: -24, Right: 24, Bottom: 24}, Damage: 2},
	}

	return &BombExplosion{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.BombExplosionSpritesImage,
			SrcRect: spriteSheet.Rect(0),
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

func (b *BombExplosion) Update(level core.Level, _ core.Player) core.UpdateResult {
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
	// Wrong. Frame() returns the FRAME NUMBER (from array).
	// Original used `b.animation.frameIndex` which is the index in array.
	// Frame() returns frames[frameIndex].
	// If frames are [0, 1, 2, 3, 4, 5], Frame() = frameIndex.
	// But safely, we want index?
	// Original `bomb.go`: `animIndex := b.animation.frameIndex`.
	// But `frameIndex` is not exported in `core.Animation`.
	// I need to export `FrameIndex` or add a getter.
	// Or change `Frame()` to return what?

	// In `core/animation.go`: `func (a *Animation) Frame() int { return a.frames[a.frameIndex] }`
	// So Frame() returns sprite index.

	// BombExplosion `attackHitboxes` map key is 1.
	// Animation frames: 0, 1, 2, 3, 4, 5.
	// So if Frame() returns 1, we match.
	// So it works if frame values are unique/sequential.
	// Here they are [0,1,2,3...].
	// So `b.animation.Frame()` works correctly.

	if config, ok := b.attackHitboxes[animIndex]; ok {
		worldHitbox := config.HitBox.Offset(b.Loc.X, b.Loc.Y)
		return core.NewDamageSource(core.TagPlayer, worldHitbox, config.Damage)
	}

	return nil
}
