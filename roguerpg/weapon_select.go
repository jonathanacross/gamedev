package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type WeaponType int

const (
	WeaponNone WeaponType = iota
	WeaponSword
	WeaponBoomerang
	WeaponBomb
	WeaponShield
	WeaponArrows
	WeaponWand

	UiIconSelect    = 0
	UiIconSword     = 1
	UiIconBoomerang = 2
	UiIconShield    = 4
	UiIconBomb      = 3
	UiIconArrows    = 5
	UiIconWand      = 6
	UiIconEmpty     = 7
)

type WeaponSlot struct {
	Type      WeaponType
	IconIndex int
}

// Implements GameState
type WeaponSelector struct {
	windowImage      *ebiten.Image
	iconsImage       *ebiten.Image
	iconsSpriteSheet *SpriteSheet
	windowLoc        Location
	weaponIndex      int
	weaponTable      []WeaponSlot
}

var WeaponSelectorInstance *WeaponSelector = NewWeaponSelector()

func NewWeaponSelector() *WeaponSelector {
	windowX := float64(ScreenWidth-WeaponSelectWindowImage.Bounds().Dx()) / 2.0
	windowY := float64(ScreenHeight-WeaponSelectWindowImage.Bounds().Dy()) / 2.0
	return &WeaponSelector{
		windowImage:      WeaponSelectWindowImage,
		iconsImage:       UiIconsImage,
		iconsSpriteSheet: NewSpriteSheet(16, 16, 8, 1),
		windowLoc:        Location{X: windowX, Y: windowY},
		weaponTable: []WeaponSlot{
			{Type: WeaponBoomerang, IconIndex: UiIconBoomerang},
			{Type: WeaponBomb, IconIndex: UiIconBomb},
			{Type: WeaponShield, IconIndex: UiIconShield},
			// TODO: Add wand and arrows when implemented
			{Type: WeaponNone, IconIndex: UiIconEmpty},
			{Type: WeaponNone, IconIndex: UiIconEmpty},
		},
	}
}

func (w *WeaponSelector) Draw(screen *ebiten.Image, context *GameContext) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.windowLoc.X, w.windowLoc.Y)

	screen.DrawImage(w.windowImage, op)

	xStart := 40
	yStart := 36
	// Draw sword
	w.drawIcon(screen, UiIconSword, 8, float64(yStart))
	// Draw other weapons
	for i, ws := range w.weaponTable {
		w.drawIcon(screen, ws.IconIndex, float64(xStart+i*TileSize), float64(yStart))
	}
	// Highlight selected weapon
	w.drawIcon(screen, UiIconSelect, float64(xStart+w.weaponIndex*TileSize), float64(yStart))
}

// Draws the icon at the given (x,y) relative to the window location
func (w *WeaponSelector) drawIcon(screen *ebiten.Image, index int, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.windowLoc.X+x, w.windowLoc.Y+y)
	icon := w.iconsImage.SubImage(w.iconsSpriteSheet.Rect(index)).(*ebiten.Image)
	screen.DrawImage(icon, op)
}

func (w *WeaponSelector) Update(context *GameContext) []Action {
	numWeapons := len(w.weaponTable)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return []Action{{Type: ActionPopState}}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		w.weaponIndex = (w.weaponIndex + 1) % numWeapons
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		w.weaponIndex = (w.weaponIndex - 1 + numWeapons) % numWeapons
	}
	return []Action{}
}
