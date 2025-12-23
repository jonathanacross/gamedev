package objects

import (
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sword struct {
	spriteSheet    *core.SpriteSheet
	animations     map[core.Direction]*core.Animation
	attackHitboxes map[core.Direction]map[int]core.DamageSourceConfig
	isAttacking    bool
}

func NewSword() *Sword {
	// ssColumns := 8, ssRows := 6 (from Player)
	// Actually we might want a dedicated spritesheet for just the weapon later,
	// but for now we are using the player spritesheet and drawing the "baked in" weapon.
	// This means we are re-drawing the player Holding the weapon.
	// This is the "Double Drawing" compromise.
	spriteSheet := core.NewSpriteSheet(48, 64, 8, 6)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false)

	// Hitboxes (Copied from Player)
	var baseSwordDamage = 1
	swordAttackBoxes := map[core.Direction]map[int]core.DamageSourceConfig{
		core.Down: {
			1: {HitBox: core.Rect{Left: 6, Top: 15, Right: 24, Bottom: 30}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: -10, Top: 25, Right: 10, Bottom: 35}, Damage: baseSwordDamage},
			3: {HitBox: core.Rect{Left: -24, Top: 10, Right: -6, Bottom: 30}, Damage: baseSwordDamage},
		},
		core.Left: {
			1: {HitBox: core.Rect{Left: -24, Top: -10, Right: -5, Bottom: 5}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: -24, Top: 5, Right: -10, Bottom: 20}, Damage: baseSwordDamage},
			3: {HitBox: core.Rect{Left: -20, Top: 15, Right: -4, Bottom: 27}, Damage: baseSwordDamage},
		},
		core.Right: {
			1: {HitBox: core.Rect{Left: 5, Top: -10, Right: 24, Bottom: 5}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: 10, Top: 5, Right: 24, Bottom: 20}, Damage: baseSwordDamage},
			3: {HitBox: core.Rect{Left: 4, Top: 15, Right: 20, Bottom: 27}, Damage: baseSwordDamage},
		},
		core.Up: {
			1: {HitBox: core.Rect{Left: -24, Top: -15, Right: -6, Bottom: 2}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: -10, Top: -24, Right: 10, Bottom: -15}, Damage: baseSwordDamage},
			3: {HitBox: core.Rect{Left: 6, Top: -15, Right: 24, Bottom: 2}, Damage: baseSwordDamage},
		},
	}

	return &Sword{
		spriteSheet:    spriteSheet,
		animations:     animations,
		attackHitboxes: swordAttackBoxes,
		isAttacking:    false,
	}
}

func (s *Sword) Update(ctx core.PlayerContext) {
	if s.isAttacking {
		anim := s.animations[ctx.GetDirection()]
		anim.Update()
		if anim.IsFinished() {
			s.isAttacking = false
			anim.Reset()
		} else {
			// Generate Damage Source
			frame := anim.CurrentFrameIndex()
			dir := ctx.GetDirection()
			if dirBoxes, ok := s.attackHitboxes[dir]; ok {
				if config, ok := dirBoxes[frame]; ok {
					playerLoc := ctx.Location()
					worldHitbox := config.HitBox.Offset(playerLoc.X, playerLoc.Y)
					ds := core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypePhysical, config.Damage)
					ctx.CreateDamageSource(ds)
				}
			}
		}
	}
}

func (s *Sword) OnAttack(ctx core.PlayerContext) {
	if !s.isAttacking {
		s.isAttacking = true
		anim := s.animations[ctx.GetDirection()]
		anim.Reset()
	}
}

func (s *Sword) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (s *Sword) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !s.isAttacking {
		return
	}

	// Use weapon's animation frame
	anim := s.animations[playerDir]
	frame := anim.Frame()
	srcRect := s.spriteSheet.Rect(frame)

	// Offset needs to be relative to player position.
	// Player draw offset is 25, 38.
	// If we are drawing the same sprite sheet, we can reuse the offset?
	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackSwordSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}

func (s *Sword) GetCooldownProgress() float64 {
	return 1.0
}
