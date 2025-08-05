package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Engine interface {
	GenMove(board *Board) Move
	RequiresHumanInput() bool
}

func (b *Board) GetLegalMoves() []Move {
	moves := []Move{}

	startingLocs := []SquareIndex{}
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

type MinimaxEngine struct {
	maxDepth int
}

// alpha represents the minimum score that white can guarantee.
// Beta is the max score that black can guarantee
func minimax(board *Board, depth int, alpha int, beta int, player Player) (bestMove Move, score int, numEvals int) {
	if depth == 0 || board.IsGameOver() {
		//fmt.Printf("depth 0, board = \n%s\n", board.String())
		w, b := board.Score()
		return Move{}, w - b, 1
	}
	bestMove = Move{}
	numEvals = 0

	if player == White {
		//fmt.Printf("evaluating white, depth: %d\n, alpha = %d, beta = %d\n", depth, alpha, beta)
		maxEval := math.MinInt
		moves := board.GetLegalMoves()
		for _, move := range moves {
			child := board.Copy()
			child.Move(move)
			//fmt.Printf("%*s  move %v for player %v\n", 2*depth, "", move, player)
			_, eval, evalCount := minimax(child, depth-1, alpha, beta, Black)
			numEvals += evalCount
			//fmt.Printf("%*s  eval = %d\n", 2*depth, "", eval)
			if eval > maxEval {
				bestMove = move
				maxEval = eval
				//fmt.Printf("found new best move %v with score %d\n", move, eval)
			}
			// update worst possible score for white
			alpha = max(alpha, eval)
			if beta <= alpha {
				// black can guarantee a score of beta,
				// so would never pick this branch.  No
				// need to explore further
				break
			}
		}
		return bestMove, maxEval, numEvals
	} else {
		//fmt.Printf("evaluating black, depth: %d\n, alpha = %d, beta = %d\n", depth, alpha, beta)
		minEval := math.MaxInt
		moves := board.GetLegalMoves()
		for _, move := range moves {
			child := board.Copy()
			child.Move(move)
			_, eval, evalCount := minimax(child, depth-1, alpha, beta, White)
			numEvals += evalCount
			//fmt.Printf("%*s  move %v for player %v\n", 2*depth, "", move, player)
			minEval = min(eval, minEval)
			//fmt.Printf("%*s  eval = %d\n", 2*depth, "", eval)
			if eval < minEval {
				bestMove = move
				minEval = eval
			}
			beta = min(beta, eval)
			if beta <= alpha {
				// white can guarantee a score of alpha,
				// so would never pick this branch.  No
				// need to explore further
				break
			}
		}
		return bestMove, minEval, numEvals
	}
}

func (e *MinimaxEngine) GenMove(board *Board) Move {
	bestMove, _, numEvals := minimax(board, e.maxDepth, math.MinInt, math.MaxInt, board.playerToMove)
	fmt.Printf("Minimax evaluated %d moves\n", numEvals)
	return bestMove
}

func (e *MinimaxEngine) RequiresHumanInput() bool {
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

type HumanEngine struct{}

func (e HumanEngine) GenMove(b *Board) Move {
	return getMoveFromUser(b)
}

func (e HumanEngine) RequiresHumanInput() bool {
	return true
}
