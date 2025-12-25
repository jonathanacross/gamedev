package objects

import (
	"math"
	"roguerpg/assets"
	"roguerpg/core"
)

// Implements GameObject
type Fire struct {
	core.BaseSprite
	spriteSheet           *core.SpriteSheet
	remainingFlightFrames int
	velocity              core.Vector
	damageSource          core.DamageSourceConfig
}

const (
	fireDownIndex = iota
	fireLeftIndex
	fireUpIndex
	fireRightIndex
)

func NewFire(location core.Location, velocity core.Vector, level int) *Fire {
	spriteSheet := core.NewSpriteSheet(48, 48, 9, 4)
	var spriteSheetDir int
	direction := core.VectorToDirection(velocity)
	hitBox := core.Rect{}
	switch direction {
	case core.Up:
		spriteSheetDir = fireUpIndex
		hitBox = core.Rect{Left: -10, Top: -10, Right: 10, Bottom: 10}
	case core.Down:
		spriteSheetDir = fireDownIndex
		hitBox = core.Rect{Left: -10, Top: -10, Right: 10, Bottom: 10}
	case core.Left:
		spriteSheetDir = fireLeftIndex
		hitBox = core.Rect{Left: -10, Top: -10, Right: 10, Bottom: 10}
	case core.Right:
		spriteSheetDir = fireRightIndex
		hitBox = core.Rect{Left: -10, Top: -10, Right: 10, Bottom: 10}
	}

	level = core.Clamp(level, 1, 3)
	fireMaxDistance := float64((5 + 2*level) * core.TileSize)
	fireSpeed := 3.0 + 0.5*float64(level)
	fireDamage := 2 * level
	remainingFlightFrames := int(math.Round(fireMaxDistance / fireSpeed))

	damageSource := core.DamageSourceConfig{
		HitBox: hitBox,
		Damage: fireDamage,
	}

	return &Fire{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.FireSpritesImage,
			SrcRect: spriteSheet.Rect(spriteSheetDir),
			DrawOffset: core.Location{
				X: 24,
				Y: 24,
			},
		},
		spriteSheet:           spriteSheet,
		remainingFlightFrames: remainingFlightFrames,
		velocity:              velocity.Normalize().Scale(fireSpeed),
		damageSource:          damageSource,
	}
}

func (b *Fire) Update(level core.Level, _ core.Player) core.UpdateResult {
	b.Loc = core.Location(core.Vector(b.Loc).Plus(b.velocity))
	b.remainingFlightFrames--

	if ds := b.getActiveDamageSource(); ds != nil {
		return core.UpdateResult{Actions: []core.Action{{
			Type:         core.ActionCreateDamageSource,
			Location:     b.Loc,
			DamageSource: ds}}}
	}
	return core.UpdateResult{}
}

func (b *Fire) CanRemove() bool {
	return b.remainingFlightFrames <= 0
}

func (b *Fire) getActiveDamageSource() *core.DamageSource {
	worldHitbox := b.damageSource.HitBox.Offset(b.Loc.X, b.Loc.Y)
	return core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypePhysical, b.damageSource.Damage)
}
