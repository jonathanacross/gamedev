package main

import (
	"math/bits"
)

type BitBoard uint64

func (b BitBoard) Set(index SquareIndex) BitBoard {
	return b | (1 << index)
}

func (b BitBoard) Clear(index SquareIndex) BitBoard {
	return b &^ (1 << index)
}

func (b BitBoard) Get(index SquareIndex) bool {
	return b&(1<<index) != 0
}

func (b BitBoard) GetSetBitIndices() []SquareIndex {
	indices := []SquareIndex{}
	for b != 0 {
		// Find the index of the least significant set bit
		idx := SquareIndex(bits.TrailingZeros64(uint64(b)))
		indices = append(indices, idx)

		// Clear the least significant set bit
		b &^= (1 << idx)
	}
	return indices
}

func (b BitBoard) GetNumSetBits() int {
	return bits.OnesCount64(uint64(b))
}

func (b BitBoard) ToString(boardSize int) string {
	var s string
	idx := 0
	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			if b.Get(SquareIndex(idx)) {
				s += "1 "
			} else {
				s += ". "
			}
			idx++
		}
		s += "\n"
	}
	return s
}
