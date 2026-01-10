package main

import (
	"fmt"
	"image/color"
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	windowWidth        = 200
	windowHeight       = 200
	buttonWidth        = 140
	buttonHeight       = 19
	buttonSpacing      = 20
	buttonGroupXOffset = 10
	buttonGroupYOffset = 24
	buttonTextOffsetY  = 2
	buttonIconOffsetY  = 1
)

type Upgrade struct {
	UpgradeType core.UpgradeType
	WeaponLevel int
	NeededExp   int
}

// Implements GameState
type UpgradeSelector struct {
	Window   *core.Window
	upgrades []Upgrade
	buttons  []*UpgradeButton
	currIdx  int
}

var UpgradeSelectorInstance *UpgradeSelector = NewUpgradeSelector()

func NewUpgradeSelector() *UpgradeSelector {
	windowX := float64(ScreenWidth-windowWidth) / 2.0
	windowY := float64(ScreenHeight-windowHeight) / 2.0
	currIdx := 0
	// TODO: put these in same order as in weapon select
	upgradeTypes := []core.UpgradeType{
		core.UpgradeTypeNone,
		core.UpgradeTypeHeart,
		core.UpgradeTypeSword,
		core.UpgradeTypeBoomerang,
		core.UpgradeTypeBow,
		core.UpgradeTypeShield,
		core.UpgradeTypeBomb,
		core.UpgradeTypeWand,
	}
	upgrades := []Upgrade{}
	for i, upgradeType := range upgradeTypes {
		upgrades = append(upgrades, Upgrade{UpgradeType: upgradeType, WeaponLevel: 1, NeededExp: 10 + 2*i})
	}
	buttons := []*UpgradeButton{}
	for i, u := range upgrades {
		loc := core.Location{X: windowX + buttonGroupXOffset, Y: windowY + float64(buttonGroupYOffset+i*buttonSpacing)}
		availableExp := 15
		maxLevel := 3
		hasUpgradeMaterials := true
		button := NewUpgradeButton(
			loc,
			u.UpgradeType,
			u.WeaponLevel,
			maxLevel,
			hasUpgradeMaterials,
			u.NeededExp,
			availableExp,
		)
		buttons = append(buttons, button)
	}

	return &UpgradeSelector{
		Window:   core.NewWindow(core.Rect{Left: windowX, Top: windowY, Right: windowX + windowWidth, Bottom: windowY + windowHeight}),
		upgrades: upgrades,
		buttons:  buttons,
		currIdx:  currIdx,
	}
}

func (w *UpgradeSelector) Draw(screen *ebiten.Image, context *core.GameContext) {
	w.Window.Draw(screen)

	titleLoc := core.Location{X: w.Window.Rect.Left + 50, Y: w.Window.Rect.Top + 5}
	textColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	drawTextAt(screen, "Choose an upgrade", titleLoc.X, titleLoc.Y, text.AlignStart, textColor)

	for _, button := range w.buttons {
		button.Draw(screen)
	}
}

func (w *UpgradeSelector) Update(context *core.GameContext) []core.Action {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return []core.Action{
			{Type: core.ActionPopState},
		}
	}
	numUpgrades := len(w.upgrades)
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		w.currIdx = (w.currIdx + numUpgrades - 1) % numUpgrades
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		w.currIdx = (w.currIdx + 1) % numUpgrades
	}

	// TODO: this is a bit hacky.  Should have function to change button state
	for _, button := range w.buttons {
		button.buttonBackground.Status = core.ButtonEnabled
	}
	w.buttons[w.currIdx].buttonBackground.Status = core.ButtonHovered
	return []core.Action{}
}

type UpgradeButton struct {
	buttonBackground    *core.Button
	location            core.Location
	isDisabled          bool
	upgradeType         core.UpgradeType
	nextLevel           int
	maxLevel            int
	hasUpgradeMaterials bool
	neededExp           int
	availableExp        int
	weaponImage         *ebiten.Image
	upgradeImage        *ebiten.Image
}

