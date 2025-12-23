package objects

import (
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type BoomerangWeapon struct {
	hasBoomerang bool
}

func NewBoomerangWeapon() *BoomerangWeapon {
	return &BoomerangWeapon{
		hasBoomerang: true,
	}
}

func (b *BoomerangWeapon) Update(ctx core.PlayerContext) {
	// No update logic needed, purely state based
}

func (b *BoomerangWeapon) OnAttack(ctx core.PlayerContext) {
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

func (b *BoomerangWeapon) Catch() {
	b.hasBoomerang = true
}

func (b *BoomerangWeapon) GetRenderOrder(playerDirection core.Direction) core.RenderOrder {
	return core.RenderOrderFront
}

func (b *BoomerangWeapon) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos core.Location, playerDir core.Direction, playerFrame int) {
	// No visual for holding yet
}

func (b *BoomerangWeapon) GetCooldownProgress() float64 {
	if b.hasBoomerang {
		return 1.0
	}
	return 0.0
}
