package core

import (
	"reflect"
	"testing"
)

// Helper function to create a small, simple arena for testing
func newTestArena(width, height int, players []*Player) *Arena {
	// A mock controller is needed for the Player struct, even if not used in Arena logic
	mockController := &mockController{}
	if players == nil {
		p1 := NewPlayer(1, Vector{X: 1, Y: 1}, Right, mockController)
		players = []*Player{p1}
	}

	// Manually create the arena to control dimensions and grid setup
	a := &Arena{
		Width:   width,
		Height:  height,
		Players: players,
	}

	// Initialize grid with Open squares
	a.Grid = make([][]Square, height)
	for y := 0; y < height; y++ {
		a.Grid[y] = make([]Square, width)
		for x := 0; x < width; x++ {
			a.Grid[y][x] = Open
		}
	}

	// Set up the outer Wall boundaries
	for x := 0; x < width; x++ {
		a.Grid[0][x] = Wall
		a.Grid[height-1][x] = Wall
	}
	for y := 0; y < height; y++ {
		a.Grid[y][0] = Wall
		a.Grid[y][width-1] = Wall
	}

	// Set initial player paths on the grid
	for _, p := range players {
		a.Grid[p.Position.Y][p.Position.X] = Square(p.ID)
	}

	return a
}

// Mock controller for testing players without running AI logic
type mockController struct{}

func (mc *mockController) GetDirection(*Arena, int) Vector { return Right }

// --- Basic Arena Functionality Tests ---

func TestNewArena(t *testing.T) {
	// Note: We test the NewArena function, which internally sets up walls.
	w, h := 10, 8
	p1 := NewPlayer(1, Vector{X: 5, Y: 5}, Right, &mockController{})

	arena := NewArena(w, h, []*Player{p1})

	if arena.Width != w || arena.Height != h {
		t.Errorf("NewArena dimensions incorrect. Got W=%d, H=%d, want W=%d, H=%d", arena.Width, arena.Height, w, h)
	}

	// Check if perimeter walls are set correctly
	if arena.Grid[0][5] != Wall || arena.Grid[7][5] != Wall || arena.Grid[4][0] != Wall || arena.Grid[4][9] != Wall {
		t.Errorf("NewArena failed to set up perimeter walls correctly.")
	}

	// Check if the player's initial position is marked on the grid
	if arena.Grid[5][5] != Square(1) {
		t.Errorf("NewArena failed to mark player 1's starting position. Got %v, want %v", arena.Grid[5][5], Square(1))
	}
}

func TestIsCollision(t *testing.T) {
	arena := newTestArena(5, 5, nil) // Arena is 5x5, walls at 0 and 4

	tests := []struct {
		name     string
		pos      Vector
		expected bool
	}{
		{"Boundary Wall (Top)", Vector{X: 2, Y: 0}, true},
		{"Boundary Wall (Right)", Vector{X: 4, Y: 2}, true},
		{"Out of Bounds (Negative)", Vector{X: -1, Y: 2}, true},
		{"Out of Bounds (Too Large)", Vector{X: 5, Y: 2}, true},
		{"Open Space", Vector{X: 2, Y: 2}, false},
	}

	// Add a static path collision
	arena.Grid[3][3] = Square(2)
	tests = append(tests, struct {
		name     string
		pos      Vector
		expected bool
	}{"Existing Path", Vector{X: 3, Y: 3}, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if arena.isCollision(tt.pos) != tt.expected {
				t.Errorf("isCollision(%v) got %v, want %v", tt.pos, arena.isCollision(tt.pos), tt.expected)
			}
		})
	}
}

func TestClearPath(t *testing.T) {
	arena := newTestArena(5, 5, nil)
	p1 := NewPlayer(1, Vector{X: 1, Y: 1}, Right, &mockController{})

	// Manually build a path for Player 1
	p1.Path = []Vector{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 3, Y: 1},
	}

	// Mark the path on the grid (simulating movement)
	for _, pos := range p1.Path {
		arena.Grid[pos.Y][pos.X] = Square(p1.ID)
	}

	// Check before clearing
	if arena.Grid[1][2] != Square(1) {
		t.Fatalf("Setup failed: Grid not marked with path.")
	}

	arena.clearPath(p1)

	// Check after clearing
	for _, pos := range p1.Path {
		if arena.Grid[pos.Y][pos.X] != Open {
			t.Errorf("Path at %v was not cleared. Got %v, want %v", pos, arena.Grid[pos.Y][pos.X], Open)
		}
	}
}

// --- Game Logic Test ---
// Mock controllers that force a specific direction for testing collisions
type forcedController struct{ Dir Vector }

