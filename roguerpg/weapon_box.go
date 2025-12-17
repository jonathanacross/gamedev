package main

import (
	"image/color"
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	UiIconSelect = iota
	UiIconSword
	UiIconBoomerang
	UiIconShield
	UiIconBomb
	UiIconBow
	UiIconWand
	UiIconEmpty
)

var progressBarImage *ebiten.Image = createProgressBarImage(assets.UiSelectRectImage.Bounds().Dx(), 2)

type WeaponBox struct {
	Loc              core.Location
	WeaponType       core.WeaponType
	DrawFrame        bool
	ReadyPercent     float64
	frameImage       *ebiten.Image
	iconsImage       *ebiten.Image
	iconsSpriteSheet *core.SpriteSheet
}

func NewWeaponBox(location core.Location, weaponType core.WeaponType, drawFrame bool) *WeaponBox {
	return &WeaponBox{
		Loc:              location,
		WeaponType:       weaponType,
		DrawFrame:        drawFrame,
		frameImage:       assets.UiSelectRectImage,
		iconsImage:       assets.UiIconsImage,
		iconsSpriteSheet: core.NewSpriteSheet(16, 16, 8, 1),
	}
}

func getWeaponIconIndex(weaponType core.WeaponType) int {
	switch weaponType {
	case core.WeaponSword:
		return UiIconSword
	case core.WeaponBoomerang:
		return UiIconBoomerang
	case core.WeaponShield:
		return UiIconShield
	case core.WeaponBomb:
		return UiIconBomb
	case core.WeaponBow:
		return UiIconBow
	case core.WeaponWand:
		return UiIconWand
	default:
		return UiIconEmpty
	}
}

func createProgressBarImage(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})
	return img
}

func (w *WeaponBox) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.Loc.X), float64(w.Loc.Y))

	iconIndex := getWeaponIconIndex(w.WeaponType)
	icon := w.iconsImage.SubImage(w.iconsSpriteSheet.Rect(iconIndex)).(*ebiten.Image)
	screen.DrawImage(icon, op)
	frameOffset := -2

	if w.DrawFrame {
		op.GeoM.Translate(float64(frameOffset), float64(frameOffset))
		screen.DrawImage(w.frameImage, op)
	}

	if w.ReadyPercent > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(w.ReadyPercent, 1.0)
		op.GeoM.Translate(float64(w.Loc.X), float64(w.Loc.Y))
		op.GeoM.Translate(float64(frameOffset), float64(20))
		screen.DrawImage(progressBarImage, op)
	}
}

func (w *WeaponBox) SetShowFrame(show bool) {
	w.DrawFrame = show
}

func (w *WeaponBox) SetWeaponType(weaponType core.WeaponType) {
	w.WeaponType = weaponType
}