func NewUpgradeButton(
	location core.Location, upgradeType core.UpgradeType,
	nextLevel int, maxLevel int,
	hasUpgradeMaterials bool,
	neededExp int, availableExp int) *UpgradeButton {

	isDisabled := (upgradeType != core.UpgradeTypeNone) && (!hasUpgradeMaterials || neededExp > availableExp)
	weaponImage := getIconImage(upgradeType, false, false)
	upgradeImage := getIconImage(upgradeType, true, !hasUpgradeMaterials)
	return &UpgradeButton{
		buttonBackground:    core.NewButton(core.Rect{Left: location.X, Top: location.Y, Right: location.X + buttonWidth, Bottom: location.Y + buttonHeight}),
		location:            location,
		isDisabled:          isDisabled,
		upgradeType:         upgradeType,
		nextLevel:           nextLevel,
		maxLevel:            maxLevel,
		hasUpgradeMaterials: hasUpgradeMaterials,
		neededExp:           neededExp,
		availableExp:        availableExp,
		weaponImage:         weaponImage,
		upgradeImage:        upgradeImage,
	}
}

func getIconIndex(upgradeType core.UpgradeType, isUpgrade bool, grayed bool) int {
	offset := 0
	if isUpgrade {
		offset += uiIconTilesetWidth
	}
	if grayed {
		offset += uiIconTilesetWidth
	}

	switch upgradeType {
	case core.UpgradeTypeHeart:
		return UiIconHeart + offset
	case core.UpgradeTypeSword:
		return UiIconSword + offset
	case core.UpgradeTypeShield:
		return UiIconShield + offset
	case core.UpgradeTypeBow:
		return UiIconBow + offset
	case core.UpgradeTypeBoomerang:
		return UiIconBoomerang + offset
	case core.UpgradeTypeBomb:
		return UiIconBomb + offset
	case core.UpgradeTypeWand:
		return UiIconWand + offset
	default:
		return UiIconEmpty
	}
}

func getIconImage(upgradeType core.UpgradeType, isUpgrade bool, isGrayed bool) *ebiten.Image {
	iconIndex := getIconIndex(upgradeType, isUpgrade, isGrayed)
	spriteSheet := core.NewSpriteSheet(16, 16, 8, 1)
	icon := assets.UiIconsImage.SubImage(spriteSheet.Rect(iconIndex)).(*ebiten.Image)
	return icon
}

func drawIcon(screen *ebiten.Image, location core.Location, icon *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(location.X, location.Y)
	screen.DrawImage(icon, op)
}

func (b *UpgradeButton) Draw(screen *ebiten.Image) {
	b.buttonBackground.Draw(screen)
	textColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	grayedTextColor := color.RGBA{R: 128, G: 128, B: 128, A: 255}

	if b.upgradeType == core.UpgradeTypeNone {
		drawTextAt(screen, "no upgrade", b.location.X+18, b.location.Y+buttonTextOffsetY, text.AlignStart, textColor)
		return
	}

	weaponBoxLoc := core.Location{X: b.location.X + 1, Y: b.location.Y + buttonIconOffsetY}
	drawIcon(screen, weaponBoxLoc, b.weaponImage)

	levString := fmt.Sprintf("L%d/%d:", b.nextLevel, b.maxLevel)
	drawTextAt(screen, levString, b.location.X+18, b.location.Y+buttonTextOffsetY, text.AlignStart, textColor)

	upgradeBoxLoc := core.Location{X: b.location.X + 48, Y: b.location.Y + buttonIconOffsetY}
	drawIcon(screen, upgradeBoxLoc, b.upgradeImage)

	drawTextAt(screen, "+", b.location.X+64, b.location.Y+buttonTextOffsetY, text.AlignStart, textColor)
	expString := fmt.Sprintf("Exp: %d", b.neededExp)
	if b.availableExp < b.neededExp {
		textColor = grayedTextColor
	}
	drawTextAt(screen, expString, b.location.X+76, b.location.Y+buttonTextOffsetY, text.AlignStart, textColor)
}