func (fc *forcedController) GetDirection(*Arena, int) Vector { return fc.Dir }

func TestUpdate(t *testing.T) {
	// Test case 1: Simple movement
	p1StartPos := Vector{X: 1, Y: 1}
	p1 := NewPlayer(1, p1StartPos, Right, &forcedController{Dir: Right})
	arena := newTestArena(10, 10, []*Player{p1})

	arena.Update()

	p1NextPos := Vector{X: 2, Y: 1}
	if p1.Position != p1NextPos {
		t.Errorf("Simple Move Failed: Got position %v, want %v", p1.Position, p1NextPos)
	}
	if arena.Grid[p1StartPos.Y][p1StartPos.X] != Square(1) {
		t.Errorf("Simple Move Failed: Old position not marked as path. Got %v, want %v", arena.Grid[p1StartPos.Y][p1StartPos.X], Square(1))
	}
	if !p1.IsAlive {
		t.Errorf("Simple Move Failed: Player died unexpectedly.")
	}

	// Test case 2: Wall Collision
	p2StartPos := Vector{X: 1, Y: 2}
	p2 := NewPlayer(2, p2StartPos, Left, &forcedController{Dir: Left})
	arena = newTestArena(10, 10, []*Player{p2})

	arena.Update() // Move to (0, 2) which is a Wall

	if p2.IsAlive {
		t.Errorf("Wall Collision Failed: Player 2 should be dead.")
	}
	// Check if the path was cleared (path only has one point, the start position)
	if arena.Grid[p2StartPos.Y][p2StartPos.X] != Open {
		t.Errorf("Wall Collision Failed: Player 2 path not cleared.")
	}

	// Test case 3: Head-on Collision
	p3StartPos := Vector{X: 5, Y: 5}
	p4StartPos := Vector{X: 6, Y: 5}
	p3 := NewPlayer(3, p3StartPos, Right, &forcedController{Dir: Right}) // Move to (6, 5)
	p4 := NewPlayer(4, p4StartPos, Left, &forcedController{Dir: Left})   // Move to (5, 5)
	arena = newTestArena(10, 10, []*Player{p3, p4})

	arena.Update()

	if p3.IsAlive || p4.IsAlive {
		t.Errorf("Head-on Collision Failed: Both players should be dead. P3 Alive: %v, P4 Alive: %v", p3.IsAlive, p4.IsAlive)
	}
}

// TestUpdatePathCollision checks that a player attempting to move into a square
// just vacated by another player's path will die on the next tick (Tick 2).
func TestUpdatePathCollision(t *testing.T) {
	// Use the existing mockController that always returns Right
	controller := &mockController{}

	// Setup: P1 starts two squares behind P2.
	// P1 (ID 1) starts at (1, 3) moving Right.
	// P2 (ID 2) starts at (3, 3) moving Right.
	p1 := NewPlayer(1, Vector{X: 1, Y: 3}, Right, controller)
	p2 := NewPlayer(2, Vector{X: 3, Y: 3}, Right, controller)
	arena := newTestArena(7, 7, []*Player{p1, p2})

	// --- Tick 1: Both move right (P1 moves into Open, P2 moves into Open) ---
	arena.Update()

	// Expected: Both survive
	if !p1.IsAlive || !p2.IsAlive {
		t.Fatalf("Tick 1 failed: Both players should have survived a simple forward move.")
	}

	// Expected Positions: P1 at (2, 3), P2 at (4, 3)
	if !p1.Position.Equals(Vector{X: 2, Y: 3}) || !p2.Position.Equals(Vector{X: 4, Y: 3}) {
		t.Fatalf("Tick 1 failed: Incorrect positions. Got P1=%v, P2=%v", p1.Position, p2.Position)
	}

	// Crucial Grid Check: P2's trail segment at (3, 3) must be marked by its ID (2).
	if arena.Grid[3][3] != 2 {
		t.Fatalf("Tick 1 FAILED THE FIX: P2's trail at (3, 3) is not marked by its ID (Got: %v, Want: 2).", arena.Grid[3][3])
	}

	// --- Tick 2: P1 moves to (3, 3) (P2's trail) while P2 moves to (5, 3) ---
	// P1's next position is (3, 3). Grid[3][3] == 2 (P2's path). Collision should occur.
	arena.Update()

	// Expected: P1 dies trying to move into P2's trail segment. P2 survives.
	if p1.IsAlive {
		t.Error("Tick 2 FAILED: P1 should have died trying to move into P2's path/trail segment.")
	}
	if !p2.IsAlive {
		t.Error("Tick 2 FAILED: P2 should have survived a simple forward move.")
	}

	// Final Grid Check:
	// P1's path (1, 3), (2, 3) should be cleared (Open = 0)
	if arena.Grid[3][1] != Open || arena.Grid[3][2] != Open {
		t.Errorf("P1's path was not cleared upon death. Grid[3][1]=%v, Grid[3][2]=%v (Want 0)", arena.Grid[3][1], arena.Grid[3][2])
	}
	// P2's trail (4, 3) should be marked by P2's ID (2)
	if arena.Grid[3][4] != 2 {
		t.Errorf("P2's trail at (4, 3) is not marked. Got: %v (Want: 2)", arena.Grid[3][4])
	}
}

