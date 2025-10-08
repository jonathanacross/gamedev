package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CharacterSelector struct {
	Characters    []*CharacterFrame
	NumSelectable int
}

func (cs *CharacterSelector) GetSelectedCount() int {
	count := 0
	for _, char := range cs.Characters {
		if char.State == StateSelected {
			count++
		}
	}
	return count
}

func (cs *CharacterSelector) Draw(screen *ebiten.Image) {
	for _, char := range cs.Characters {
		char.Draw(screen)
	}
}

func (cs *CharacterSelector) Update() {
	for _, char := range cs.Characters {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if char.HitBox().Contains(float64(x), float64(y)) {
				if char.State == StateSelected {
					char.State = StateUnselected
				} else if cs.GetSelectedCount() < cs.NumSelectable {
					char.State = StateSelected
				}
			}
		}
	}
}

func NewCharacterSelector(
	x, y float64, spaceX, spaceY float64, rows, cols int, numSelectable int,
	// TODO: remove this parameter after ui update
	characterIndices []int) *CharacterSelector {

	chars := []*CharacterFrame{}

	for i, charIdx := range characterIndices {
		x := float64(i%cols)*spaceX + x
		y := float64(i/cols)*spaceY + y
		char := NewCharacterFrame(&Characters[charIdx], x, y, CharacterNeutral, true)
		char.State = StateUnselected
		chars = append(chars, char)
	}

	return &CharacterSelector{
		Characters:    chars,
		NumSelectable: numSelectable,
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
