package main

import (
	"reflect"
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	// Check initial white pieces
	if !board.white.Get(GetIndex(0, 0)) {
		t.Errorf("Expected white piece at (0,0)")
	}
	if !board.white.Get(GetIndex(BoardSize-1, BoardSize-1)) {
		t.Errorf("Expected white piece at (%d,%d)", BoardSize-1, BoardSize-1)
	}

	// Check initial black pieces
	if !board.black.Get(GetIndex(0, BoardSize-1)) {
		t.Errorf("Expected black piece at (0,%d)", BoardSize-1)
	}
	if !board.black.Get(GetIndex(BoardSize-1, 0)) {
		t.Errorf("Expected black piece at (%d,0)", BoardSize-1)
	}

	// Check empty squares
	expectedEmptyCount := NumSquares - 4
	actualEmptyCount := 0
	for i := 0; i < NumSquares; i++ {
		if board.empty.Get(i) {
			actualEmptyCount++
		}
	}
	if actualEmptyCount != expectedEmptyCount {
		t.Errorf("Expected %d empty squares, got %d", expectedEmptyCount, actualEmptyCount)
	}
}

func TestBoardMove(t *testing.T) {
	type TestCase struct {
		name       string
		startBoard string
		move       Move
		expect     string
	}

	testcases := []TestCase{
		{
			name: "white step",
			startBoard: `
				o . . . o o x 
				. . x . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
			move: Move{
				from: GetIndex(0, 0),
				to:   GetIndex(1, 1),
				jump: false,
			},
			expect: `
				o . . . o o x 
				. o o . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
		},
		{
			name: "black step",
			startBoard: `
				o . . . o o x 
				. o o . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
			move: Move{
				from: GetIndex(0, 6),
				to:   GetIndex(1, 5),
				jump: false,
			},
			expect: `
				o . . . x x x 
				. o o . . x . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
		},
		{
			name: "white jump",
			startBoard: `
				o . . . x x x 
				. o o . . x . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
			move: Move{
				from: GetIndex(1, 1),
				to:   GetIndex(0, 3),
				jump: false,
			},
			expect: `
				o . . o o x x 
				. . o . . x . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
		},
		{
			name: "black jump",
			startBoard: `
				o . . o x x x 
				. . o . . x . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
			move: Move{
				from: GetIndex(0, 5),
				to:   GetIndex(1, 3),
				jump: false,
			},
			expect: `
				o . . x x . x 
				. . x x . x . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				. . . . . . . 
				x . . . . . o`,
		},
	}
	for _, tc := range testcases {
		board, err := NewBoardFromText(tc.startBoard)
		if err != nil {
			t.Fatalf("%s: Failed to create inital board from text: %v", tc.name, err)
		}
		board.Move(tc.move)
		got := board
		want, err := NewBoardFromText(tc.expect)
		if err != nil {
			t.Fatalf("%s: Failed to create final board from text: %v", tc.name, err)
		}
		if reflect.DeepEqual(want, got) {
			t.Errorf("%s: want\n%s, got\n%s", tc.name, want.String(), got.String())
		}
	}
}
