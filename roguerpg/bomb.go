package main

// Implements GameObject
type Bomb struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation
	exploded    bool
}

func NewBomb(location Location) *Bomb {
	animation := NewAnimation([]int{
		0, 1, 0, 1, 2, 3, 2, 3,
		0, 1, 0, 1, 2, 3, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 3, 0, 3,
	}, 8, false)
	spriteSheet := NewSpriteSheet(16, 16, 4, 1)

	return &Bomb{
		BaseSprite: BaseSprite{
			Location: location,
			image:    BombSpritesImage,
			srcRect:  spriteSheet.Rect(0),
			drawOffset: Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		exploded:    false,
	}
}

func (b *Bomb) Update(level *Level, _ *Player) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())

	if b.animation.IsFinished() {
		b.exploded = true
		return UpdateResult{Actions: []Action{{
			Type:         ActionExplosion,
			Location:     b.Location,
			Direction:    ZeroVector(),
			DamageSource: nil,
		}}}
	}
	return UpdateResult{}
}

func (b *Bomb) CanRemove() bool {
	return b.exploded
}

type BombExplosion struct {
	BaseSprite
	spriteSheet    *SpriteSheet
	animation      *Animation
	finished       bool
	attackHitboxes map[int]DamageSourceConfig
}

func NewBombExplosion(location Location) *BombExplosion {
	animation := NewAnimation([]int{0, 1, 2, 3, 4, 5}, 5, false)
	spriteSheet := NewSpriteSheet(48, 48, 6, 1)

	attackHitboxes := map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -24, Top: -24, Right: 24, Bottom: 24}, Damage: 2},
	}

	return &BombExplosion{
		BaseSprite: BaseSprite{
			Location: location,
			image:    BombExplosionSpritesImage,
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

func (b *BombExplosion) Update(level *Level, _ *Player) UpdateResult {
	b.animation.Update()
	b.srcRect = b.spriteSheet.Rect(b.animation.Frame())
	if b.animation.IsFinished() {
		b.finished = true
	}

	if ds := b.getActiveDamageSource(); ds != nil {
		return UpdateResult{Actions: []Action{{
			Type:         ActionCreateDamageSource,
			Location:     b.Location,
			Direction:    ZeroVector(),
			DamageSource: ds}}}
	}
	return UpdateResult{}
}

func (b *BombExplosion) CanRemove() bool {
	return b.finished
}

func (b *BombExplosion) getActiveDamageSource() *DamageSource {
	animIndex := b.animation.frameIndex

	if config, ok := b.attackHitboxes[animIndex]; ok {
		worldHitbox := config.HitBox.Offset(b.X, b.Y)
		return NewDamageSource(TagPlayer, worldHitbox, config.Damage)
	}

	return nil
}
