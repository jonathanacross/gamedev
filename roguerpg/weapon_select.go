package main

import (
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Implements GameState
type WeaponSelector struct {
	windowImage      *ebiten.Image
	windowLoc        core.Location
	weaponIndex      int
	weaponTable      []core.WeaponType
	primaryWeapon    *WeaponBox
	secondaryWeapons []*WeaponBox
}

func getWeaponType(player core.Player, upgradeType core.UpgradeType, weaponType core.WeaponType) core.WeaponType {
	playerStats := player.GetCurrentStats()
	if playerStats[upgradeType] > 0 {
		return weaponType
	}
	return core.WeaponNone
}

func NewWeaponSelector(player core.Player) *WeaponSelector {
	weaponTable := []core.WeaponType{
		getWeaponType(player, core.UpgradeTypeBoomerang, core.WeaponBoomerang),
		getWeaponType(player, core.UpgradeTypeBomb, core.WeaponBomb),
		getWeaponType(player, core.UpgradeTypeShield, core.WeaponShield),
		getWeaponType(player, core.UpgradeTypeBow, core.WeaponBow),
		getWeaponType(player, core.UpgradeTypeWand, core.WeaponWand),
	}

	windowX := float64(ScreenWidth-assets.UiWeaponSelectWindowImage.Bounds().Dx()) / 2.0
	windowY := float64(ScreenHeight-assets.UiWeaponSelectWindowImage.Bounds().Dy()) / 2.0

	firstStart := 8
	xStart := 40
	yStart := 36
	weaponSpacing := assets.UiSelectRectImage.Bounds().Dx()

	primaryWeapon := NewWeaponBox(core.Location{
		X: windowX + float64(firstStart),
		Y: windowY + float64(yStart)},
		core.WeaponSword, false)
	secondaryWeapons := []*WeaponBox{}
	for i, weapon := range weaponTable {
		wb := NewWeaponBox(
			core.Location{X: windowX + float64(xStart+i*weaponSpacing),
				Y: windowY + float64(yStart)}, weapon, false)
		secondaryWeapons = append(secondaryWeapons, wb)
	}
	return &WeaponSelector{
		windowImage:      assets.UiWeaponSelectWindowImage,
		windowLoc:        core.Location{X: windowX, Y: windowY},
		weaponTable:      weaponTable,
		primaryWeapon:    primaryWeapon,
		secondaryWeapons: secondaryWeapons,
	}
}

func (w *WeaponSelector) Draw(screen *ebiten.Image, context *core.GameContext) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.windowLoc.X, w.windowLoc.Y)

	screen.DrawImage(w.windowImage, op)

	// Draw sword
	w.primaryWeapon.Draw(screen)

	// Draw other weapons
	for i, _ := range w.weaponTable {
		w.secondaryWeapons[i].SetShowFrame(i == w.weaponIndex)
		w.secondaryWeapons[i].Draw(screen)
	}
}

func (w *WeaponSelector) Update(context *core.GameContext) []core.Action {
	numWeapons := len(w.weaponTable)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return []core.Action{
			{Type: core.ActionPopState},
			{Type: core.ActionSwitchWeapon, WeaponType: w.weaponTable[w.weaponIndex]},
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		w.weaponIndex = (w.weaponIndex + 1) % numWeapons
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		w.weaponIndex = (w.weaponIndex - 1 + numWeapons) % numWeapons
	}
	return []core.Action{}
}
