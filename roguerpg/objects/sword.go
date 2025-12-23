package objects

import (
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sword struct {
	spriteSheet      *core.SpriteSheet
	animations       map[core.Direction]*core.Animation
	attackHitboxes   map[core.Direction]map[int]core.DamageSourceConfig
	isAttacking      bool
	currentAnimation *core.Animation
}

func NewSword() *Sword {
	ssColumns := 8
	ssRows := 6
	spriteSheet := core.NewSpriteSheet(48, 64, ssColumns, ssRows)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false)

	// Copied from player.go
	baseSwordDamage := 1
	swordAttackBoxes := map[core.Direction]map[int]core.DamageSourceConfig{
		core.Down: {
			1: {HitBox: core.Rect{Left: -16, Top: 0, Right: 16, Bottom: 22}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: -16, Top: 0, Right: 16, Bottom: 22}, Damage: baseSwordDamage},
		},
		core.Left: {
			1: {HitBox: core.Rect{Left: -22, Top: -24, Right: 0, Bottom: 8}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: -22, Top: -24, Right: 0, Bottom: 8}, Damage: baseSwordDamage},
		},
		core.Right: {
			1: {HitBox: core.Rect{Left: 0, Top: -24, Right: 22, Bottom: 8}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: 0, Top: -24, Right: 22, Bottom: 8}, Damage: baseSwordDamage},
		},
		core.Up: {
			1: {HitBox: core.Rect{Left: -16, Top: -22, Right: 16, Bottom: 0}, Damage: baseSwordDamage},
			2: {HitBox: core.Rect{Left: -16, Top: -22, Right: 16, Bottom: 0}, Damage: baseSwordDamage},
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
	if s.isAttacking && s.currentAnimation != nil {
		s.currentAnimation.Update()
		if s.currentAnimation.IsFinished() {
			s.isAttacking = false
			s.currentAnimation.Reset()
			s.currentAnimation = nil
		} else {
			// Check for active hitboxes in current frame
			frame := s.currentAnimation.CurrentFrameIndex()
			dir := ctx.GetDirection()
			if dirBoxes, ok := s.attackHitboxes[dir]; ok {
				if config, ok := dirBoxes[frame]; ok {
					// Create DamageSource
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
		dir := ctx.GetDirection()
		s.currentAnimation = s.animations[dir]
		s.currentAnimation.Reset()
	}
}

func (s *Sword) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (s *Sword) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !s.isAttacking || s.currentAnimation == nil {
		return
	}

	// Use specific weapon sprites if available, otherwise reuse the player sheet
	// (Note: Currently the player sheet HAS the sword drawn on it, so we are double drawing.
	// This is known and acceptable for this refactor phase).
	// We will draw from the same sprite sheet using the sword animation frame.

	frame := s.currentAnimation.Frame()
	srcRect := s.spriteSheet.Rect(frame)

	// Draw Offset (Matches Player logic)
	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	// Using the PlayerAttackSwordSpritesImage (which is what the player used)
	// Ideally we would have a 'SwordOnlySpritesImage', but for now we use the full sheet
	// which effectively draws the player holding the sword again.
	// As per plan: "The existing sprites might look double-drawn... This is known."
	img := assets.PlayerAttackSwordSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}
