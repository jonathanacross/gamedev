package objects

import (
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Shield struct {
	spriteSheet    *core.SpriteSheet
	animations     map[core.Direction]*core.Animation
	attackHitboxes map[core.Direction]map[int]core.DamageSourceConfig

	isBlocking bool // Input received this frame
	isActive   bool // State for Update/Draw (latched from Input)

	currentAnimation *core.Animation
}

func NewShield() *Shield {
	spriteSheet := core.NewSpriteSheet(48, 64, 8, 6)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{1}, 0, directionOffsets, 6, true)

	var baseShieldDamage = 0
	shieldAttackBoxes := map[core.Direction]map[int]core.DamageSourceConfig{
		core.Down: {
			0: {HitBox: core.Rect{Left: -3, Top: -6, Right: 7, Bottom: 7}, Damage: baseShieldDamage},
			1: {HitBox: core.Rect{Left: -3, Top: -6, Right: 7, Bottom: 7}, Damage: baseShieldDamage},
		},
		core.Left: {
			0: {HitBox: core.Rect{Left: -9, Top: -7, Right: -1, Bottom: 6}, Damage: baseShieldDamage},
			1: {HitBox: core.Rect{Left: -9, Top: -7, Right: -1, Bottom: 6}, Damage: baseShieldDamage},
		},
		core.Right: {
			0: {HitBox: core.Rect{Left: 1, Top: -7, Right: 9, Bottom: 6}, Damage: baseShieldDamage},
			1: {HitBox: core.Rect{Left: 1, Top: -7, Right: 9, Bottom: 7}, Damage: baseShieldDamage},
		},
		core.Up: {
			0: {HitBox: core.Rect{Left: -9, Top: -12, Right: 1, Bottom: 1}, Damage: baseShieldDamage},
			1: {HitBox: core.Rect{Left: -9, Top: -12, Right: 1, Bottom: 1}, Damage: baseShieldDamage},
		},
	}

	return &Shield{
		spriteSheet:    spriteSheet,
		animations:     animations,
		attackHitboxes: shieldAttackBoxes,
		isBlocking:     false,
		isActive:       false,
	}
}

func (s *Shield) Update(ctx core.PlayerContext) {
	// Latch input state to active state for this frame logic
	s.isActive = s.isBlocking
	// Reset input state for next frame
	s.isBlocking = false

	if s.isActive {
		ctx.SetBlocking(true)

		if s.currentAnimation != nil {
			s.currentAnimation.Update()
		}

		frame := 0
		dir := ctx.GetDirection()
		if dirBoxes, ok := s.attackHitboxes[dir]; ok {
			if config, ok := dirBoxes[frame]; ok {
				playerLoc := ctx.Location()
				worldHitbox := config.HitBox.Offset(playerLoc.X, playerLoc.Y)
				ds := core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypeImpact, config.Damage)
				ctx.CreateDamageSource(ds)
			}
		}
	} else {
		ctx.SetBlocking(false)
		if s.currentAnimation != nil {
			s.currentAnimation.Reset()
			s.currentAnimation = nil
		}
	}
}

func (s *Shield) OnAttack(ctx core.PlayerContext) {
	s.isBlocking = true
	if s.currentAnimation == nil {
		dir := ctx.GetDirection()
		s.currentAnimation = s.animations[dir]
		s.currentAnimation.Reset()
	}
}

func (s *Shield) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (s *Shield) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !s.isActive {
		return
	}

	if s.currentAnimation == nil {
		return
	}

	frame := s.currentAnimation.Frame()
	srcRect := s.spriteSheet.Rect(frame)

	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackShieldSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}
