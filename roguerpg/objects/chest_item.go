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

func NewChestItem(upgradeType core.UpgradeType, location core.Location) *ChestItem {
	spriteSheet := core.NewSpriteSheet(16, 16, 8, 3)
	typeToIndex := map[core.UpgradeType]int{
		core.UpgradeTypeHeart:     8,
		core.UpgradeTypeSword:     9,
		core.UpgradeTypeBoomerang: 10,
		core.UpgradeTypeShield:    11,
		core.UpgradeTypeBomb:      12,
		core.UpgradeTypeBow:       13,
		core.UpgradeTypeWand:      14,
	}
	startIdx := typeToIndex[upgradeType]

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
