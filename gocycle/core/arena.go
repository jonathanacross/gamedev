package core

import "math"

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

	// Mark player starting positions on the grid
	for _, p := range players {
		grid[p.Position.Y][p.Position.X] = Square(p.ID)
	}

	return &Arena{
		Grid:    grid,
		Width:   w,
		Height:  h,
		Players: players,
	}
}

func NewArenaFromGrid(grid [][]Square, players []*Player) *Arena {
	// Make a copy of the grid to avoid mutating the input
	newGrid := make([][]Square, len(grid))
	for i := range grid {
		newGrid[i] = make([]Square, len(grid[i]))
		copy(newGrid[i], grid[i])
	}

	// Mark player starting positions on the grid
	for _, p := range players {
		newGrid[p.Position.Y][p.Position.X] = Square(p.ID)
	}

	return &Arena{
		Grid:    newGrid,
		Width:   len(grid[0]),
		Height:  len(grid),
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
		nextPos := p.Position.Add(p.Direction)
		newPositions[p.ID] = nextPos

		// Check for collision with walls or paths
		// This checks collision with *existing* paths/walls/heads on the grid.
		if a.isCollision(nextPos) {
			collidedIDs[p.ID] = true
		}
	}

	// Resolve simultaneous head-on collisions and apply movement

	// Check for players attempting to move into the same square.
	// This handles head-to-head collisions that bypass the grid check.
	for id, pos := range newPositions {
		if collidedIDs[id] {
			continue
		}

		// Check for multiple players attempting the same square
		for otherID, otherPos := range newPositions {
			if id == otherID {
				continue
			}
			if pos.Equals(otherPos) {
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
		// Store old position (the trail segment)
		oldPos := p.Position

		// Update the player's position to the new head location
		p.Position = newPositions[p.ID]

		// Mark the *old* position as the permanent trail/path.
		a.Grid[oldPos.Y][oldPos.X] = Square(p.ID)

		// Mark the *new* head position on the grid.
		// This prevents other players from moving into a surviving player's head
		// on subsequent turns.
		a.Grid[p.Position.Y][p.Position.X] = Square(p.ID)

		// Add the new head position to the Path history.
		p.Path = append(p.Path, p.Position)
	}
}

// bfsQueueItem is a helper struct for the multi-source BFS.
type bfsQueueItem struct {
	ID   int
	Pos  Vector
	Dist int
}

// findClosestAssignments uses a multi-source BFS (Voronoi partitioning) to determine
// which player controls which square based on Manhattan distance.
func (a *Arena) findClosestAssignments() ([][]int, [][]int) {
	// Initialize the distance and assignment grids
	distanceGrid := make([][]int, a.Height)
	assignmentGrid := make([][]int, a.Height)
	for y := 0; y < a.Height; y++ {
		distanceGrid[y] = make([]int, a.Width)
		assignmentGrid[y] = make([]int, a.Width)
		for x := 0; x < a.Width; x++ {
			// Initialize all squares to MaxInt (unvisited).
			distanceGrid[y][x] = math.MaxInt
		}
	}

	queue := []bfsQueueItem{}

	// Multi-Source Seeding: Start the BFS from all alive player heads
	for _, p := range a.Players {
		if !p.IsAlive {
			continue
		}

		pos := p.Position
		distanceGrid[pos.Y][pos.X] = 0
		assignmentGrid[pos.Y][pos.X] = p.ID

		queue = append(queue, bfsQueueItem{ID: p.ID, Pos: pos, Dist: 0})
	}

	// Run the BFS
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		nextDist := current.Dist + 1

		for _, dir := range []Vector{Up, Down, Left, Right} {
			nextPos := current.Pos.Add(dir)
			x, y := nextPos.X, nextPos.Y

			// Check bounds
			if y < 0 || y >= a.Height || x < 0 || x >= a.Width {
				continue
			}

			// Do not traverse into walls or existing paths.
			if a.Grid[y][x] != Open {
				continue
			}

			// --- Distance Comparison
			if nextDist < distanceGrid[y][x] {
				// Found a shorter path (new controlling player)
				distanceGrid[y][x] = nextDist
				assignmentGrid[y][x] = current.ID
				queue = append(queue, bfsQueueItem{ID: current.ID, Pos: nextPos, Dist: nextDist})
			} else if nextDist == distanceGrid[y][x] {
				// Found an equally short path (Neutral zone)
				if assignmentGrid[y][x] != current.ID {
					assignmentGrid[y][x] = 0 // Mark as neutral tie zone
				}
			}
		}
	}

	return distanceGrid, assignmentGrid
}

// based on the result of the BFS (the assignment grid).
func calculateScores(assignmentGrid [][]int, players []*Player) map[int]int {
	playerScores := make(map[int]int)

	// Initialize scores for all players
	for _, p := range players {
		playerScores[p.ID] = 0
	}

	// Iterate over the AssignmentGrid and tally the scores.
	for _, row := range assignmentGrid {
		for _, playerID := range row {
			if playerID > 0 { // Only count squares assigned to a player (ID > 0)
				playerScores[playerID]++
			}
		}
	}

	return playerScores
}

// ComputePlayerScores performs a multi-source BFS to calculate the controlled
// area for each player (Voronoi partitioning) and returns the scores.
func (a *Arena) ComputePlayerScores() map[int]int {
	// Run the Multi-Source BFS to get the assignment grid.
	_, assignmentGrid := a.findClosestAssignments()

	// Calculate the total area score based on the assignment grid.
	return calculateScores(assignmentGrid, a.Players)
}

// DeepCopy creates a complete copy of the Arena, including deep copies of the Grid and Players.
func (a *Arena) DeepCopy() *Arena {
	// 1. Copy the Grid (deep copy)
	newGrid := make([][]Square, a.Height)
	for y := range a.Height {
		newGrid[y] = make([]Square, a.Width)
		copy(newGrid[y], a.Grid[y])
	}

	// Copy the Players (deep copy)
	newPlayers := make([]*Player, len(a.Players))
	for i, p := range a.Players {
		// Deep copy the Player struct
		newP := *p
		// Deep copy the Path slice
		newP.Path = make([]Vector, len(p.Path))
		copy(newP.Path, p.Path)

		// The controller field cannot be deep-copied cleanly, but for a search
		// where we only care about its *state* (IsAlive, Position),
		// we can keep a reference to the original controller for the ID/type.
		// NOTE: Controller state (like InputQueue in HumanController) is irrelevant
		// in the sandbox, as we're injecting a direction, so this shallow copy is safe.
		newPlayers[i] = &newP
	}

	return &Arena{
		Grid:    newGrid,
		Width:   a.Width,
		Height:  a.Height,
		Players: newPlayers,
	}
}
