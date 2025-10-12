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
	InputQueue []Vector
}

func (hc *HumanController) EnqueueDirection(dir Vector) {
	hc.InputQueue = append(hc.InputQueue, dir)
}

// GetDirection returns the direction stored in the controller.
// It also performs the check to prevent a 180-degree turn.
func (hc *HumanController) GetDirection(arena *Arena, playerID int) Vector {
	player := arena.Players[playerID-1]
	currentDir := player.Direction

	if len(hc.InputQueue) == 0 {
		// If empty, continue in the current direction.
		return currentDir
	}

	// Dequeue the first requested direction
	requestedDir := hc.InputQueue[0]

	// Update the queue: Remove the consumed request
	hc.InputQueue = hc.InputQueue[1:]

	// Continue in the same direction if the player tries to do a 180 degree turn
	if IsOpposite(currentDir, requestedDir) {
		return currentDir
	}

	return requestedDir
}

// NewHumanController initializes the controller
func NewHumanController() *HumanController {
	return &HumanController{}
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
		if isPossiblePlayerCollision(arena, playerID, dir) {
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

// isPossiblePlayerCollision checks if moving in a direction that may
// result in a simultaneous head-on collision with another player.
func isPossiblePlayerCollision(arena *Arena, playerID int, nextDir Vector) bool {
	player := arena.Players[playerID-1]

	// Calculate the current player's proposed next position
	nextPos := player.Position.Add(nextDir)

	// Check against all other alive players
	for _, otherPlayer := range arena.Players {
		// Skip self and dead players
		if !otherPlayer.IsAlive || otherPlayer.ID == playerID {
			continue
		}

		// Avoid moving into a square where the other player might move the
		// next turn
		for _, dir := range []Vector{Up, Down, Left, Right} {
			otherNextPos := otherPlayer.Position.Add(dir)
			if nextPos.Equals(otherNextPos) {
				return true
			}
		}
	}
	return false
}

type RandomTurnerController struct {
	TurnProb float64
}

func getRandomDir(m map[Vector]struct{}) Vector {
	keys := []Vector{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}

func (rt *RandomTurnerController) GetDirection(arena *Arena, playerID int) Vector {
	player := arena.Players[playerID-1]

	dirs := []Vector{Up, Down, Left, Right}
	safeDirs := make(map[Vector]struct{})
	maybeSafeDirs := make(map[Vector]struct{})
	for _, dir := range dirs {
		nextPos := player.Position.Add(dir)
		if arena.isCollision(nextPos) {
			continue
		}

		if isPossiblePlayerCollision(arena, playerID, dir) {
			maybeSafeDirs[dir] = struct{}{}
		} else {
			safeDirs[dir] = struct{}{}
		}
	}

	// Occasionally move randomly
	if rand.Float64() < rt.TurnProb && len(safeDirs) > 0 {
		return getRandomDir(safeDirs)
	}

	// Try to go forward
	if _, ok := safeDirs[player.Direction]; ok {
		return player.Direction
	}

	if len(safeDirs) > 0 {
		// Pick any random safe direction
		return getRandomDir(safeDirs)
	}

	if len(maybeSafeDirs) > 0 {
		// Pick any random maybe-safe direction
		return getRandomDir(maybeSafeDirs)
	}

	// Player is doomed, just go forward
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
		if isPossiblePlayerCollision(arena, playerID, dir) {
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

type WallHuggerController struct {
}

func (wh *WallHuggerController) GetDirection(arena *Arena, playerID int) Vector {
	player := arena.Players[playerID-1]
	if len(player.Path) < 2 {
		dirs := []Vector{Up, Down, Left, Right}
		return dirs[rand.Intn(len(dirs))]
	}

	dir := player.Direction
	leftDir := dir.TurnLeft()
	toLeftLoc := player.Position.Add(leftDir)
	toBackLeftLoc := toLeftLoc.Subtract(dir)

	rightDir := dir.TurnRight()
	toRightLoc := player.Position.Add(rightDir)
	toBackRightLoc := toRightLoc.Subtract(dir)

	dirs := []Vector{}
	if arena.isCollision(toBackLeftLoc) {
		// See a wall on the back left, try to follow it
		dirs = append(dirs, leftDir)
		dirs = append(dirs, dir)
		dirs = append(dirs, rightDir)
	} else if arena.isCollision(toBackRightLoc) {
		// See a wall on the back right, try to follow it
		dirs = append(dirs, rightDir)
		dirs = append(dirs, dir)
		dirs = append(dirs, leftDir)
	} else {
		// No walls to follow, prefer to go straight
		dirs = append(dirs, dir)
		dirs = append(dirs, rightDir)
		dirs = append(dirs, leftDir)
	}
	for _, dir := range dirs {
		nextPos := player.Position.Add(dir)
		if !arena.isCollision(nextPos) {
			return dir
		}
	}

	// No safe turns, just go forward and die
	return player.Direction
}
