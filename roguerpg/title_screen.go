package main

import (
	"image/color"
	"roguerpg/assets"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TitleScreenState implements core.GameState
type TitleScreenState struct {
}

func (t *TitleScreenState) Update(context *core.GameContext) []core.Action {
	if len(inpututil.AppendJustPressedKeys(nil)) > 0 {
		return []core.Action{
			{
				Type: core.ActionPushState,
				GameState: NewFadeOutState([]core.Action{
					{Type: core.ActionStartGame}, // Clears the stack and adds MainGameState
					{Type: core.ActionPushState, GameState: NewFadeInState()},
				}),
			},
		}
	}
	return nil
}

func (t *TitleScreenState) Draw(screen *ebiten.Image, context *core.GameContext) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(assets.TitleScreenImage, op)

	msg := "Press Any Key"
	x := float64(ScreenWidth) / 2.0
	y := float64(ScreenHeight) - 40.0

	drawTextAt(screen, msg, x, y, text.AlignCenter, color.White)
}
