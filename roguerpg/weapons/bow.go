package weapons

import (
	"time"

	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bow struct {
	spriteSheet   *core.SpriteSheet
	animations    map[core.Direction]*core.Animation
	cooldownTimer *core.Timer
	isAttacking   bool
	level         int
}

func NewBow() *Bow {
	spriteSheet := core.NewSpriteSheet(48, 64, 8, 6)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{0, 1}, 0, directionOffsets, 6, false)

	return &Bow{
		spriteSheet:   spriteSheet,
		animations:    animations,
		cooldownTimer: core.NewTimer(500 * time.Millisecond),
		isAttacking:   false,
		level:         1,
	}
}

func (b *Bow) SetLevel(level int) {
	b.level = level
	var arrowCooldown time.Duration
	switch level {
	case 1:
		arrowCooldown = 500 * time.Millisecond
	case 2:
		arrowCooldown = 400 * time.Millisecond
	default:
		arrowCooldown = 300 * time.Millisecond
	}
	b.cooldownTimer = core.NewTimer(arrowCooldown)
}

func (b *Bow) Update(ctx core.PlayerContext, level core.Level) {
	b.cooldownTimer.Update()

	if b.isAttacking {
		anim := b.animations[ctx.GetDirection()]
		anim.Update()

		if anim.IsFinished() {
			b.isAttacking = false
			anim.Reset()

			ctx.AddAction(core.Action{
				Type:      core.ActionShootArrow,
				Location:  ctx.Location(),
				Direction: core.DirectionToVector(ctx.GetDirection()),
				Level:     b.level,
			})
		}
	}
}

func (b *Bow) OnAttack(ctx core.PlayerContext) {
	if b.cooldownTimer.IsReady() && !b.isAttacking {
		b.cooldownTimer.Reset()
		b.isAttacking = true
		anim := b.animations[ctx.GetDirection()]
		anim.Reset()
	}
}

func (b *Bow) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (b *Bow) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !b.isAttacking {
		return
	}

	anim := b.animations[playerDir]
	frame := anim.Frame()
	srcRect := b.spriteSheet.Rect(frame)

	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackBowSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}

func (b *Bow) GetCooldownProgress() float64 {
	return b.cooldownTimer.GetProgress()
}
