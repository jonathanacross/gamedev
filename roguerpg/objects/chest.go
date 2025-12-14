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

	ChestClosedIdx = 4
	ChestOpenIdx   = 7
)

// An openable chest.  Implements Interactable.
type Chest struct {
	core.BasePhysical
	State ChestState

	spriteSheet   *core.SpriteSheet
	openAnimation *core.Animation
}

func NewChest(location core.Location) *Chest {
	spriteSheet := core.NewSpriteSheet(core.TileSize, core.TileSize, 4, 2)
	openAnimation := core.NewAnimation([]int{5, 6}, 20, false)

	return &Chest{
		BasePhysical: core.BasePhysical{
			BaseSprite: core.BaseSprite{
				Loc:     location,
				Image:   assets.DungeonObjectsTileset,
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
		spriteSheet:   spriteSheet,
		openAnimation: openAnimation,
	}
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
	case ChestOpen:
		c.SrcRect = c.spriteSheet.Rect(ChestOpenIdx)
	}

	return core.UpdateResult{}
}

func (c *Chest) CanRemove() bool {
	return false
}
