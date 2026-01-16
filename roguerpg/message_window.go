package main

import (
	"image/color"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	messageWindowWidth  = 200
	messageWindowHeight = 40
)

// Implements GameState
type MessageWindow struct {
	Window  *core.Window
	message string
}

func NewMessageWindow(message string) *MessageWindow {
	windowX := float64(ScreenWidth-messageWindowWidth) / 2.0
	windowY := 0.7*float64(ScreenHeight) - 0.5*float64(messageWindowHeight)

	return &MessageWindow{
		Window: core.NewWindow(core.Rect{
			Left:   windowX,
			Top:    windowY,
			Right:  windowX + messageWindowWidth,
			Bottom: windowY + messageWindowHeight,
		}),
		message: message,
	}
}

func (w *MessageWindow) Draw(screen *ebiten.Image, context *core.GameContext) {
	w.Window.Draw(screen)

	titleLoc := core.Location{X: w.Window.Rect.Left + 8, Y: w.Window.Rect.Top + 8}
	drawTextAt(screen, w.message, titleLoc.X, titleLoc.Y, text.AlignStart, color.White)
}

func (w *MessageWindow) Update(context *core.GameContext) []core.Action {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyTab) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return []core.Action{
			{Type: core.ActionPopState},
		}
	}
	return []core.Action{}
}
