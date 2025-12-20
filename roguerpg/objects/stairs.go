package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
)

// Stairs that lead up or down between levels.  Implements Interactable.
type Stairs struct {
	core.BasePhysical

	IsUpstairs bool
}

func NewStairs(location core.Location, isUpstairs bool) *Stairs {
	spriteSheet := core.NewSpriteSheet(core.TileSize, core.TileSize, 4, 2)
	spriteIdx := 0
	if !isUpstairs {
		spriteIdx = 1
	}

	s := &Stairs{
		BasePhysical: core.BasePhysical{
			BaseSprite: core.BaseSprite{
				Loc:     location,
				Image:   assets.StairsSpritesImage,
				SrcRect: spriteSheet.Rect(spriteIdx),
				DrawOffset: core.Location{
					X: 8,
					Y: 8,
				},
			},
			PushBoxOffset: core.Rect{
				Left:   -8,
				Top:    -8,
				Right:  core.TileSize - 8,
				Bottom: core.TileSize - 8,
			},
		},
		IsUpstairs: isUpstairs,
	}
	return s
}

func (s *Stairs) Touch(level core.Level, p core.Player) []core.Action {
	return []core.Action{}
}

func (s *Stairs) Interact(level core.Level, p core.Player) []core.Action {
	if s.IsUpstairs && level != nil {
		return []core.Action{{Type: core.ActionGoUpLevel}}
	} else if !s.IsUpstairs {
		return []core.Action{{Type: core.ActionGoDownLevel}}
	}

	return []core.Action{}
}

func (s *Stairs) Update(level core.Level, p core.Player) core.UpdateResult {
	return core.UpdateResult{}
}

func (s *Stairs) CanRemove() bool {
	return false
}
