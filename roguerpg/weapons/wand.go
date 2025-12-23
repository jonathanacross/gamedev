package weapons

import (
	"time"

	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Wand struct {
	spriteSheet   *core.SpriteSheet
	animations    map[core.Direction]*core.Animation
	cooldownTimer *core.Timer
	isAttacking   bool
}

func NewWand() *Wand {
	spriteSheet := core.NewSpriteSheet(48, 64, 8, 6)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{0, 1, 2}, 0, directionOffsets, 4, false)

	return &Wand{
		spriteSheet:   spriteSheet,
		animations:    animations,
		cooldownTimer: core.NewTimer(1000 * time.Millisecond),
		isAttacking:   false,
	}
}

func (w *Wand) Update(ctx core.PlayerContext, level core.Level) {
	w.cooldownTimer.Update()

	if w.isAttacking {
		anim := w.animations[ctx.GetDirection()]
		anim.Update()

		if anim.IsFinished() {
			w.isAttacking = false
			anim.Reset()

			loc := ctx.Location()
			dir := ctx.GetDirection()
			spawnLoc := core.Vector(loc).Plus(core.DirectionToVector(dir).Scale(float64(core.TileSize)))
			target := level.FindNearestEnemy(loc, core.TileSize*6)

			ctx.AddAction(core.Action{
				Type:      core.ActionCreateStar,
				Location:  core.Location(spawnLoc),
				Direction: core.DirectionToVector(dir),
				Target:    target,
			})
		}
	}
}

func (w *Wand) OnAttack(ctx core.PlayerContext) {
	if w.cooldownTimer.IsReady() && !w.isAttacking {
		w.cooldownTimer.Reset()
		w.isAttacking = true
		anim := w.animations[ctx.GetDirection()]
		anim.Reset()
	}
}

func (w *Wand) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (w *Wand) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !w.isAttacking {
		return
	}

	anim := w.animations[playerDir]
	frame := anim.Frame()
	srcRect := w.spriteSheet.Rect(frame)

	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackWandSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}

func (w *Wand) GetCooldownProgress() float64 {
	return w.cooldownTimer.GetProgress()
}
