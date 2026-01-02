package weapons

import (
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sword struct {
	spriteSheet    *core.SpriteSheet
	animations     map[core.Direction]*core.Animation
	attackHitboxes map[core.Direction]map[int]core.Rect
	level          int
	damage         int
	isAttacking    bool
}

func NewSword() *Sword {
	spriteSheet := core.NewSpriteSheet(48, 64, 8, 6)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false)

	swordAttackBoxes := map[core.Direction]map[int]core.Rect{
		core.Up: {
			1: core.Rect{Left: -16, Top: -26, Right: 16, Bottom: -4},
			2: core.Rect{Left: -16, Top: -26, Right: 16, Bottom: -4},
		},
		core.Down: {
			1: core.Rect{Left: -16, Top: -4, Right: 16, Bottom: 18},
			2: core.Rect{Left: -16, Top: -4, Right: 16, Bottom: 18},
		},
		core.Left: {
			1: core.Rect{Left: -22, Top: -24, Right: 0, Bottom: 8},
			2: core.Rect{Left: -22, Top: -24, Right: 0, Bottom: 8},
		},
		core.Right: {
			1: core.Rect{Left: 0, Top: -24, Right: 22, Bottom: 8},
			2: core.Rect{Left: 0, Top: -24, Right: 22, Bottom: 8},
		},
	}

	return &Sword{
		spriteSheet:    spriteSheet,
		animations:     animations,
		attackHitboxes: swordAttackBoxes,
		level:          1,
		damage:         1,
		isAttacking:    false,
	}
}

func (s *Sword) SetLevel(level int) {
	level = core.Clamp(level, 1, 3)
	s.level = level
	s.damage = 1 << level
}

func (s *Sword) Update(ctx core.PlayerContext, level core.Level) {
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
				if hitBox, ok := dirBoxes[frame]; ok {
					playerLoc := ctx.Location()
					worldHitbox := hitBox.Offset(playerLoc.X, playerLoc.Y)
					ds := core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypePhysical, s.damage)
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

	// Offset same as player offset.
	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackSwordSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}

func (s *Sword) GetCooldownProgress() float64 {
	// Sword has no cooldown. Always ready.
	return 1.0
}

func (s *Sword) IsAttacking() bool {
	return s.isAttacking
}
