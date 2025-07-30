package main

import (
	"fmt"
	"strconv"
)

type Square int

const (
	SquareUndefined = iota
	SquareEmpty
	SquareRed
	SquareBlue
)

func (s Square) ToString() string {
	switch s {
	case SquareEmpty:
		return "."
	case SquareRed:
		return "x"
	case SquareBlue:
		return "o"
	default:
		return "?"
	}
}

const BoardSize = 7

type Board [BoardSize * BoardSize]Square

func NewBoard() Board {
	b := Board{}
	for i := range BoardSize * BoardSize {
		b[i] = SquareEmpty
	}
	return b
}

func (b *Board) ToString() string {
	result := ""

	for row := range BoardSize {
		result += strconv.Itoa(BoardSize-row) + " "
		for col := range BoardSize {
			result += b[row+col*BoardSize].ToString() + " "
		}
		result += "\n"
	}
	result += "  "
	for col := range BoardSize {
		result += string('A'+rune(col)) + " "
	}
	result += "\n"

	return result
}

func main() {
	b := NewBoard()
	fmt.Println(b.ToString())
}
