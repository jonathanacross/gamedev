package core

import (
	"math/rand"
)

type PlayerController interface {
	GetDirection(arena *Arena, playerID int) Vector
}

type Player struct {
	ID        int
	IsAlive   bool
	Position  Vector
	Direction Vector
	Path      []Vector

	Controller PlayerController
}

func NewPlayer(id int, position Vector, direction Vector, controller PlayerController) *Player {
	return &Player{
		ID:         id,
		IsAlive:    true,
		Position:   position,
		Direction:  direction,
		Path:       []Vector{position},
		Controller: controller,
	}
}

// HumanController implements PlayerController for human input.
type HumanController struct {
	RequestedDirection Vector
}

// GetDirection returns the direction stored in the controller.
// It also performs the check to prevent a 180-degree turn.
func (hc *HumanController) GetDirection(arena *Arena, playerID int) Vector {
	player := arena.Players[playerID-1]

	// If the requested direction is NOT a 180-degree turn from the current direction,
	// use the requested direction.
	if !IsOpposite(player.Direction, hc.RequestedDirection) {
		return hc.RequestedDirection
	}

	// If the keypress was an attempt to turn 180 degrees, ignore it
	// and continue moving in the current direction.
	return player.Direction
}

// NewHumanController initializes the controller
func NewHumanController(startDir Vector) *HumanController {
	return &HumanController{
		RequestedDirection: startDir,
	}
}

type RandomController struct{}

// Picks a completely random direction
func (rc *RandomController) GetDirection(arena *Arena, playerID int) Vector {
	dirs := []Vector{Up, Down, Left, Right}
	for {
		dir := dirs[rand.Intn(len(dirs))]
		// don't back into ourselves
		if !IsOpposite(arena.Players[playerID-1].Direction, dir) {
			return dir
		}
	}
}

type RandomAvoidingController struct{}

// Picks a random direction, but tries to avoid a collision on the next step
func (rc *RandomAvoidingController) GetDirection(arena *Arena, playerID int) Vector {
	dirs := []Vector{Up, Down, Left, Right}

	safeDirs := []Vector{}
	for _, dir := range dirs {
		nextPos := arena.Players[playerID-1].Position.Add(dir)
		if arena.isCollision(nextPos) {
			continue
		}
		if isImmediateFatalCollision(arena, playerID, dir) {
			continue
		}

		safeDirs = append(safeDirs, dir)
	}

	if len(safeDirs) == 0 {
		// Going to die, just pick anything
		return dirs[rand.Intn(len(dirs))]
	}

	// Pick a safe direction
	return safeDirs[rand.Intn(len(safeDirs))]
}

// isImmediateFatalCollision checks if moving in the proposed direction (nextDir)
// will result in a simultaneous head-on collision with any other player.
// It predicts the other player's move using their Path history for direction.
func isImmediateFatalCollision(arena *Arena, playerID int, nextDir Vector) bool {
	player := arena.Players[playerID-1]

	// 1. Calculate the current player's proposed next position
	nextPosA := player.Position.Add(nextDir)

	// 2. Check against all other alive players
	for _, otherPlayer := range arena.Players {
		// Skip self and dead players
		if !otherPlayer.IsAlive || otherPlayer.ID == playerID {
			continue
		}

		// A player's path starts with their initial position, so a player always
		// has at least one position in their Path slice.
		// If Path.Length == 1, they haven't moved yet, so we can't reliably predict
		// a direction based on history, and the collision check would rely on the
		// simultaneous collision logic in Arena.Update() later. We skip the prediction.
		if len(otherPlayer.Path) < 2 {
			continue
		}

		// Calculate the other player's predicted direction based on their last two positions.
		// Last element is current position (head). Second to last is the previous position.
		currentPosB := otherPlayer.Path[len(otherPlayer.Path)-1] // Head position
		prevPosB := otherPlayer.Path[len(otherPlayer.Path)-2]    // Previous position

		// The direction is (Current X - Previous X, Current Y - Previous Y)
		predictedDirB := Vector{
			X: currentPosB.X - prevPosB.X,
			Y: currentPosB.Y - prevPosB.Y,
		}

		// Calculate the other player's predicted next position
		nextPosB := currentPosB.Add(predictedDirB)

		// 3. Check for a simultaneous head-on collision (same next square)
		if nextPosA.Equals(nextPosB) {
			return true // Fatal collision detected
		}
	}
	return false
}

type RandomTurnerController struct {
	TurnProb float64
}

func (rt *RandomTurnerController) GetDirection(arena *Arena, playerID int) Vector {
	player := arena.Players[playerID-1]

	dirs := []Vector{Up, Down, Left, Right}
	safeDirs := []Vector{}
	for _, dir := range dirs {
		nextPos := player.Position.Add(dir)
		if arena.isCollision(nextPos) {
			continue
		}
		if isImmediateFatalCollision(arena, playerID, dir) {
			continue
		}

		safeDirs = append(safeDirs, dir)
	}

	// Going to die; pick anything
	if len(safeDirs) == 0 {
		return dirs[rand.Intn(len(dirs))]
	}

	// If we would die going forward, then force a turn
	straightNextPos := player.Position.Add(player.Direction)
	mustTurn := arena.isCollision(straightNextPos)
	if mustTurn {
		// Pick a random safe direction
		return safeDirs[rand.Intn(len(safeDirs))]
	}

	// Going forward is safe; but with some probability, make a turn anyway
	if rand.Float64() < rt.TurnProb {
		return safeDirs[rand.Intn(len(safeDirs))]
	}
	return player.Direction
}

// AreaController is a computer player that chooses the direction
// that maximizes its controlled area (Voronoi score).
type AreaController struct{}

func (ac *AreaController) GetDirection(arena *Arena, playerID int) Vector {
	player := arena.Players[playerID-1]
	bestDir := player.Direction
	maxScore := -1

	// Check four directions: Up, Down, Left, Right
	for _, dir := range []Vector{Up, Down, Left, Right} {
		// 1. Do not allow 180-degree turn
		if IsOpposite(player.Direction, dir) {
			continue
		}

		nextPos := player.Position.Add(dir)

		// 2. Check for immediate collision with walls/paths
		if arena.isCollision(nextPos) {
			continue
		}
		if isImmediateFatalCollision(arena, playerID, dir) {
			continue
		}

		// --- Core AI Logic: Simulate the move and check the score ---

		// Create a temporary "sandbox" arena for score calculation (shallow copy)
		sandboxArena := *arena // Copy the Arena struct itself

		// Create a temporary player slice for the sandbox arena (deep copy of players)
		sandboxArena.Players = make([]*Player, len(arena.Players))
		for i, p := range arena.Players {
			newP := *p // Shallow copy of player struct
			sandboxArena.Players[i] = &newP
		}

		// Find the player in the sandbox and update its state for the simulation
		sandboxPlayer := sandboxArena.Players[playerID-1]
		sandboxPlayer.Position = nextPos

		// Calculate the new controlled area score *after* the simulated move
		scores := sandboxArena.ComputePlayerScores()
		currentScore := scores[playerID]

		// Check if this move is better than the current best
		if currentScore > maxScore {
			maxScore = currentScore
			bestDir = dir
		}
	}

	// Fallback: If all directions lead to death/low score, or if maxScore remains -1,
	// we return the current direction (or the best direction found).
	if maxScore == -1 {
		return player.Direction
	}

	return bestDir
}
