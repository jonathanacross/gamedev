package core

import (
	"fmt"
)

// SimulateRound runs a complete game on the given grid with the provided players
// until only one (or zero) player remains.
// It returns a map where the key is the Player ID and the value is the final score.
func SimulateRound(grid [][]Square, players []*Player) map[int]int {
	arena := NewArenaFromGrid(grid, players)

	// A pool of ranks (scores) to distribute.
	// We initialize this list from lowest score to highest score
	// (0, 2, 4, 6 for 4 players) to match the time order of death (last place first).
	numPlayers := len(players)
	remainingRanks := make([]int, numPlayers)
	for rank := range numPlayers {
		remainingRanks[rank] = rank * 2
	}

	// Map to store the final score for each player (ID -> Score).
	// An initial value of -1 means 'not yet scored'.
	roundScores := make(map[int]int, numPlayers)
	for _, p := range players {
		roundScores[p.ID] = -1
	}

	// Tracks the 'IsAlive' status from the *previous* update for collision-based scoring.
	previousIsAlive := make([]bool, numPlayers)
	for i := range numPlayers {
		previousIsAlive[i] = players[i].IsAlive
	}

	// Game loop
	for {
		// 1. Check for end condition
		numActivePlayers := 0
		for _, p := range players {
			if p.IsAlive {
				numActivePlayers++
			}
		}

		// The round ends when one or zero players are active
		if numActivePlayers <= 1 {
			break
		}

		// 2. Game Tick
		arena.Update()

		// 3. Score recently dead players
		HandleScoreUpdate(players, previousIsAlive, roundScores, &remainingRanks)

		// 4. Update previous status for the next tick
		for i, p := range players {
			previousIsAlive[i] = p.IsAlive
		}

		// Safety break (optional but good practice for simulations)
		if len(arena.Players[0].Path) > 50*50 {
			// Arena is full or loop is too long (should not happen in a normal run)
			fmt.Println("Warning: Simulation exceeded max path length. Breaking.")
			break
		}
	}

	// 5. Final Scoring (assign remaining scores to winner(s) or the last tie group)
	ScoreRemainingPlayers(players, roundScores, remainingRanks)

	return roundScores
}

// handleScoreUpdate checks for players who died *this tick* and scores them.
func HandleScoreUpdate(players []*Player, previousIsAlive []bool, roundScores map[int]int, remainingRanks *[]int) {
	diedThisTurn := []int{} // List of IDs of players who died in this tick

	for i, p := range players {
		// If the player was alive last tick but is now dead: they died this turn
		if previousIsAlive[i] && !p.IsAlive {
			diedThisTurn = append(diedThisTurn, p.ID)
		}
	}

	if len(diedThisTurn) > 0 {
		// Calculate the score for the group that died in a tie
		numDied := len(diedThisTurn)
		if numDied > len(*remainingRanks) {
			// Should not happen, but a safeguard against index out of bounds
			numDied = len(*remainingRanks)
		}

		score := calculateAverageScore(*remainingRanks, numDied)

		// Assign the score to all players who died in the tie
		for _, id := range diedThisTurn {
			roundScores[id] = score
		}

		// Remove the ranks that were consumed by this death group
		*remainingRanks = (*remainingRanks)[numDied:]
	}
}

// ScoreRemainingPlayers scores the final active player(s) or the last tie group.
func ScoreRemainingPlayers(players []*Player, roundScores map[int]int, remainingRanks []int) {
	unscoredPlayers := []*Player{}
	for _, p := range players {
		// Only consider players who are still alive OR were the last to die but haven't been scored.
		if roundScores[p.ID] == -1 {
			unscoredPlayers = append(unscoredPlayers, p)
		}
	}

	if len(unscoredPlayers) > 0 {
		score := calculateAverageScore(remainingRanks, len(unscoredPlayers))

		// Assign the score to all remaining players (the winner(s) or the last tie group)
		for _, p := range unscoredPlayers {
			roundScores[p.ID] = score
		}
	}
}

// calculateAverageScore finds the integer-averaged score for a group of tied players.
// It sums the top 'num' ranks and divides by 'num' (rounded down).
func calculateAverageScore(ranks []int, num int) int {
	if num == 0 || len(ranks) == 0 {
		return 0
	}

	if num > len(ranks) {
		num = len(ranks)
	}

	sum := 0
	for i := 0; i < num; i++ {
		sum += ranks[i]
	}

	// The division must be integer division (floor) to match the game logic.
	return sum / num
}
