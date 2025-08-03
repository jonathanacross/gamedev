package main

import (
	"math/rand"
	"time"
)

type Engine interface {
	GenMove(board *Board) Move
	RequiresHumanInput() bool
}

func (b *Board) GetLegalMoves() []Move {
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
			moves = append(moves, Move{
				from: from,
				to:   to,
				jump: false,
				pass: false,
			})
		}
		jumpLocs := (jumpBitboards[from] & b.empty).GetSetBitIndices()
		for _, to := range jumpLocs {
			moves = append(moves, Move{
				from: from,
				to:   to,
				jump: true,
				pass: false,
			})
		}
	}

	if len(moves) == 0 {
		moves = append(moves, Move{
			from: -1,
			to:   -1,
			jump: false,
			pass: true,
		})
	}

	return moves
}

type RandomEngine struct{}

func (e *RandomEngine) GenMove(board *Board) Move {
	moves := board.GetLegalMoves()
	idx := rand.Intn(len(moves))
	return moves[idx]
}

func (e *RandomEngine) RequiresHumanInput() bool {
	return false
}

type GreedyEngine struct{}

func ScoreForPlayer(b *Board, p Player) int {
	whiteScore, blackScore := b.Score()
	if p == White {
		return whiteScore - blackScore
	}
	return blackScore - whiteScore
}

func (e *GreedyEngine) GenMove(board *Board) Move {
	moves := board.GetLegalMoves()
	if len(moves) == 1 {
		return moves[0]
	}
	rand.Shuffle(len(moves), func(i, j int) {
		moves[i], moves[j] = moves[j], moves[i]
	})

	currPlayer := board.playerToMove
	bestMove := moves[0]
	bestScore := -9999
	for _, move := range moves {
		tmpBoard := board.Copy()
		tmpBoard.Move(move)
		score := ScoreForPlayer(tmpBoard, currPlayer)
		if score > bestScore {
			bestMove = move
			bestScore = score
		}
	}

	return bestMove
}

func (e *GreedyEngine) RequiresHumanInput() bool {
	return false
}

type SlowEngine struct{}

func (e *SlowEngine) GenMove(board *Board) Move {
	dummyEngine := RandomEngine{}
	time.Sleep(1 * time.Second)
	return dummyEngine.GenMove(board)
}

func (e *SlowEngine) RequiresHumanInput() bool {
	return false
}

type Human struct{}

func (h *Human) GenMove(board *Board) Move {
	return getMoveFromUser(board)
}

func (e *Human) RequiresHumanInput() bool {
	return true
}
