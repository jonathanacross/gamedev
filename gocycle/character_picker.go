package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	PickerRows   = 2
	PickerCols   = 5
	PickerSpaceX = 74
	PickerSpaceY = 90
	PickerStartX = 12
	PickerStartY = 40
)

type CharacterPicker struct {
	Characters         []*CharacterFrame
	MaxNumSelectable   int
	NumPlayer1Selected int
	NumPlayer2Selected int
	NumSelected        int
}

func NewCharacterPicker() *CharacterPicker {

	chars := []*CharacterFrame{}
	characterIndices := []int{6, 8, 0, 1, 2, 7, 9, 3, 4, 5}

	for i, charIdx := range characterIndices {
		x := float64(i%PickerCols)*PickerSpaceX + PickerStartX
		y := float64(i/PickerCols)*PickerSpaceY + PickerStartY
		char := NewCharacterFrame(&Characters[charIdx], x, y, CharacterNeutral, true)
		char.State = StateUnselected
		chars = append(chars, char)
	}

	return &CharacterPicker{
		Characters:         chars,
		MaxNumSelectable:   4,
		NumPlayer1Selected: 0,
		NumPlayer2Selected: 0,
		NumSelected:        0,
	}
}

func (cs *CharacterPicker) Draw(screen *ebiten.Image) {
	for _, char := range cs.Characters {
		char.Draw(screen)
	}
	drawTextAt(screen, "Pick your characters", 130, 5, text.AlignStart, color.White)
	drawTextAt(screen, "Player 1", 60, 25, text.AlignStart, color.White)
	drawTextAt(screen, "Player 2", 60, 115, text.AlignStart, color.White)
	drawTextAt(screen, "Computer opponents", 200, 25, text.AlignStart, color.White)
	if cs.IsValid() {
		drawTextAt(screen, "Press space to continue", 250, 205, text.AlignStart, color.White)
	}
}

func (cs *CharacterPicker) Update() {
	for _, char := range cs.Characters {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if char.HitBox().Contains(float64(x), float64(y)) {
				if char.State == StateSelected {
					char.State = StateUnselected
					cs.updateCounts(char, -1)
				} else if cs.canSelect(char) {
					char.State = StateSelected
					cs.updateCounts(char, 1)
				}
			}
		}
	}
}

func (cs *CharacterPicker) GetSelectedCharacters() []*CharData {
	selectedChars := []*CharData{}
	for _, ch := range cs.Characters {
		if ch.State == StateSelected {
			selectedChars = append(selectedChars, ch.CharData)
		}
	}
	return selectedChars
}

func (cs *CharacterPicker) canSelect(char *CharacterFrame) bool {
	// can only pick one human player 1
	if char.CharData.ControllerType == HumanFirstPlayer && cs.NumPlayer1Selected >= 1 {
		return false
	}

	// can only pick one human player 2
	if char.CharData.ControllerType == HumanSecondPlayer && cs.NumPlayer2Selected >= 1 {
		return false
	}

	// Limit total number of players
	if cs.NumSelected >= cs.MaxNumSelectable {
		return false
	}

	return true
}

func (cs *CharacterPicker) updateCounts(char *CharacterFrame, delta int) {
	switch char.CharData.ControllerType {
	case HumanFirstPlayer:
		cs.NumPlayer1Selected += delta
	case HumanSecondPlayer:
		cs.NumPlayer2Selected += delta
	}
	cs.NumSelected += delta
}

func (cs *CharacterPicker) IsValid() bool {
	return cs.NumSelected >= 2
}
