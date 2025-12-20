package main

import (
	"roguerpg/assets"
	"roguerpg/core"
)

const (
	UpgradeTypeHeartIdx = iota
	UpgradeTypeSwordIdx
	UpgradeTypeBoomerangIdx
	UpgradeTypeShieldIdx
	UpgradeTypeBombIdx
	UpgradeTypeBowIdx
	UpgradeTypeWandIdx
)

const uiIconTilesetWidth = 8
const uiIconTilesetHeight = 3

type UpgradeIcon struct {
	core.BaseSprite
	UpgradeType core.UpgradeType
	spriteSheet *core.SpriteSheet
}

func getUpgradeSrcIndex(upgradeType core.UpgradeType, grayed bool) int {
	typeToIndex := map[core.UpgradeType]int{
		core.UpgradeTypeHeart:     UpgradeTypeHeartIdx,
		core.UpgradeTypeSword:     UpgradeTypeSwordIdx,
		core.UpgradeTypeBoomerang: UpgradeTypeBoomerangIdx,
		core.UpgradeTypeShield:    UpgradeTypeShieldIdx,
		core.UpgradeTypeBomb:      UpgradeTypeBombIdx,
		core.UpgradeTypeBow:       UpgradeTypeBowIdx,
		core.UpgradeTypeWand:      UpgradeTypeWandIdx,
	}
	startIdx := uiIconTilesetWidth
	if grayed {
		startIdx = 2 * uiIconTilesetWidth
	}

	return startIdx + typeToIndex[upgradeType]
}

func NewUpgradeIcon(upgradeType core.UpgradeType, loc core.Location) *UpgradeIcon {
	spriteSheet := core.NewSpriteSheet(16, 16, uiIconTilesetWidth, uiIconTilesetHeight)
	srcIndex := getUpgradeSrcIndex(upgradeType, false)

	return &UpgradeIcon{
		BaseSprite: core.BaseSprite{
			Loc:     loc,
			Image:   assets.UiIconsImage,
			SrcRect: spriteSheet.Rect(srcIndex),
			DrawOffset: core.Location{
				X: 8,
				Y: 8,
			},
		},
	}
}
