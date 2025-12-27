package objects

import (
	"math"
	"roguerpg/assets"
	"roguerpg/core"
)

const (
	boulderSpeed       = 2.5
	boulderDamage      = 2
	boulderMaxDistance = 7 * core.TileSize
)

// Implements GameObject
type Boulder struct {
	core.BaseSprite
	spriteSheet           *core.SpriteSheet
	animation             *core.Animation
	remainingFlightFrames int
	velocity              core.Vector
	damageSource          core.DamageSourceConfig
}

func NewBoulder(location core.Location, velocity core.Vector) *Boulder {
	spriteSheet := core.NewSpriteSheet(16, 16, 4, 1)
	hitBox := core.Rect{Left: -6, Top: -6, Right: 6, Bottom: 6}
	animation := core.NewAnimation([]int{0, 1, 2, 3}, 5, true)

	remainingFlightFrames := int(math.Round(float64(boulderMaxDistance) / boulderSpeed))

	damageSource := core.DamageSourceConfig{
		HitBox: hitBox,
		Damage: boulderDamage,
	}

	return &Boulder{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.BoulderSpritesImage,
			SrcRect: spriteSheet.Rect(0),
			DrawOffset: core.Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet:           spriteSheet,
		animation:             animation,
		remainingFlightFrames: remainingFlightFrames,
		velocity:              velocity.Normalize().Scale(boulderSpeed),
		damageSource:          damageSource,
	}
}

func (b *Boulder) Update(level core.Level, _ core.Player) core.UpdateResult {
	b.animation.Update()
	b.SrcRect = b.spriteSheet.Rect(b.animation.Frame())
	b.Loc = core.Location(core.Vector(b.Loc).Plus(b.velocity))
	b.remainingFlightFrames--

	// Create a damage source at the boulder's location
	worldHitbox := b.damageSource.HitBox.Offset(b.Loc.X, b.Loc.Y)
	ds := core.NewDamageSource(core.TagEnemy, worldHitbox, core.DamageTypePhysical, b.damageSource.Damage)
	actions := []core.Action{{
		Type:         core.ActionCreateDamageSource,
		DamageSource: ds,
	}}

	return core.UpdateResult{Actions: actions}
}

func (b *Boulder) CanRemove() bool {
	return b.remainingFlightFrames <= 0
}
