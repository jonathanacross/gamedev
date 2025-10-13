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
	ID             int
	Name           string
	Image          *ebiten.Image
	DarkColor      color.Color
	BrightColor    color.Color
	Controller     core.PlayerController
	ControllerType ControllerType
}

var Characters []CharData = loadCharData()
var NumCharacters = len(Characters)

// GetCharDataByID returns the CharData for a given 1-based core.Player ID.
func GetCharDataByID(playerID int) *CharData {
	// The Characters slice is 0-indexed, so we subtract 1.
	index := playerID - 1

	if index < 0 || index >= len(Characters) {
		return nil
	}

	return &Characters[index]
}

func loadCharData() []CharData {
	return []CharData{
		{
			ID:             1,
			Name:           "Milo",
			Image:          MiloCharImage,
			DarkColor:      color.RGBA{29, 70, 125, 255},
			BrightColor:    color.RGBA{3, 166, 224, 255},
			Controller:     &core.RandomAvoidingController{},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             2,
			Name:           "Sara",
			Image:          SaraCharImage,
			DarkColor:      color.RGBA{167, 151, 50, 255},
			BrightColor:    color.RGBA{248, 243, 79, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.10},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             3,
			Name:           "Dr. Q",
			Image:          DrQCharImage,
			DarkColor:      color.RGBA{23, 110, 114, 255},
			BrightColor:    color.RGBA{74, 199, 198, 255},
			Controller:     &core.WallHuggerController{},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             4,
			Name:           "Erica",
			Image:          EricaCharImage,
			DarkColor:      color.RGBA{156, 20, 38, 255},
			BrightColor:    color.RGBA{231, 64, 71, 255},
			Controller:     &core.RandomTurnerController{TurnProb: 0.005},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             5,
			Name:           "Biff",
			Image:          BiffCharImage,
			DarkColor:      color.RGBA{182, 70, 37, 255},
			BrightColor:    color.RGBA{238, 156, 50, 255},
			Controller:     &core.AreaController{},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             6,
			Name:           "Elara",
			Image:          ElaraCharImage,
			DarkColor:      color.RGBA{67, 67, 130, 255},
			BrightColor:    color.RGBA{121, 121, 203, 255},
			Controller:     &core.MinimaxAreaController{MaxDepth: 3},
			ControllerType: ComputerPlayer,
		},
		{
			ID:             7,
			Name:           "Mike Green",
			Image:          MikeGCharImage,
			DarkColor:      color.RGBA{20, 104, 20, 255},
			BrightColor:    color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
		},
		{
			ID:             8,
			Name:           "Mike Violet",
			Image:          MikeVCharImage,
			DarkColor:      color.RGBA{94, 33, 72, 255},
			BrightColor:    color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
		},
		{
			ID:             9,
			Name:           "Heather Green",
			Image:          HeatherGCharImage,
			DarkColor:      color.RGBA{20, 104, 20, 255},
			BrightColor:    color.RGBA{156, 224, 42, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanFirstPlayer,
		},
		{
			ID:             10,
			Name:           "Heather Violet",
			Image:          HeatherVCharImage,
			DarkColor:      color.RGBA{94, 33, 72, 255},
			BrightColor:    color.RGBA{193, 92, 153, 255},
			Controller:     core.NewHumanController(),
			ControllerType: HumanSecondPlayer,
		},
	}
}
