package objects

import (
	"time"

	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bow struct {
	cooldownTimer    *core.Timer
	isAttacking      bool
	currentAnimation *core.Animation
	spriteSheet      *core.SpriteSheet
	animations       map[core.Direction]*core.Animation
}

func NewBow() *Bow {
	// Reusing Player Bow Animation Logic
	// ssColumns := 8, ssRows := 6 (from Player)
	spriteSheet := core.NewSpriteSheet(48, 64, 8, 6)

	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := core.NewDirectionAnimationMap([]int{0, 1}, 0, directionOffsets, 6, false)

	return &Bow{
		// BowCooldown was 300ms in player.go
		cooldownTimer: core.NewTimer(300 * time.Millisecond),
		spriteSheet:   spriteSheet,
		animations:    animations,
		isAttacking:   false,
	}
}

func (b *Bow) Update(ctx core.PlayerContext) {
	b.cooldownTimer.Update()

	if b.isAttacking && b.currentAnimation != nil {
		b.currentAnimation.Update()
		if b.currentAnimation.IsFinished() {
			// Animation finished -> Shoot
			// Calculate spawn location
			loc := ctx.Location()
			dir := ctx.GetDirection()
			arrowStartOffset := core.Vector{X: 0, Y: -6}
			spawnLoc := core.Vector(loc).Plus(arrowStartOffset).Plus(core.DirectionToVector(dir).Scale(float64(core.TileSize)))

			ctx.AddAction(core.Action{
				Type:      core.ActionShootArrow,
				Location:  core.Location(spawnLoc),
				Direction: core.DirectionToVector(dir),
			})

			b.isAttacking = false
			b.currentAnimation.Reset()
			b.currentAnimation = nil
		}
	}
}

func (b *Bow) OnAttack(ctx core.PlayerContext) {
	if b.cooldownTimer.IsReady() && !b.isAttacking {
		b.cooldownTimer.Reset()
		b.isAttacking = true
		dir := ctx.GetDirection()
		b.currentAnimation = b.animations[dir]
		b.currentAnimation.Reset()
	}
}

func (b *Bow) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	if playerDirection == core.Up {
		return core.RenderOrderBehind
	}
	return core.RenderOrderFront
}

func (b *Bow) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	if !b.isAttacking || b.currentAnimation == nil {
		return
	}

	frame := b.currentAnimation.Frame()
	srcRect := b.spriteSheet.Rect(frame)

	// Draw Offset (Matches Player logic)
	drawIdxX := 25.0
	drawIdxY := 38.0

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerPos.X-drawIdxX, playerPos.Y-drawIdxY)
	op.GeoM.Concat(cameraMatrix)

	img := assets.PlayerAttackBowSpritesImage.SubImage(srcRect).(*ebiten.Image)
	screen.DrawImage(img, op)
}
