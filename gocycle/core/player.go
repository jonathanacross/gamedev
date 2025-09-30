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
		if !arena.isCollision(nextPos) {
			safeDirs = append(safeDirs, dir)
		}
	}

	if len(safeDirs) == 0 {
		// Going to die, just pick anything
		return dirs[rand.Intn(len(dirs))]
	}

	// Pick a safe direction
	return safeDirs[rand.Intn(len(safeDirs))]
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
		if !arena.isCollision(nextPos) {
			safeDirs = append(safeDirs, dir)
		}
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
