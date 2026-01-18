package main

import (
	"image/color"
	"roguerpg/core"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	fadeDuration = 500 * time.Millisecond
)

var (
	// Shared 1x1 white image for fading
	fadeImage = ebiten.NewImage(1, 1)
)

func init() {
	fadeImage.Fill(color.White)
}

// Implements GameState
type FadeOutState struct {
	timer             *core.Timer
	onCompleteActions []core.Action
}

func NewFadeOutState(nextActions []core.Action) *FadeOutState {
	return &FadeOutState{
		timer:             core.NewTimer(fadeDuration),
		onCompleteActions: nextActions,
	}
}

func (s *FadeOutState) Update(ctx *core.GameContext) []core.Action {
	s.timer.Update()
	if s.timer.IsReady() {
		// When done, return the sequence: Pop myself, then do the "next" things
		return append([]core.Action{{Type: core.ActionPopState}}, s.onCompleteActions...)
	}
	return nil
}

func (s *FadeOutState) Draw(screen *ebiten.Image, ctx *core.GameContext) {
	progress := s.timer.GetProgress()
	alpha := float32(progress)

	op := &ebiten.DrawImageOptions{}
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	op.GeoM.Scale(float64(w), float64(h))

	// Fade to black: set color to (0, 0, 0) and alpha to progress
	op.ColorScale.Scale(0, 0, 0, alpha)

	screen.DrawImage(fadeImage, op)
}

// FadeInState: Fades from black to transparent
// Implements GameState
type FadeInState struct {
	timer *core.Timer
}

func NewFadeInState() *FadeInState {
	return &FadeInState{timer: core.NewTimer(fadeDuration)}
}

func (s *FadeInState) Update(ctx *core.GameContext) []core.Action {
	s.timer.Update()
	if s.timer.IsReady() {
		return []core.Action{{Type: core.ActionPopState}}
	}
	return nil
}

func (s *FadeInState) Draw(screen *ebiten.Image, ctx *core.GameContext) {
	progress := s.timer.GetProgress()
	alpha := float32(1.0 - progress)

	op := &ebiten.DrawImageOptions{}
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	op.GeoM.Scale(float64(w), float64(h))

	// Fade from black: set color to (0, 0, 0) and alpha to (1-progress)
	op.ColorScale.Scale(0, 0, 0, alpha)

	screen.DrawImage(fadeImage, op)
}
