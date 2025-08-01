package main

import "math/rand"

type Engine interface {
	GenMove(board *Board) Move
}

func (b *Board) GetLegalMovees() []Move {
	moves := []Move{}

	startingLocs := []int{}
	if b.playerToMove == White {
		startingLocs = b.white.GetSetBitIndices()
	} else {
		startingLocs = b.black.GetSetBitIndices()
	}

	for _, from := range startingLocs {
		stepLocs := (adjacentBitboards[from] & b.empty).GetSetBitIndices()
		for _, to := range stepLocs {
			moves = append(moves, Move{from, to, false})
		}
		jumpLocs := (jumpBitboards[from] & b.empty).GetSetBitIndices()
		for _, to := range jumpLocs {
			moves = append(moves, Move{from, to, true})
		}
	}

	return moves
}

type RandomEngine struct{}

func (e *RandomEngine) GenMove(board *Board) Move {
	moves := board.GetLegalMovees()
	if len(moves) == 0 {
		// TODO: handle pass
		return Move{}
	}

	idx := rand.Intn(len(moves))
	return moves[idx]
}

type Human struct{}

func (h *Human) GenMove(board *Board) Move {
	return getMoveFromUser(board)
}
