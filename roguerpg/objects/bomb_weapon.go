package objects

import (
	"time"

	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type BombWeapon struct {
	cooldownTimer *core.Timer
}

func NewBombWeapon() *BombWeapon {
	// Cooldown was 750ms in player.go
	return &BombWeapon{
		cooldownTimer: core.NewTimer(750 * time.Millisecond),
	}
}

func (b *BombWeapon) Update(ctx core.PlayerContext) {
	b.cooldownTimer.Update()
}

func (b *BombWeapon) OnAttack(ctx core.PlayerContext) {
	if b.cooldownTimer.IsReady() {
		b.cooldownTimer.Reset()

		loc := ctx.Location()
		dir := ctx.GetDirection()

		// Calculate spawn location (logic from player.go)
		spawnLoc := core.Vector(loc).Plus(core.DirectionToVector(dir).Scale(float64(core.TileSize)))

		ctx.AddAction(core.Action{
			Type:      core.ActionDropBomb,
			Location:  core.Location(spawnLoc),
			Direction: core.DirectionToVector(dir),
		})
	}
}

func (b *BombWeapon) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	return core.RenderOrderFront
}

func (b *BombWeapon) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	// No visual for holding the bomb yet (it just appears)
}

func (b *BombWeapon) GetCooldownProgress() float64 {
	return b.cooldownTimer.GetProgress()
}
