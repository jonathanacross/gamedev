package weapons

import (
	"time"

	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bomb struct {
	cooldownTimer *core.Timer
	level         int
}

func NewBomb() *Bomb {
	return &Bomb{
		cooldownTimer: core.NewTimer(750 * time.Millisecond),
		level:         1,
	}
}

func (b *Bomb) SetLevel(level int) {
	b.level = level
}

func (b *Bomb) Update(ctx core.PlayerContext, level core.Level) {
	b.cooldownTimer.Update()
}

func (b *Bomb) OnAttack(ctx core.PlayerContext) {
	if b.cooldownTimer.IsReady() {
		b.cooldownTimer.Reset()

		loc := ctx.Location()
		dir := ctx.GetDirection()

		spawnLoc := core.Vector(loc).Plus(core.DirectionToVector(dir).Scale(float64(core.TileSize)))

		ctx.AddAction(core.Action{
			Type:     core.ActionDropBomb,
			Location: core.Location(spawnLoc),
			Level:    b.level,
		})
	}
}

func (b *Bomb) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	return core.RenderOrderFront
}

func (b *Bomb) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	// No visual for holding a bomb.
}

func (b *Bomb) GetCooldownProgress() float64 {
	return b.cooldownTimer.GetProgress()
}
