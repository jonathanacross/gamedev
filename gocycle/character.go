package main

import (
	"gocycle/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Cast of characters
const (
	CharacterMilo int = iota
	CharacterSara
	CharacterDrQ
	CharacterErica
	CharacterBiff
	CharacterElara
	CharacterMikeG
	CharacterMikeV
	CharacterHeatherG
	CharacterHeatherV
	NumCharacters
)

type ControllerType int

const (
	HumanFirstPlayer ControllerType = iota
	HumanSecondPlayer
	ComputerPlayer
)

type CharData struct {
	Name  string
	Image *ebiten.Image
	// TODO: rename these to BrightColor, DarkColor
	SelectedColor  color.Color
	FrameColor     color.Color
	Controller     core.PlayerController
	ControllerType ControllerType
}

var Characters []CharData = loadCharData()

func loadCharData() []CharData {
	return []CharData{
		{
			Name:           "Milo",
			Image:          MiloCharImage,
			SelectedColor:  color.RGBA{22, 48, 83, 255},
			FrameColor:     color.RGBA{3, 166, 224, 255},
			Controller:     &core.RandomAvoidingController{},
			ControllerType: ComputerPlayer,
		},
		{
			Name:           "Sara",
			Image:          SaraCharImage,
			SelectedColor:  color.RGBA{146, 132, 51, 255},
			FrameColor:     color.RGBA{248, 243, 79, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.30},
			ControllerType: ComputerPlayer,
		},
		{
			Name:           "Dr. Q",
			Image:          DrQCharImage,
			SelectedColor:  color.RGBA{20, 75, 78, 255},
			FrameColor:     color.RGBA{74, 199, 198, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.10},
			ControllerType: ComputerPlayer,
		},
		{
			Name:           "Erica",
			Image:          EricaCharImage,
			SelectedColor:  color.RGBA{104, 9, 13, 255},
			FrameColor:     color.RGBA{231, 64, 71, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.01},
			ControllerType: ComputerPlayer,
		},
		{
			Name:           "Biff",
			Image:          BiffCharImage,
			SelectedColor:  color.RGBA{99, 26, 3, 255},
			FrameColor:     color.RGBA{238, 156, 50, 255},
			Controller:     &core.AreaController{},
			ControllerType: ComputerPlayer,
		},
		{
			Name:           "Elara",
			Image:          ElaraCharImage,
			SelectedColor:  color.RGBA{58, 59, 94, 255},
			FrameColor:     color.RGBA{121, 121, 203, 255},
			Controller:     &core.AreaController{},
			ControllerType: ComputerPlayer,
		},
		{
			Name:           "Mike Green",
			Image:          MikeGCharImage,
			SelectedColor:  color.RGBA{20, 104, 20, 255},
			FrameColor:     color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
		},
		{
			Name:           "Mike Violet",
			Image:          MikeVCharImage,
			SelectedColor:  color.RGBA{68, 21, 51, 255},
			FrameColor:     color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
		},
		{
			Name:           "Heather Green",
			Image:          HeatherGCharImage,
			SelectedColor:  color.RGBA{20, 104, 20, 255},
			FrameColor:     color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
		},
		{
			Name:           "Heather Violet",
			Image:          HeatherVCharImage,
			SelectedColor:  color.RGBA{68, 21, 51, 255},
			FrameColor:     color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
		},
	}
}
