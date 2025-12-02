package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Stairs struct {
	BasePhysical

	IsUpstairs bool
}

func NewStairs(location Location, isUpstairs bool) *Stairs {
	spriteSheet := NewSpriteSheet(TileSize, TileSize, 4, 2)
	spriteIdx := 0
	if !isUpstairs {
		spriteIdx = 1
	}

	s := &Stairs{
		BasePhysical: BasePhysical{
			BaseSprite: BaseSprite{
				Location: location,
				image:    DungeonObjectsTileset,
				srcRect:  spriteSheet.Rect(spriteIdx),
				drawOffset: Location{
					X: 8,
					Y: 8,
				},
			},
			pushBoxOffset: Rect{
				Left:   -8,
				Top:    -8,
				Right:  TileSize - 8,
				Bottom: TileSize - 8,
			},
		},
		IsUpstairs: isUpstairs,
	}
	return s
}

func (s *Stairs) Update(level *Level, p *Player) UpdateResult {
	// Check for player interaction
	if p.GetPushBox().Intersects(s.GetPushBox()) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if s.IsUpstairs && level.UpLevel != nil {
			return UpdateResult{Actions: []Action{{Type: ActionGoUpLevel}}}
		} else if !s.IsUpstairs && level.DownLevel != nil {
			return UpdateResult{Actions: []Action{{Type: ActionGoDownLevel}}}
		}
	}
	return UpdateResult{}
}

func (s *Stairs) CanRemove() bool {
	return false
}

type ChestState int

const (
	ChestClosed ChestState = iota
	ChestOpening
	ChestOpen

	ChestClosedIdx = 4
	ChestOpenIdx   = 7
)

type Chest struct {
	BasePhysical
	State ChestState

	spriteSheet   *SpriteSheet
	openAnimation *Animation
}

func NewChest(location Location) *Chest {
	spriteSheet := NewSpriteSheet(TileSize, TileSize, 4, 2)
	openAnimation := NewAnimation([]int{5, 6}, 20, false)

	return &Chest{
		BasePhysical: BasePhysical{
			BaseSprite: BaseSprite{
				Location: location,
				image:    DungeonObjectsTileset,
				srcRect:  spriteSheet.Rect(ChestClosedIdx),
				drawOffset: Location{
					X: 7,
					Y: 9,
				},
			},
			pushBoxOffset: Rect{
				Left:   -8,
				Top:    -8,
				Right:  TileSize - 8,
				Bottom: TileSize - 8,
			},
		},
		State:         ChestClosed,
		spriteSheet:   spriteSheet,
		openAnimation: openAnimation,
	}
}

func (c *Chest) Update(level *Level, p *Player) UpdateResult {
	switch c.State {
	case ChestClosed:
		c.srcRect = c.spriteSheet.Rect(ChestClosedIdx)
		// Check for player interaction
		if p.GetPushBox().Intersects(c.GetPushBox()) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			c.State = ChestOpening
		}

	case ChestOpening:
		c.openAnimation.Update()
		c.srcRect = c.spriteSheet.Rect(c.openAnimation.Frame())
		if c.openAnimation.IsFinished() {
			c.State = ChestOpen
		}
	case ChestOpen:
		c.srcRect = c.spriteSheet.Rect(ChestOpenIdx)
	}

	return UpdateResult{}
}

func (c *Chest) CanRemove() bool {
	return false
}
