package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
	"time"
)

const (
	riseTime  = 2000 * time.Millisecond
	riseSpeed = 0.10
)

// Implements GameObject
type ChestItem struct {
	core.BaseSprite
	spriteSheet *core.SpriteSheet
	timer       *core.Timer
}

func NewChestItem(upgradeType core.UpgradeType, location core.Location, firstFound bool) *ChestItem {
	spriteSheet := core.NewSpriteSheet(16, 16, 8, 3)
	typeToIndex := map[core.UpgradeType]int{
		core.UpgradeTypeHeart:     0,
		core.UpgradeTypeSword:     1,
		core.UpgradeTypeBoomerang: 2,
		core.UpgradeTypeShield:    3,
		core.UpgradeTypeBomb:      4,
		core.UpgradeTypeBow:       5,
		core.UpgradeTypeWand:      6,
	}
	startIdx := typeToIndex[upgradeType]
	if !firstFound {
		// If this isn't the first time we've seen this upgrade type,
		// offset to show the upgrade version of the icon.
		startIdx += 8
	}
	return &ChestItem{
		BaseSprite: core.BaseSprite{
			Loc:     location,
			Image:   assets.UiIconsImage,
			SrcRect: spriteSheet.Rect(startIdx),
			DrawOffset: core.Location{
				X: 8,
				Y: 8,
			},
		},
		spriteSheet: spriteSheet,
		timer:       core.NewTimer(riseTime),
	}
}

func (b *ChestItem) Update(_ core.Level, _ core.Player) core.UpdateResult {
	b.timer.Update()
	b.Loc.Y -= riseSpeed
	return core.UpdateResult{}
}

func (b *ChestItem) CanRemove() bool {
	return b.timer.IsReady()
}
