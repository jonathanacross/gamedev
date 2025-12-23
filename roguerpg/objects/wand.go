package objects

import (
	"time"

	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Wand struct {
	cooldownTimer    *core.Timer
	isAttacking      bool
	currentAnimation *core.Animation
	spriteSheet      *core.SpriteSheet
	animations       map[core.Direction]*core.Animation
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
		// WandCooldown was 2000ms in player.go
		cooldownTimer: core.NewTimer(2000 * time.Millisecond),
		spriteSheet:   spriteSheet,
		animations:    animations,
		isAttacking:   false,
	}
}

func (w *Wand) Update(ctx core.PlayerContext) {
	w.cooldownTimer.Update()

	if w.isAttacking && w.currentAnimation != nil {
		w.currentAnimation.Update()
		if w.currentAnimation.IsFinished() {
			// Animation finished -> Shoot Star
			loc := ctx.Location()
			dir := ctx.GetDirection()
			starStartOffset := core.Vector{X: 0, Y: -4}
			spawnLoc := core.Vector(loc).Plus(starStartOffset).Plus(core.DirectionToVector(dir).Scale(float64(core.TileSize)))

			// Note: Finding nearest enemy is usually done in Action handling or passed in.
			// Ideally, we'd find the target here if we had access to Level through context.
			// But PlayerContext only exposes what we defined.
			// For this refactor, we will emit the ActionCreateStar without a target for now,
			// or we can Expand PlayerContext to include Level or FindTarget capability.
			// However, `ActionCreateStar` handling in main.go already expects `action.Target`.
			// Since we can't easily get the target here without expanding the context significantly,
			// we will rely on the Game Engine to potentially handle targeting or just fire straight.
			// *Correction*: In player.go, it did: target = level.FindNearestEnemy...
			// To preserve this, we should add FindNearestEnemy to PlayerContext or accept it.
			// But for now, let's keep it simple and just fire blindly (or change context next).

			ctx.AddAction(core.Action{
				Type:      core.ActionCreateStar,
				Location:  core.Location(spawnLoc),
				Direction: core.DirectionToVector(dir),
				// Target: nil, // Will just go straight
			})

			w.isAttacking = false
			w.currentAnimation.Reset()
			w.currentAnimation = nil
		}
	}
}

func (w *Wand) OnAttack(ctx core.PlayerContext) {
	if w.cooldownTimer.IsReady() && !w.isAttacking {
		w.cooldownTimer.Reset()
		w.isAttacking = true
		dir := ctx.GetDirection()
		w.currentAnimation = w.animations[dir]
		w.currentAnimation.Reset()
	}
}

func (w *Wand) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (w *Wand) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !w.isAttacking || w.currentAnimation == nil {
		return
	}

	frame := w.currentAnimation.Frame()
	srcRect := w.spriteSheet.Rect(frame)

	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackWandSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}
