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

type RandomController struct{}

// Picks a completely random direction
func (rc *RandomController) GetDirection(arena *Arena, playerID int) Vector {
	dirs := []Vector{Up, Down, Left, Right}
	for {
		dir := dirs[rand.Intn(len(dirs))]
		// don't back into ourselves
		if !isOpposite(arena.Players[playerID-1].Direction, dir) {
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
