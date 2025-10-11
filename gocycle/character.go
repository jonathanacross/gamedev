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
	Score          int
}

var Characters []CharData = loadCharData()

func loadCharData() []CharData {
	return []CharData{
		{
			Name:           "Milo",
			Image:          MiloCharImage,
			SelectedColor:  color.RGBA{29, 70, 125, 255},
			FrameColor:     color.RGBA{3, 166, 224, 255},
			Controller:     &core.RandomAvoidingController{},
			ControllerType: ComputerPlayer,
			Score:          0,
		},
		{
			Name:           "Sara",
			Image:          SaraCharImage,
			SelectedColor:  color.RGBA{167, 151, 50, 255},
			FrameColor:     color.RGBA{248, 243, 79, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.30},
			ControllerType: ComputerPlayer,
			Score:          0,
		},
		{
			Name:           "Dr. Q",
			Image:          DrQCharImage,
			SelectedColor:  color.RGBA{23, 110, 114, 255},
			FrameColor:     color.RGBA{74, 199, 198, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.10},
			ControllerType: ComputerPlayer,
			Score:          0,
		},
		{
			Name:           "Erica",
			Image:          EricaCharImage,
			SelectedColor:  color.RGBA{156, 20, 38, 255},
			FrameColor:     color.RGBA{231, 64, 71, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.01},
			ControllerType: ComputerPlayer,
			Score:          0,
		},
		{
			Name:           "Biff",
			Image:          BiffCharImage,
			SelectedColor:  color.RGBA{182, 70, 37, 255},
			FrameColor:     color.RGBA{238, 156, 50, 255},
			Controller:     &core.AreaController{},
			ControllerType: ComputerPlayer,
			Score:          0,
		},
		{
			Name:           "Elara",
			Image:          ElaraCharImage,
			SelectedColor:  color.RGBA{67, 67, 130, 255},
			FrameColor:     color.RGBA{121, 121, 203, 255},
			Controller:     &core.AreaController{},
			ControllerType: ComputerPlayer,
			Score:          0,
		},
		{
			Name:           "Mike Green",
			Image:          MikeGCharImage,
			SelectedColor:  color.RGBA{20, 104, 20, 255},
			FrameColor:     color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
			Score:          0,
		},
		{
			Name:           "Mike Violet",
			Image:          MikeVCharImage,
			SelectedColor:  color.RGBA{94, 33, 72, 255},
			FrameColor:     color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
			Score:          0,
		},
		{
			Name:           "Heather Green",
			Image:          HeatherGCharImage,
			SelectedColor:  color.RGBA{20, 104, 20, 255},
			FrameColor:     color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
			Score:          0,
		},
		{
			Name:           "Heather Violet",
			Image:          HeatherVCharImage,
			SelectedColor:  color.RGBA{94, 33, 72, 255},
			FrameColor:     color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
			Score:          0,
		},
	}
}
