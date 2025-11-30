package main

import "github.com/hajimehoshi/ebiten/v2"

// Implements GameState
type WeaponSelector struct {
	windowImage      *ebiten.Image
	iconsImage       *ebiten.Image
	iconsSpriteSheet *SpriteSheet
	windowLoc        Location
}

// TODO: make this a singleton
func NewWeaponSelector() *WeaponSelector {
	windowX := float64(ScreenWidth-WeaponSelectWindowImage.Bounds().Dx()) / 2.0
	windowY := float64(ScreenHeight-WeaponSelectWindowImage.Bounds().Dy()) / 2.0
	return &WeaponSelector{
		windowImage:      WeaponSelectWindowImage,
		iconsImage:       UiIconsImage,
		iconsSpriteSheet: NewSpriteSheet(16, 16, 8, 1),
		windowLoc:        Location{X: windowX, Y: windowY},
	}
}

func (w *WeaponSelector) Draw(screen *ebiten.Image, context *GameContext) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.windowLoc.X, w.windowLoc.Y)

	screen.DrawImage(w.windowImage, op)

	// op.GeoM.Concat(cameraMatrix)
	// currImage := bs.image.SubImage(bs.srcRect).(*ebiten.Image)
	// screen.DrawImage(currImage, op)
}

func (w *WeaponSelector) Update(context *GameContext) []Action {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return []Action{
			Action{
				Type:         ActionPopState,
				GameState:    nil,
				Location:     Location{},
				Direction:    ZeroVector(),
				DamageSource: nil,
			},
		}
	}
	return []Action{}
}
