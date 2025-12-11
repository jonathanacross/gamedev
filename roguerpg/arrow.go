package main

import "math"

// Implements GameObject
type Arrow struct {
	BaseSprite
	spriteSheet           *SpriteSheet
	remainingFlightFrames int
	velocity              Vector
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
	switch direction {
	case Up:
		spriteSheetDir = arrowUpIndex
	case Down:
		spriteSheetDir = arrowDownIndex
	case Left:
		spriteSheetDir = arrowLeftIndex
	case Right:
		spriteSheetDir = arrowRightIndex
	}
	remainingFlightFrames := int(math.Round(arrowMaxDistance / arrowSpeed))

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
	}
}

func (b *Arrow) Update(level *Level, _ *Player) UpdateResult {
	b.Location = Location(Vector(b.Location).Plus(b.velocity))
	b.remainingFlightFrames--

	return UpdateResult{}
}

func (b *Arrow) CanRemove() bool {
	return b.remainingFlightFrames <= 0
}

type ArrowExplosion struct {
	BaseSprite
	spriteSheet    *SpriteSheet
	animation      *Animation
	finished       bool
	attackHitboxes map[int]DamageSourceConfig
}

func NewArrowExplosion(location Location) *ArrowExplosion {
	animation := NewAnimation([]int{0, 1, 2, 3, 4, 5}, 5, false)
	spriteSheet := NewSpriteSheet(48, 48, 6, 1)

	attackHitboxes := map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -24, Top: -24, Right: 24, Bottom: 24}, Damage: 2},
	}

	return &ArrowExplosion{
		BaseSprite: BaseSprite{
			Location: location,
			image:    ArrowSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
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

func (b *ArrowExplosion) Update(level *Level, _ *Player) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())
	if b.animation.IsFinished() {
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

func (b *ArrowExplosion) CanRemove() bool {
	return b.finished
}

func (b *ArrowExplosion) getActiveDamageSource() *DamageSource {
	animIndex := b.animation.frameIndex

	if config, ok := b.attackHitboxes[animIndex]; ok {
		worldHitbox := config.HitBox.Offset(b.X, b.Y)
		return NewDamageSource(TagPlayer, worldHitbox, config.Damage)
	}

	return nil
}
