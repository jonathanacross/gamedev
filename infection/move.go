package main

import (
	"strconv"
)

// indexToHumanCoords converts plain integer to human-coordinates.
// E.g., if the board size were 3, then the board would look like this:
// |  3 | 0 1 2
// |  2 | 3 4 5
// |  1 | 6 7 8
// |    +------
// |      A B C
//
// And so index 0 becomes the coordinate string "A3", etc.
func IndexToHumanCoords(index int) string {
	col := index % BoardSize
	row := BoardSize - (index / BoardSize)
	return string('A'+rune(col)) + strconv.Itoa(row)
}

func HumanCoordsToIndex(coords string) int {
	col := int(coords[0] - 'A')
	row, _ := strconv.Atoi(coords[1:])
	return (BoardSize-row)*BoardSize + col
}

type Move struct {
	from int
	to   int
}

func (m Move) ToString() string {
	return IndexToHumanCoords(m.from) + "-" + IndexToHumanCoords(m.to)
}
