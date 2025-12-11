package main

import "math"

// Implements GameObject
type Arrow struct {
	BaseSprite
	spriteSheet           *SpriteSheet
	remainingFlightFrames int
	velocity              Vector
	damageSource          DamageSourceConfig
}

const (
	arrowDownIndex = iota
	arrowLeftIndex
	arrowUpIndex
	arrowRightIndex
)

const (
	arrowMaxDistance = 5 * TileSize
	arrowSpeed       = 3.0
)

func NewArrow(location Location, velocity Vector) *Arrow {
	spriteSheet := NewSpriteSheet(16, 16, 4, 1)
	var spriteSheetDir int
	direction := VectorToDirection(velocity)
	hitBox := Rect{}
	switch direction {
	case Up:
		spriteSheetDir = arrowUpIndex
		hitBox = Rect{Left: -2, Top: -8, Right: 2, Bottom: 0}
	case Down:
		spriteSheetDir = arrowDownIndex
		hitBox = Rect{Left: -2, Top: 0, Right: 2, Bottom: 8}
	case Left:
		spriteSheetDir = arrowLeftIndex
		hitBox = Rect{Left: -8, Top: -2, Right: 0, Bottom: 2}
	case Right:
		spriteSheetDir = arrowRightIndex
		hitBox = Rect{Left: 0, Top: -2, Right: 8, Bottom: 2}
	}
	remainingFlightFrames := int(math.Round(arrowMaxDistance / arrowSpeed))

	damageSource := DamageSourceConfig{
		HitBox: hitBox,
		Damage: 1,
	}

	return &Arrow{
		BaseSprite: BaseSprite{
			Location: location,
			image:    ArrowSpritesImage,
			srcRect:  spriteSheet.Rect(spriteSheetDir),
			drawOffset: Location{
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

func (b *Arrow) Update(level *Level, _ *Player) UpdateResult {
	b.Location = Location(Vector(b.Location).Plus(b.velocity))
	b.remainingFlightFrames--

	if ds := b.getActiveDamageSource(); ds != nil {
		return UpdateResult{Actions: []Action{{
			Type:         ActionCreateDamageSource,
			Location:     b.Location,
			DamageSource: ds}}}
	}
	return UpdateResult{}
}

func (b *Arrow) CanRemove() bool {
	return b.remainingFlightFrames <= 0
}

func (b *Arrow) getActiveDamageSource() *DamageSource {
	worldHitbox := b.damageSource.HitBox.Offset(b.X, b.Y)
	return NewDamageSource(TagPlayer, worldHitbox, b.damageSource.Damage)
}