// --- Voronoi Scoring Algorithm Tests ---

func TestFindClosestAssignments(t *testing.T) {
	w, h := 7, 7
	p1 := NewPlayer(1, Vector{X: 1, Y: 1}, Right, &mockController{})
	p2 := NewPlayer(2, Vector{X: 5, Y: 5}, Left, &mockController{})
	arena := newTestArena(w, h, []*Player{p1, p2})

	// Add an internal wall to block P1's access to the center
	arena.Grid[3][3] = Wall

	_, assignmentGrid := arena.findClosestAssignments()

	// Assignments check for all key squares.
	tests := []struct {
		name     string
		pos      Vector
		expected int
		failLine int // Reference to original line for context
	}{
		// These assignments were incorrectly set to 0 due to the self-tie bug (now fixed)
		{"P1: (2, 2)", Vector{2, 2}, 1, 217},
		{"P2: (4, 4)", Vector{4, 4}, 2, 222},

		// These assignments were caused by propagation from the self-tie squares
		{"P1: (3, 2)", Vector{3, 2}, 1, 228},
		{"P2: (3, 4)", Vector{3, 4}, 2, 232},

		// FIX: (4,3) is distance 3 from P2 and 5 from P1. It should be P2 (2), not Neutral (0).
		{"P2: (4, 3) (Fixed Expectation)", Vector{4, 3}, 2, 237},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assignmentGrid[tt.pos.Y][tt.pos.X]
			if got != tt.expected {
				t.Errorf("Assignment at (%d,%d) expected %v, got %v (Original line: %d)",
					tt.pos.X, tt.pos.Y, tt.expected, got, tt.failLine)
			}
		})
	}

	// Additional check for the expected true tie squares (distance 4 from both)
	// (2,4) is a tie
	if assignmentGrid[4][2] != 0 {
		t.Errorf("Assignment at (4,2) expected Neutral (0), got %v", assignmentGrid[4][2])
	}
	// (4,2) is a tie
	if assignmentGrid[2][4] != 0 {
		t.Errorf("Assignment at (2,4) expected Neutral (0), got %v", assignmentGrid[2][4])
	}
}

func TestCalculateScores(t *testing.T) {
	// Mock players
	p1 := NewPlayer(1, Vector{X: 1, Y: 1}, Right, &mockController{})
	p2 := NewPlayer(2, Vector{X: 5, Y: 5}, Left, &mockController{})
	players := []*Player{p1, p2}

	// Clean 3x3 assignment (9 open squares):
	mockAssignments := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 2, 0},
		{0, 1, 0, 2, 0},
		{0, 1, 2, 2, 0},
		{0, 0, 0, 0, 0},
	}

	expectedScores := map[int]int{
		1: 4,
		2: 4,
	}

	// Calling the unexported function directly
	scores := calculateScores(mockAssignments, players)

	if !reflect.DeepEqual(scores, expectedScores) {
		t.Errorf("calculateScores failed.\nGot: %v\nWant: %v", scores, expectedScores)
	}
}

func TestComputePlayerScores(t *testing.T) {
	// Integration test for findClosestAssignments + calculateScores
	w, h := 7, 7
	p1 := NewPlayer(1, Vector{X: 1, Y: 1}, Right, &mockController{})
	p2 := NewPlayer(2, Vector{X: 5, Y: 5}, Left, &mockController{})
	arena := newTestArena(w, h, []*Player{p1, p2})

	// Arena size is 7x7. Playable inner area is 5x5 (25 squares).
	// Player heads take 2 squares. Available area for assignment is 23 squares.
	// Neutral squares (equidistant) are 3: (3,3), (2,4), (4,2).
	// Assigned squares: 23 - 3 = 20. Each player gets 10.

	expectedScores := map[int]int{
		1: 10,
		2: 10,
	}

	scores := arena.ComputePlayerScores()

	if !reflect.DeepEqual(scores, expectedScores) {
		t.Errorf("ComputePlayerScores integration failed (Symmetry case).\nGot: %v\nWant: %v", scores, expectedScores)
	}
}
