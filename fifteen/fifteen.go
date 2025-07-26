package main

import (
	"fmt"
	"math/rand/v2"
)

type Fifteen struct {
	Grid [16]int
}

func NewFifteen() *Fifteen {
	return &Fifteen{
		Grid: [16]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, -1},
	}
}

func (f *Fifteen) Randomize() {
	for {
		for range 100 {
			i := rand.IntN(16)
			j := rand.IntN(16)
			f.Grid[i], f.Grid[j] = f.Grid[j], f.Grid[i]
		}
		if isSolvable(f.Grid) {
			break
		}
	}
}

func (f *Fifteen) ToString() string {
	return fmt.Sprintf(
		"%2d %2d %2d %2d\n"+
			"%2d %2d %2d %2d\n"+
			"%2d %2d %2d %2d\n"+
			"%2d %2d %2d %2d\n",
		f.Grid[0], f.Grid[1], f.Grid[2], f.Grid[3],
		f.Grid[4], f.Grid[5], f.Grid[6], f.Grid[7],
		f.Grid[8], f.Grid[9], f.Grid[10], f.Grid[11],
		f.Grid[12], f.Grid[13], f.Grid[14], f.Grid[15])
}

// GetEmptyLoc gets the x, y coordinates of the empty square
// (the square with value -1).
func (f *Fifteen) GetEmptyLoc() (int, int) {
	for i, v := range f.Grid {
		if v == -1 {
			return i % 4, i / 4
		}
	}

	return -1, -1
}

func gridLoc(x, y int) int {
	return y*4 + x
}

// Moves the tile (x,y) (and other squares if needed)
// toward the empty square.
func (f *Fifteen) Move(x, y int) error {
	emptyX, emptyY := f.GetEmptyLoc()

	if x < 0 || x > 3 || y < 0 || y > 3 {
		return fmt.Errorf("invalid location %d %d", x, y)
	}

	if x == emptyX {
		if y > emptyY {
			for i := emptyY; i < y; i++ {
				f.Grid[gridLoc(x, i)] = f.Grid[gridLoc(x, i+1)]
			}
		} else if y < emptyY {
			for i := emptyY; i > y; i-- {
				f.Grid[gridLoc(x, i)] = f.Grid[gridLoc(x, i-1)]
			}
		} else {
			return fmt.Errorf("cannot move empty square %d %d", x, y)
		}
	} else if y == emptyY {
		if x > emptyX {
			for i := emptyX; i < x; i++ {
				f.Grid[gridLoc(i, y)] = f.Grid[gridLoc(i+1, y)]
			}
		} else if x < emptyX {
			for i := emptyX; i > x; i-- {
				f.Grid[gridLoc(i, y)] = f.Grid[gridLoc(i-1, y)]
			}
		} else {
			return fmt.Errorf("cannot move empty square %d %d", x, y)
		}
	} else {
		return fmt.Errorf("can't move square %d %d as not in same row or col as empty square %d %d", x, y, emptyX, emptyY)
	}

	f.Grid[gridLoc(x, y)] = -1

	return nil
}

func (f *Fifteen) IsSolved() bool {
	soln := NewFifteen()
	return f.Grid == soln.Grid
}

// isSolvable determines if a 15-puzzle arrangement is solvable.
//
// A 15-puzzle (4x4 grid) is solvable if:
//  1. If the blank tile is on an odd row counting from the bottom (row 1 or 3),
//     the number of inversions must be even.
//  2. If the blank tile is on an even row counting from the bottom (row 2 or 4),
//     the number of inversions must be odd.
//
// In simpler terms, the parity of the number of inversions must be different
// from the parity of the blank tile's row number (counting from the bottom).
func isSolvable(tiles [16]int) bool {
	inversions := 0
	blankIndex := -1 // To store the 0-indexed position of the blank tile

	// Step 1: Calculate the number of inversions
	// An inversion is when a larger number precedes a smaller number.
	// We ignore the blank tile (-1) for inversion calculation.
	for i := 0; i < len(tiles)-1; i++ {
		if tiles[i] == -1 {
			blankIndex = i // Store the blank tile's index
			continue
		}
		for j := i + 1; j < len(tiles); j++ {
			if tiles[j] == -1 {
				if blankIndex == -1 { // Ensure blankIndex is set if it's the first element encountered
					blankIndex = j
				}
				continue
			}
			if tiles[i] > tiles[j] {
				inversions++
			}
		}
	}

	// If blankIndex wasn't found (shouldn't happen with valid input), handle it.
	if blankIndex == -1 {
		fmt.Println("Error: Blank tile (-1) not found in the puzzle.")
		return false // Or panic, depending on desired error handling
	}

	// Step 2: Determine the row of the blank space from the bottom
	// The puzzle is a 4x4 grid.
	// 0-indexed rows: 0, 1, 2, 3
	// Row 0 (top) -> 4th row from bottom
	// Row 1       -> 3rd row from bottom
	// Row 2       -> 2nd row from bottom
	// Row 3 (bottom) -> 1st row from bottom
	//
	// (blankIndex / 4) gives the 0-indexed row from the top.
	// 4 - (blankIndex / 4) gives the 1-indexed row from the bottom.
	blankRowFromBottom := 4 - (blankIndex / 4)

	// Step 3: Apply the solvability rule
	// For a 4x4 grid (even width):
	// - If the blank is on an even row from the bottom, inversions must be odd.
	// - If the blank is on an odd row from the bottom, inversions must be even.
	// This means (inversions % 2) != (blankRowFromBottom % 2) must be true for a solvable puzzle.
	return (inversions % 2) != (blankRowFromBottom % 2)
}
