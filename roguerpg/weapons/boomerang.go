package weapons

import (
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Boomerang struct {
	hasBoomerang bool
}

func NewBoomerang() *Boomerang {
	return &Boomerang{
		hasBoomerang: true,
	}
}

func (b *Boomerang) Update(ctx core.PlayerContext) {
	// No update logic needed, purely state based
}

func (b *Boomerang) OnAttack(ctx core.PlayerContext) {
	if b.hasBoomerang {
		b.hasBoomerang = false

		loc := ctx.Location()
		dir := ctx.GetDirection()

		spawnLoc := core.Vector(loc).Plus(core.DirectionToVector(dir).Scale(float64(core.TileSize)))

		ctx.AddAction(core.Action{
			Type:      core.ActionThrowBoomerang,
			Location:  core.Location(spawnLoc),
			Direction: core.DirectionToVector(dir),
		})
	}
}

func (b *Boomerang) Catch() {
	b.hasBoomerang = true
}

func (b *Boomerang) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	return core.RenderOrderFront
}

func (b *Boomerang) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	// No visual for holding yet
}

func (b *Boomerang) GetCooldownProgress() float64 {
	if b.hasBoomerang {
		return 1.0
	}
	return 0.0
}
