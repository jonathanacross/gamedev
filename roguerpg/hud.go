package main

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"image/color"
	"roguerpg/assets"
	"roguerpg/core"
)

const (
	HeartWidth        = 11
	HeartHeight       = 10
	HeartSubdivisions = 2
)

var heartInstance = newHeart()

type Heart struct {
	image       *ebiten.Image
	spriteSheet *core.SpriteSheet
}

func newHeart() Heart {
	return Heart{
		image:       assets.UiHealthHeartImage,
		spriteSheet: core.NewSpriteSheet(HeartWidth, HeartHeight, HeartSubdivisions+1, 1),
	}
}

func DrawHeart(screen *ebiten.Image, x float64, y float64, frame int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	subRect := heartInstance.spriteSheet.Rect(frame)
	currImage := heartInstance.image.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

func DrawPlayerHeath(screen *ebiten.Image, currHeath int, maxHeath int) {
	numHearts := (maxHeath + (HeartSubdivisions - 1)) / HeartSubdivisions

	for i := range numHearts {
		x := float64(20 + i*HeartWidth)
		y := float64(20)
		frame := core.Clamp(currHeath-HeartSubdivisions*i, 0, HeartSubdivisions)
		DrawHeart(screen, x, y, frame)
	}
}

func DrawCurrentWeapons(screen *ebiten.Image,
	primaryWeapon core.WeaponType,
	primaryWeaponProgress float64,
	secondaryWeapon core.WeaponType,
	secondaryWeaponProgress float64) {
	primaryWeaponBox := NewWeaponBox(core.Location{X: 100, Y: 14}, primaryWeapon, true)
	secondaryWeaponBox := NewWeaponBox(core.Location{X: 135, Y: 14}, secondaryWeapon, true)

	primaryWeaponBox.ReadyPercent = primaryWeaponProgress
	secondaryWeaponBox.ReadyPercent = secondaryWeaponProgress

	primaryWeaponBox.Draw(screen)
	secondaryWeaponBox.Draw(screen)
}

// TODO: move this to a utils file
func drawTextAt(screen *ebiten.Image, message string, x float64, y float64, align text.Align, c color.Color) {
	fontSize := float64(16)
	fontFace := &text.GoTextFace{
		Source: assets.TextFaceSource,
		Size:   fontSize,
	}

	// Manually handle alignment to ensure pixel-perfect rendering
	textWidth, _ := text.Measure(message, fontFace, 1.0)
	if align == text.AlignCenter {
		x -= float64(textWidth) / 2
	} else if align == text.AlignEnd {
		x -= float64(textWidth)
	}
	x = float64(int(x))
	y = float64(int(y))

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignStart

	text.Draw(screen, message, fontFace, op)
}

func DrawExperience(screen *ebiten.Image, experience int) {
	expString := strconv.Itoa(experience)
	drawTextAt(screen, "Exp: "+expString, ScreenWidth-80, 18, text.AlignStart, color.White)
}

func DrawHeadsUpDisplay(screen *ebiten.Image, player core.Player) {
	DrawPlayerHeath(screen, player.GetHealth(), player.GetMaxHealth())
	DrawExperience(screen, player.GetExperience())
	primaryWeapon := player.GetPrimaryWeapon()
	secondaryWeapon := player.GetSecondaryWeapon()
	DrawCurrentWeapons(screen,
		primaryWeapon, player.GetWeaponProgress(primaryWeapon),
		secondaryWeapon, player.GetWeaponProgress(secondaryWeapon))
}
