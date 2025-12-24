package objects

import (
	"math"
	"roguerpg/assets"
	"roguerpg/core"
)

// Implements GameObject
type Arrow struct {
	core.BaseSprite
	spriteSheet           *core.SpriteSheet
	remainingFlightFrames int
	velocity              core.Vector
	damageSource          core.DamageSourceConfig
}

const (
	arrowDownIndex = iota
	arrowLeftIndex
	arrowUpIndex
	arrowRightIndex
)

func NewArrow(location core.Location, velocity core.Vector, level int) *Arrow {
	spriteSheet := core.NewSpriteSheet(16, 16, 4, 1)
	var spriteSheetDir int
	direction := core.VectorToDirection(velocity)
	hitBox := core.Rect{}
	switch direction {
	case core.Up:
		spriteSheetDir = arrowUpIndex
		hitBox = core.Rect{Left: -2, Top: -8, Right: 2, Bottom: 0}
	case core.Down:
		spriteSheetDir = arrowDownIndex
		hitBox = core.Rect{Left: -2, Top: 0, Right: 2, Bottom: 8}
	case core.Left:
		spriteSheetDir = arrowLeftIndex
		hitBox = core.Rect{Left: -8, Top: -2, Right: 0, Bottom: 2}
	case core.Right:
		spriteSheetDir = arrowRightIndex
		hitBox = core.Rect{Left: 0, Top: -2, Right: 8, Bottom: 2}
	}

	level = core.Clamp(level, 1, 3)
	arrowMaxDistance := float64((5 + 2*level) * core.TileSize)
	arrowSpeed := 3.0 + 0.5*float64(level)
	arrowDamage := 2 * level
	remainingFlightFrames := int(math.Round(arrowMaxDistance / arrowSpeed))

	damageSource := core.DamageSourceConfig{
		HitBox: hitBox,
		Damage: arrowDamage,
	}

	return &Arrow{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.ArrowSpritesImage,
			SrcRect: spriteSheet.Rect(spriteSheetDir),
			DrawOffset: core.Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet:           spriteSheet,
		remainingFlightFrames: remainingFlightFrames,
		velocity:              velocity.Normalize().Scale(arrowSpeed),
		damageSource:          damageSource,
	}
}

func (b *Arrow) Update(level core.Level, _ core.Player) core.UpdateResult {
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

func (b *Arrow) CanRemove() bool {
	return b.remainingFlightFrames <= 0
}

func (b *Arrow) getActiveDamageSource() *core.DamageSource {
	worldHitbox := b.damageSource.HitBox.Offset(b.Loc.X, b.Loc.Y)
	return core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypePhysical, b.damageSource.Damage)
}
