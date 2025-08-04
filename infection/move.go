package main

import (
	"fmt"
)

type SquareIndex int

// Move represents a single action on the board.
type Move struct {
	from SquareIndex
	to   SquareIndex
	jump bool
	pass bool
}

// CreateMove attempts to create a valid Move from a start and end square.
// It returns a Move and an error.
func CreateMove(from, to SquareIndex) (Move, error) {
	if from < 0 || from >= NumSquares {
		return Move{}, fmt.Errorf("from index %v out of bounds", from)
	}
	if to < 0 || to >= NumSquares {
		return Move{}, fmt.Errorf("to index %v out of bounds", to)
	}

	if adjacentBitboards[from].Get(to) {
		return Move{from: from, to: to, jump: false, pass: false}, nil
	}
	if jumpBitboards[from].Get(to) {
		return Move{from: from, to: to, jump: true, pass: false}, nil
	}

	return Move{}, fmt.Errorf("invalid move: from %d to %d is neither a step nor a jump", from, to)
}

// IsValid checks if a move is legal on the given board.
// It returns true and an empty string if the move is legal,
// otherwise false and an error message.
func (m *Move) IsValid(b *Board) (bool, string) {
	if m.pass {
		// A pass is only valid if there are no legal moves for the current player.
		legalMoves := b.GetLegalMoves()
		// Check for moves that are not 'pass' moves.
		for _, move := range legalMoves {
			if !move.pass {
				return false, "cannot pass when other legal moves are available"
			}
		}
		return true, ""
	}

	// Check if the from square is occupied by the current player.
	if b.playerToMove == White {
		if !b.white.Get(m.from) {
			return false, "from square not occupied by white"
		}
	} else { // Black
		if !b.black.Get(m.from) {
			return false, "from square not occupied by black"
		}
	}

	// Check if the to square is empty.
	if !b.empty.Get(m.to) {
		return false, "to square is not empty"
	}

	// Check if the move is a valid step or jump.
	// This is already handled by the CreateMove function, so we can re-use that logic.
	// We don't need to check again if we assume CreateMove has already been used.
	// However, for a truly robust IsValid, we can add a check here.
	if (m.jump && !jumpBitboards[m.from].Get(m.to)) || (!m.jump && !adjacentBitboards[m.from].Get(m.to)) {
		return false, "invalid move type (step or jump)"
	}

	return true, ""
}
