package main

import (
	"flag"
	"fmt"
	"gocycle/core" // Assuming your core package is correctly imported
	"math/rand/v2"
	"os"
	"sort"
)

// --- 1. Stats and Controller Definitions ---

// PlayerStats tracks the performance of a specific Controller type.
type PlayerStats struct {
	Name        string
	GamesPlayed int
	TotalScore  int
}

// AIControllers is the master list of all controllers to be benchmarked.
var AIControllers = map[string]core.PlayerController{
	//"Random":             &core.RandomController{},
	"RandomAvoiding":          &core.RandomAvoidingController{},
	"RandomTurner_0.1":        &core.RandomTurnerController{TurnProb: 0.10},
	"RandomTurner_0.005":      &core.RandomTurnerController{TurnProb: 0.005},
	"WallHugger":              &core.WallHuggerController{},
	"AreaController":          &core.AreaController{},
	"MinimaxAreaController_3": &core.MinimaxAreaController{MaxDepth: 3},
}

// newControllerInstance creates a fresh, non-aliased instance of a controller.
func newControllerInstance(name string) core.PlayerController {
	switch name {
	//case "Random":
	//	return &core.RandomController{}
	case "RandomAvoiding":
		return &core.RandomAvoidingController{}
	case "RandomTurner_0.1":
		return &core.RandomTurnerController{TurnProb: 0.10}
	case "RandomTurner_0.005":
		return &core.RandomTurnerController{TurnProb: 0.005}
	case "WallHugger":
		return &core.WallHuggerController{}
	case "AreaController":
		return &core.AreaController{}
	case "MinimaxAreaController_3":
		return &core.MinimaxAreaController{MaxDepth: 3}
	}
	return &core.RandomController{} // Default fallback
}

// --- 2. Simulation Setup ---

// setupRound prepares the players for a round based on random selections.
func setupRound(selectedNames []string) []*core.Player {
	numPlayers := len(selectedNames)

	// Get starting positions from core.GetStartVectors
	arenaLocs := core.GetStartVectors(numPlayers)
	// Fixed initial directions for the 4 starting slots
	initialDirections := []core.Vector{core.Right, core.Left, core.Up, core.Down}

	players := make([]*core.Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		// Create a fresh controller instance for each player
		controller := newControllerInstance(selectedNames[i])
		players[i] = core.NewPlayer(i+1, arenaLocs[i], initialDirections[i], controller)
	}

	return players
}

func runBenchmark() {
	// --- Command Line Setup ---
	runsPtr := flag.Int("runs", 1000, "Number of simulation runs to execute.")
	flag.Parse()

	if *runsPtr <= 0 {
		fmt.Println("Error: Please specify a positive number of runs using -runs.")
		os.Exit(1)
	}

	fmt.Printf("Starting AI Benchmark: %d Runs\n", *runsPtr)
	fmt.Println("---------------------------------")

	// --- Initialization ---
	stats := make(map[string]*PlayerStats)
	controllerNames := make([]string, 0, len(AIControllers))
	for name := range AIControllers {
		stats[name] = &PlayerStats{Name: name}
		controllerNames = append(controllerNames, name)
	}

	// Assuming 4 grids (0, 1, 2, 3) are defined in core/grids.go
	const numGrids = 4
	const playersPerGame = 4

	// Check if there are enough controllers to run a 4-player game
	if len(controllerNames) < playersPerGame {
		fmt.Printf("Error: Need at least %d controllers defined. Found %d.\n", playersPerGame, len(controllerNames))
		os.Exit(1)
	}

	// --- Main Simulation Loop ---
	for run := 0; run < *runsPtr; run++ {
		// 1. Randomly select 4 unique AI controllers
		rand.Shuffle(len(controllerNames), func(i, j int) {
			controllerNames[i], controllerNames[j] = controllerNames[j], controllerNames[i]
		})
		selectedNames := controllerNames[:playersPerGame]

		// 2. Randomly select an arena grid (0 to numGrids-1)
		gridID := rand.IntN(numGrids)
		grid := core.GetGrid(gridID)

		// 3. Setup the game and run the simulation
		players := setupRound(selectedNames)

		// SimulateRound returns map[PlayerID]Score
		roundScores := core.SimulateRound(grid, players)

		// 4. Update Stats
		for i, p := range players {
			name := selectedNames[i]
			score := roundScores[p.ID]

			// Only update if the player was properly scored (score >= 0)
			if score >= 0 {
				stats[name].TotalScore += score
				stats[name].GamesPlayed++
			}
		}
	}

	// --- Display Results ---
	fmt.Println("\nBenchmark Results:")

	// Convert map to slice for sorting
	results := make([]*PlayerStats, 0, len(stats))
	for _, s := range stats {
		// Only show controllers that actually played
		if s.GamesPlayed > 0 {
			results = append(results, s)
		}
	}

	// Sort results by TotalScore descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalScore > results[j].TotalScore
	})

	fmt.Printf("%-24s %-15s %-10s %-11s\n", "Controller", "Total Score", "Games Played", "Norm Score")
	fmt.Println("------------------------ --------------- ---------- ----------")

	for _, s := range results {
		fmt.Printf("%-24s %-15d %-10d %-10f\n", s.Name, s.TotalScore, s.GamesPlayed, float64(s.TotalScore)/float64(s.GamesPlayed))
	}
}

/*
func main() {
	runBenchmark()
}
*/
