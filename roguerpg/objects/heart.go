package objects

import (
	"roguerpg/assets"
	"roguerpg/core"
	"time"
)

const (
	HeartDuration = 6 * time.Second
)

// Implements GameObject
type Heart struct {
	core.BasePhysical
	spriteSheet *core.SpriteSheet
	animation   *core.Animation
	collected   bool
	timer       *core.Timer
}

func NewHeart(location core.Location) *Heart {
	animation := core.NewAnimation([]int{2, 3, 2, 1}, 16, true)
	spriteSheet := core.NewSpriteSheet(9, 8, 6, 1)

	return &Heart{
		BasePhysical: core.BasePhysical{
			BaseSprite: core.BaseSprite{
				Loc:     location,
				Image:   assets.HeartSpritesImage,
				SrcRect: spriteSheet.Rect(0),
				DrawOffset: core.Location{
					X: 5,
					Y: 4,
				},
			},
			PushBoxOffset: core.Rect{Left: -5, Top: -4, Right: 4, Bottom: 4},
		},
		collected:   false,
		animation:   animation,
		spriteSheet: spriteSheet,
		timer:       core.NewTimer(HeartDuration),
	}
}

func (h *Heart) Touch(_ core.Level, p core.Player) []core.Action {
	if !h.collected {
		h.collected = true
		p.AddHealth(2)
	}
	return []core.Action{}
}

func (h *Heart) Interact(_ core.Level, _ core.Player) []core.Action {
	return []core.Action{}
}

func (h *Heart) Update(_ core.Level, _ core.Player) core.UpdateResult {
	h.animation.Update()
	h.SrcRect = h.spriteSheet.Rect(h.animation.Frame())
	h.timer.Update()
	return core.UpdateResult{}
}

func (h *Heart) CanRemove() bool {
	return h.timer.IsReady() || h.collected
}
