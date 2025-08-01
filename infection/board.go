package main

import (
	"fmt"
	"regexp"
)

// BoardSize must be <= 8
const BoardSize = 7
const NumSquares = BoardSize * BoardSize

type Board struct {
	white BitBoard
	black BitBoard
	empty BitBoard
}

// used for making moves (and generating moves)
var adjacentBitboards [NumSquares]BitBoard
var jumpBitboards [NumSquares]BitBoard

func init() {
	adjacentBitboards = genAdjacentBitboards()
	jumpBitboards = genJumpBitboards()
}

func GetIndex(row, col int) int {
	return row*BoardSize + col
}

func IndexToRowCol(index int) (int, int) {
	return index / BoardSize, index % BoardSize
}

func (b *Board) syncEmpty() {
	b.empty = ^(b.white | b.black)
}

func NewBoard() *Board {
	b := Board{
		white: BitBoard(0),
		black: BitBoard(0),
		empty: BitBoard(0),
	}
	// Set white/black pieces in opposite corners
	b.white = b.white.Set(GetIndex(0, 0))
	b.white = b.white.Set(GetIndex(BoardSize-1, BoardSize-1))
	b.black = b.black.Set(GetIndex(0, BoardSize-1))
	b.black = b.black.Set(GetIndex(BoardSize-1, 0))
	b.syncEmpty()
	return &b
}

func (b *Board) String() string {
	result := ""
	for r := range BoardSize {
		for c := range BoardSize {
			idx := GetIndex(r, c)
			if b.white.Get(idx) {
				result += "o "
			} else if b.black.Get(idx) {
				result += "x "
			} else {
				result += ". "
			}
		}
		result += "\n"
	}

	return result
}

type Offset struct {
	dx int
	dy int
}

// location index to bitboard of adjacent squares
func genOffsetBitboards(offsets []Offset) [NumSquares]BitBoard {
	bitboards := [NumSquares]BitBoard{}

	for idx := range NumSquares {
		row, col := IndexToRowCol(idx)

		bb := BitBoard(0)
		for _, offset := range offsets {
			new_row := row + offset.dy
			new_col := col + offset.dx
			if (new_row >= 0 && new_row < BoardSize) && (new_col >= 0 && new_col < BoardSize) {
				new_idx := GetIndex(new_row, new_col)
				bb = bb.Set(new_idx)
			}
		}
		bitboards[idx] = bb
	}
	return bitboards
}

func genAdjacentBitboards() [NumSquares]BitBoard {
	offsets := []Offset{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
	return genOffsetBitboards(offsets)
}

func genJumpBitboards() [NumSquares]BitBoard {
	offsets := []Offset{
		{-2, -2},
		{-2, -1},
		{-2, 0},
		{-2, 1},
		{-2, 2},
		{-1, -2},
		{-1, 2},
		{0, -2},
		{0, 2},
		{1, -2},
		{1, 2},
		{2, -2},
		{2, -1},
		{2, 0},
		{2, 1},
		{2, 2},
	}
	return genOffsetBitboards(offsets)
}

type Move struct {
	from int
	to   int
	jump bool
}

func (b *Board) Move(m Move) {
	if b.white.Get(m.from) {
		b.white = b.white.Set(m.to)
		if m.jump {
			b.white = b.white.Clear(m.from)
		}
		// change adjacent black squares to white
		infectedSquares := b.black & adjacentBitboards[m.to]
		b.black &^= infectedSquares
		b.white |= infectedSquares
	} else if b.black.Get(m.from) {
		b.black = b.black.Set(m.to)
		if m.jump {
			b.black = b.black.Clear(m.from)
		}
		// change adjacent white squares to black
		infectedSquares := b.white & adjacentBitboards[m.to]
		b.white &^= infectedSquares
		b.black |= infectedSquares
	}

	b.syncEmpty()
}

func NewBoardFromText(text string) (*Board, error) {
	reg := regexp.MustCompile(`[^xo.]+`)
	importantChars := reg.ReplaceAllString(text, "")
	if len(importantChars) != NumSquares {
		return nil, fmt.Errorf("string representation has wrong number of characters")
	}

	b := Board{}
	for i, char := range importantChars {
		switch char {
		case 'o':
			b.white = b.white.Set(i)
		case 'x':
			b.black = b.black.Set(i)
		default:
			b.empty = b.empty.Set(i)
		}
	}
	return &b, nil
}

func main() {
	b := NewBoard()
	b.black = b.black.Set(GetIndex(1, 2))
	b.white = b.white.Set(GetIndex(0, 5))
	b.white = b.white.Set(GetIndex(0, 4))
	b.Move(Move{
		from: GetIndex(0, 0),
		to:   GetIndex(1, 1),
		jump: false,
	})
	fmt.Println(b.String())
}
