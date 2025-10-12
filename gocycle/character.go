package main

import (
	"gocycle/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ControllerType int

const (
	HumanFirstPlayer ControllerType = iota
	HumanSecondPlayer
	ComputerPlayer
)

type CharData struct {
	ID    int
	Name  string
	Image *ebiten.Image
	// TODO: rename these to BrightColor, DarkColor
	SelectedColor  color.Color
	FrameColor     color.Color
	Controller     core.PlayerController
	ControllerType ControllerType
}

var Characters []CharData = loadCharData()
var NumCharacters = len(Characters)

func loadCharData() []CharData {
	return []CharData{
		{
			ID:             1,
			Name:           "Milo",
			Image:          MiloCharImage,
			SelectedColor:  color.RGBA{29, 70, 125, 255},
			FrameColor:     color.RGBA{3, 166, 224, 255},
			Controller:     &core.RandomAvoidingController{},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             2,
			Name:           "Sara",
			Image:          SaraCharImage,
			SelectedColor:  color.RGBA{167, 151, 50, 255},
			FrameColor:     color.RGBA{248, 243, 79, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.10},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             3,
			Name:           "Dr. Q",
			Image:          DrQCharImage,
			SelectedColor:  color.RGBA{23, 110, 114, 255},
			FrameColor:     color.RGBA{74, 199, 198, 255},
			Controller:     &core.WallHuggerController{},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             4,
			Name:           "Erica",
			Image:          EricaCharImage,
			SelectedColor:  color.RGBA{156, 20, 38, 255},
			FrameColor:     color.RGBA{231, 64, 71, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.005},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             5,
			Name:           "Biff",
			Image:          BiffCharImage,
			SelectedColor:  color.RGBA{182, 70, 37, 255},
			FrameColor:     color.RGBA{238, 156, 50, 255},
			Controller:     &core.AreaController{},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             6,
			Name:           "Elara",
			Image:          ElaraCharImage,
			SelectedColor:  color.RGBA{67, 67, 130, 255},
			FrameColor:     color.RGBA{121, 121, 203, 255},
			Controller:     &core.MinimaxAreaController{MaxDepth: 3},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             7,
			Name:           "Mike Green",
			Image:          MikeGCharImage,
			SelectedColor:  color.RGBA{20, 104, 20, 255},
			FrameColor:     color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
		},
		{
			ID:             8,
			Name:           "Mike Violet",
			Image:          MikeVCharImage,
			SelectedColor:  color.RGBA{94, 33, 72, 255},
			FrameColor:     color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
		},
		{
			ID:             9,
			Name:           "Heather Green",
			Image:          HeatherGCharImage,
			SelectedColor:  color.RGBA{20, 104, 20, 255},
			FrameColor:     color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
		},
		{
			ID:             10,
			Name:           "Heather Violet",
			Image:          HeatherVCharImage,
			SelectedColor:  color.RGBA{94, 33, 72, 255},
			FrameColor:     color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
		},
	}
}
