package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type CharacterSelector struct {
	Characters         []*CharacterFrame
	MaxNumSelectable   int
	NumPlayer1Selected int
	NumPlayer2Selected int
	NumSelected        int
}

func NewCharacterSelector(
	x, y float64, spaceX, spaceY float64, rows, cols int) *CharacterSelector {

	chars := []*CharacterFrame{}
	characterIndices := []int{6, 8, 0, 1, 2, 7, 9, 3, 4, 5}

	for i, charIdx := range characterIndices {
		x := float64(i%cols)*spaceX + x
		y := float64(i/cols)*spaceY + y
		char := NewCharacterFrame(&Characters[charIdx], x, y, CharacterNeutral, true)
		char.State = StateUnselected
		chars = append(chars, char)
	}

	return &CharacterSelector{
		Characters:         chars,
		MaxNumSelectable:   4,
		NumPlayer1Selected: 0,
		NumPlayer2Selected: 0,
		NumSelected:        0,
	}
}

func (cs *CharacterSelector) Draw(screen *ebiten.Image) {
	for _, char := range cs.Characters {
		char.Draw(screen)
	}
	drawTextAt(screen, "Pick your characters", 130, 5, text.AlignStart, color.White)
	drawTextAt(screen, "Player 1", 60, 15, text.AlignStart, color.White)
	drawTextAt(screen, "Player 2", 60, 105, text.AlignStart, color.White)
	drawTextAt(screen, "Computer opponents", 200, 15, text.AlignStart, color.White)
	if cs.IsValid() {
		drawTextAt(screen, "Press space to continue", 250, 200, text.AlignStart, color.White)
	}
}

func (cs *CharacterSelector) Update() {
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

func (cs *CharacterSelector) GetSelectedCharacters() []*CharData {
	selectedChars := []*CharData{}
	for _, ch := range cs.Characters {
		if ch.State == StateSelected {
			selectedChars = append(selectedChars, ch.CharData)
		}
	}
	return selectedChars
}

func (cs *CharacterSelector) canSelect(char *CharacterFrame) bool {
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

func (cs *CharacterSelector) updateCounts(char *CharacterFrame, delta int) {
	switch char.CharData.ControllerType {
	case HumanFirstPlayer:
		cs.NumPlayer1Selected += delta
	case HumanSecondPlayer:
		cs.NumPlayer2Selected += delta
	}
	cs.NumSelected += delta
}

func (cs *CharacterSelector) IsValid() bool {
	return cs.NumSelected >= 2
}
