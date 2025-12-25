package objects

import (
	"math"
	"roguerpg/assets"
	"roguerpg/core"
)

type FireState int

const (
	fireStateFlying FireState = iota
	fireStateHitting
	fireStateDone
)

// Implements GameObject
type Fire struct {
	core.BaseSprite
	state                 FireState
	spriteSheet           *core.SpriteSheet
	remainingFlightFrames int
	flyAnimation          *core.Animation
	hitAnimation          *core.Animation
	velocity              core.Vector
	damageSource          core.DamageSourceConfig
}

const (
	fireDownIndex = iota
	fireLeftIndex
	fireUpIndex
	fireRightIndex
)

func NewFire(location core.Location, velocity core.Vector) *Fire {
	spriteSheet := core.NewSpriteSheet(48, 48, 9, 4)

	directionOffsets := map[core.Direction]int{
		core.Down:  0,
		core.Up:    9,
		core.Left:  18,
		core.Right: 27,
	}
	flightAnimations := core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 8, false)
	hitAnimations := core.NewDirectionAnimationMap([]int{4, 5, 6, 7, 8}, 0, directionOffsets, 8, false)

	var spriteSheetDir int
	direction := core.VectorToDirection(velocity)
	hitBox := core.Rect{}
	switch direction {
	case core.Up:
		spriteSheetDir = fireUpIndex
		hitBox = core.Rect{Left: -9, Top: -9, Right: 9, Bottom: 9}
	case core.Down:
		spriteSheetDir = fireDownIndex
		hitBox = core.Rect{Left: -9, Top: -9, Right: 9, Bottom: 9}
	case core.Left:
		spriteSheetDir = fireLeftIndex
		hitBox = core.Rect{Left: -9, Top: -9, Right: 9, Bottom: 9}
	case core.Right:
		spriteSheetDir = fireRightIndex
		hitBox = core.Rect{Left: -9, Top: -9, Right: 9, Bottom: 9}
	}

	fireMaxDistance := float64(7 * core.TileSize)
	fireSpeed := 2.5
	fireDamage := 1
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
		state:                 fireStateFlying,
		spriteSheet:           spriteSheet,
		flyAnimation:          flightAnimations[direction],
		hitAnimation:          hitAnimations[direction],
		remainingFlightFrames: remainingFlightFrames,
		velocity:              velocity.Normalize().Scale(fireSpeed),
		damageSource:          damageSource,
	}
}

func (b *Fire) Update(level core.Level, _ core.Player) core.UpdateResult {
	b.Loc = core.Location(core.Vector(b.Loc).Plus(b.velocity))
	b.remainingFlightFrames--

	switch b.state {
	case fireStateFlying:
		b.flyAnimation.Update()
		b.SrcRect = b.spriteSheet.Rect(b.flyAnimation.Frame())
		if b.remainingFlightFrames <= 0 {
			b.state = fireStateHitting
		}
	case fireStateHitting:
		b.hitAnimation.Update()
		b.SrcRect = b.spriteSheet.Rect(b.hitAnimation.Frame())
		if b.hitAnimation.IsFinished() {
			b.state = fireStateDone
		}
	}

	if ds := b.getActiveDamageSource(); ds != nil {
		return core.UpdateResult{Actions: []core.Action{{
			Type:         core.ActionCreateDamageSource,
			Location:     b.Loc,
			DamageSource: ds}}}
	}
	return core.UpdateResult{}
}

func (b *Fire) CanRemove() bool {
	return b.state == fireStateDone
}

func (b *Fire) getActiveDamageSource() *core.DamageSource {
	worldHitbox := b.damageSource.HitBox.Offset(b.Loc.X, b.Loc.Y)
	return core.NewDamageSource(core.TagEnemy, worldHitbox, core.DamageTypePhysical, b.damageSource.Damage)
}
