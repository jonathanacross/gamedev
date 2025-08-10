package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayerView struct {
	previousPlayer *Button
	nextPlayer     *Button
	playerName     string
	X              int
	Y              int
	index          int
	activeColor    color.Color
}

type namedEngine struct {
	name   string
	engine Engine
}

var namedEngines = []namedEngine{
	{"Human", &HumanEngine{}},
	{"Random", &RandomEngine{}},
	{"Greedy", &GreedyEngine{}},
	{"Minimax2", &MinimaxEngine{maxDepth: 2}},
	{"Minimax3", &MinimaxEngine{maxDepth: 3}},
	{"Minimax4", &MinimaxEngine{maxDepth: 4}},
}

const (
	pvSpacing    = 15
	pvLabelWidth = 55
	pvButtonSize = 32
)

func NewPlayerView(x, y int, index int, activeColor color.Color) *PlayerView {
	ps := &PlayerView{
		previousPlayer: nil,
		nextPlayer:     nil,
		playerName:     namedEngines[index].name,
		X:              x,
		Y:              y,
		index:          index,
		activeColor:    activeColor,
	}
	previousPlayer := NewButton(x+pvSpacing, y+pvSpacing, pvButtonSize, pvButtonSize,
		LeftArrowIdleImage,
		LeftArrowPressedImage,
		func() { ps.SwitchPlayer(-1) })
	nextPlayer := NewButton(x+pvButtonSize+pvLabelWidth+3*pvSpacing, y+pvSpacing, pvButtonSize, pvButtonSize,
		RightArrowIdleImage,
		RightArrowPressedImage,
		func() { ps.SwitchPlayer(1) })
	ps.previousPlayer = previousPlayer
	ps.nextPlayer = nextPlayer
	return ps
}

func (p *PlayerView) Update() {
	p.previousPlayer.Update()
	p.nextPlayer.Update()
}

func (p *PlayerView) Draw(screen *ebiten.Image, score int, active bool) {
	p.previousPlayer.Draw(screen)
	p.nextPlayer.Draw(screen)

	pvHeight := 2*pvButtonSize + pvSpacing
	pvWidth := pvLabelWidth + 4*pvSpacing + 2*pvButtonSize
	scoreX := p.X + pvWidth/2
	scoreY := p.Y + pvSpacing + pvButtonSize
	labelX := p.X + pvWidth/2
	labelY := p.Y + pvSpacing + 5

	op := &text.DrawOptions{}
	fontSize := float64(16)
	op.GeoM.Translate(float64(labelX), float64(labelY))
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, p.playerName, &text.GoTextFace{
		Source: DisplayFont,
		Size:   fontSize,
	}, op)

	op.GeoM.Reset()
	op.GeoM.Translate(float64(scoreX), float64(scoreY))
	text.Draw(screen, fmt.Sprintf("%d", score), &text.GoTextFace{
		Source: DisplayFont,
		Size:   fontSize,
	}, op)

	strokeWidth := float32(1.0)
	if active {
		strokeWidth = float32(3.0)
	}
	vector.StrokeRect(
		screen, float32(p.X), float32(p.Y),
		float32(pvWidth), float32(pvHeight), strokeWidth, p.activeColor, false)
}

func (p *PlayerView) SwitchPlayer(offset int) {
	p.index = (p.index + offset + len(namedEngines)) % len(namedEngines)
	p.playerName = namedEngines[p.index].name
}

func (p *PlayerView) GetEngine() Engine {
	return namedEngines[p.index].engine
}
