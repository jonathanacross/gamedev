package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
)

type ChestState int

const (
	ChestClosed ChestState = iota
	ChestOpening
	ChestOpen

	ChestClosedIdx = 0
	ChestOpenIdx   = 3
)

// An openable chest.  Implements Interactable.
type Chest struct {
	core.BasePhysical
	State         ChestState
	Contents      core.UpgradeType
	spriteSheet   *core.SpriteSheet
	openAnimation *core.Animation
	releasedItem  bool
}

func NewChest(location core.Location, upgradeType core.UpgradeType) *Chest {
	spriteSheet := core.NewSpriteSheet(core.TileSize, core.TileSize, 4, 2)
	openAnimation := core.NewAnimation([]int{1, 2}, 20, false)

	return &Chest{
		BasePhysical: core.BasePhysical{
			BaseSprite: core.BaseSprite{
				Loc:     location,
				Image:   assets.ChestSpritesImage,
				SrcRect: spriteSheet.Rect(ChestClosedIdx),
				DrawOffset: core.Location{
					X: 7,
					Y: 9,
				},
			},
			PushBoxOffset: core.Rect{
				Left:   -8,
				Top:    -8,
				Right:  core.TileSize - 8,
				Bottom: core.TileSize - 8,
			},
		},
		State:         ChestClosed,
		Contents:      upgradeType,
		spriteSheet:   spriteSheet,
		openAnimation: openAnimation,
		releasedItem:  false,
	}
}

func (c *Chest) Touch(level core.Level, p core.Player) []core.Action {
	return []core.Action{}
}

func (c *Chest) Interact(level core.Level, p core.Player) []core.Action {
	if c.State == ChestClosed {
		c.State = ChestOpening
	}
	return []core.Action{}
}

func (c *Chest) Update(level core.Level, p core.Player) core.UpdateResult {
	switch c.State {
	case ChestClosed:
		c.SrcRect = c.spriteSheet.Rect(ChestClosedIdx)

	case ChestOpening:
		c.openAnimation.Update()
		c.SrcRect = c.spriteSheet.Rect(c.openAnimation.Frame())
		if c.openAnimation.IsFinished() {
			c.State = ChestOpen
		}
		if !c.releasedItem {
			upgradeType := c.Contents
			firstFound := p.GetCurrentStats()[upgradeType] == 0
			c.releasedItem = true
			return core.UpdateResult{Actions: []core.Action{{
				Type:        core.ActionShowChestItem,
				Location:    c.Loc,
				UpgradeType: upgradeType,
				FirstFound:  firstFound,
			}}}
		}
	case ChestOpen:
		c.SrcRect = c.spriteSheet.Rect(ChestOpenIdx)
	}

	return core.UpdateResult{}
}

func (c *Chest) CanRemove() bool {
	return false
}
