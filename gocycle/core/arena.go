package core

const (
	ArenaWidth  = 50
	ArenaHeight = 50
)

type Arena struct {
	// Grid stores the state of each square (Open, Wall, or PlayerID)
	Grid   [][]Square
	Width  int
	Height int

	// Slice of all players, indexed by ID - 1 (e.g., PlayerID 1 is at index 0)
	Players []*Player
}

// isCollision checks if a given position is blocked by a Wall or an existing path.
func (a *Arena) isCollision(pos Vector) bool {
	// Check bounds collision. If the position is outside the grid, it's a collision.
	if pos.Y < 0 || pos.Y >= a.Height || pos.X < 0 || pos.X >= a.Width {
		return true
	}

	// Check grid content collision (Wall or another player's path/head)
	// Any value that is not Open (0) indicates a collision.
	squareState := a.Grid[pos.Y][pos.X]

	return squareState != Open
}

// clearPath iterates through a player's path and sets the corresponding Grid squares back to Open.
// This is used when a player dies to remove their trail.
func (a *Arena) clearPath(p *Player) {
	for _, pos := range p.Path {
		// Ensure the path position is within the arena bounds before clearing
		if pos.Y >= 0 && pos.Y < a.Height && pos.X >= 0 && pos.X < a.Width {
			a.Grid[pos.Y][pos.X] = Open
		}
	}
}

// RectangleGrid creates a grid initialized with a perimeter wall.
func rectangleGrid(width, height int) [][]Square {
	grid := make([][]Square, height)

	for y := range height {
		grid[y] = make([]Square, width)

		for x := range width {
			if y == 0 || y == height-1 || x == 0 || x == width-1 {
				grid[y][x] = Wall
			} else {
				grid[y][x] = Open
			}
		}
	}

	return grid
}

// NewArena is a constructor for setting up the initial grid and players.
func NewArena(w, h int, players []*Player) *Arena {
	grid := rectangleGrid(w, h)

	return &Arena{
		Grid:    grid,
		Width:   w,
		Height:  h,
		Players: players,
	}
}

// Update moves all living players and handles collisions.
func (a *Arena) Update() {
	// Collect all desired new positions and determine potential casualties
	newPositions := make(map[int]Vector) // map[PlayerID]NewPosition
	collidedIDs := make(map[int]bool)    // map[PlayerID]collided

	for _, p := range a.Players {
		if !p.IsAlive {
			continue
		}

		// Get next direction from the player's controller (Human/AI)
		nextDir := p.Controller.GetDirection(a, p.ID)
		p.Direction = nextDir
		nextPos := p.Position.Add(p.Direction) // Assuming a Vector.Add() method exists
		newPositions[p.ID] = nextPos

		// Check for collision with WALLS or STATIC PATHS (paths of already confirmed moving players)
		if a.isCollision(nextPos) {
			collidedIDs[p.ID] = true
		}
	}

	// Resolve simultaneous head-on collisions and apply movement

	// Check for players attempting to move into the same square
	// or crossing paths (A moves to B's old spot, B moves to A's old spot)
	for id, pos := range newPositions {
		if collidedIDs[id] {
			continue
		}

		// Check for multiple players attempting the same square
		for otherID, otherPos := range newPositions {
			if id == otherID {
				continue
			}
			if pos.Equals(otherPos) { // Assuming a Vector.Equals() method exists
				collidedIDs[id] = true
				collidedIDs[otherID] = true
			}
		}
	}

	// Finalize State Changes
	playersToMove := make(map[int]*Player)

	for _, p := range a.Players {
		if !p.IsAlive {
			continue
		}

		if collidedIDs[p.ID] {
			// Player dies and their trail is cleared *before* they update their position.
			a.clearPath(p)
			p.IsAlive = false
		} else {
			playersToMove[p.ID] = p
		}
	}

	// Only after all collisions are determined, move the survivors and draw their new path
	for _, p := range playersToMove {
		// Update the player's position first.
		oldPos := p.Position
		p.Position = newPositions[p.ID]

		// Now that the player has moved, mark the *old* position as the trail/path.
		a.Grid[oldPos.Y][oldPos.X] = Square(p.ID)

		// Add the new head position to the Path history.
		p.Path = append(p.Path, p.Position)
	}
}
