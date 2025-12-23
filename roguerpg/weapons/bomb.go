package weapons

import (
	"time"

	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bomb struct {
	cooldownTimer *core.Timer
}

func NewBomb() *Bomb {
	// Cooldown was 750ms in player.go
	return &Bomb{
		cooldownTimer: core.NewTimer(750 * time.Millisecond),
	}
}

func (b *Bomb) Update(ctx core.PlayerContext) {
	b.cooldownTimer.Update()
}

func (b *Bomb) OnAttack(ctx core.PlayerContext) {
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

func (b *Bomb) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	return core.RenderOrderFront
}

func (b *Bomb) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	// No visual for holding the bomb yet (it just appears)
}

func (b *Bomb) GetCooldownProgress() float64 {
	return b.cooldownTimer.GetProgress()
}
