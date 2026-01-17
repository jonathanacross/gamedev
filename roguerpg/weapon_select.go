package main

import (
	"image/color"
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	weaponSelectWindowWidth  = 148
	weaponSelectWindowHeight = 72
)

// Implements GameState
type WeaponSelector struct {
	Window           *core.Window
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

	windowX := float64(ScreenWidth-weaponSelectWindowWidth) / 2.0
	windowY := float64(ScreenHeight-weaponSelectWindowHeight) / 2.0

	firstStart := 8
	xStart := 40
	yStart := 40
	weaponSpacing := assets.UiSelectRectImage.Bounds().Dx()

	primaryWeapon := NewWeaponBox(core.Location{
		X: windowX + float64(firstStart),
		Y: windowY + float64(yStart)},
		core.WeaponSword, false)
	secondaryWeapons := []*WeaponBox{}
	weaponIndex := 0
	foundWeapon := false
	currentWeapon := player.GetSecondaryWeapon()
	for i, weapon := range weaponTable {
		wb := NewWeaponBox(
			core.Location{X: windowX + float64(xStart+i*weaponSpacing),
				Y: windowY + float64(yStart)}, weapon, false)
		secondaryWeapons = append(secondaryWeapons, wb)

		if weapon == currentWeapon && !foundWeapon {
			weaponIndex = i
			foundWeapon = true
		}
	}
	return &WeaponSelector{
		Window: core.NewWindow(core.Rect{
			Left:   windowX,
			Top:    windowY,
			Right:  windowX + weaponSelectWindowWidth,
			Bottom: windowY + weaponSelectWindowHeight,
		}),
		weaponIndex:      weaponIndex,
		weaponTable:      weaponTable,
		primaryWeapon:    primaryWeapon,
		secondaryWeapons: secondaryWeapons,
	}
}

func (w *WeaponSelector) Draw(screen *ebiten.Image, context *core.GameContext) {
	w.Window.Draw(screen)

	// Draw Text
	titleLoc := core.Location{X: w.Window.Rect.Left + 8, Y: w.Window.Rect.Top + 8}
	drawTextAt(screen, "Select your weapon", titleLoc.X, titleLoc.Y, text.AlignStart, color.White)

	primaryButtonLoc := core.Location{X: w.Window.Rect.Left + 14, Y: w.Window.Rect.Top + 24}
	drawTextAt(screen, "X", primaryButtonLoc.X, primaryButtonLoc.Y, text.AlignStart, color.White)

	secondaryButtonLoc := core.Location{X: w.Window.Rect.Left + 80, Y: w.Window.Rect.Top + 24}
	drawTextAt(screen, "Z", secondaryButtonLoc.X, secondaryButtonLoc.Y, text.AlignStart, color.White)

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
