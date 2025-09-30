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
