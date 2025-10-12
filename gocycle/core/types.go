package core

// Square represents the state of a cell on the Arena grid.
type Square int

const (
	// Open (0): The cell is empty and available for movement.
	Open Square = 0
	// Wall (-1): The cell is permanently blocked by a border or internal obstacle.
	Wall Square = -1
	// Player paths are tracked by their positive Player ID (1, 2, 3...)
)

// Vector represents a position (x, y) or a direction (dx, dy).
type Vector struct {
	X int
	Y int
}

// Directions
var (
	Up    = Vector{0, -1}
	Down  = Vector{0, 1}
	Left  = Vector{-1, 0}
	Right = Vector{1, 0}
)

// Add returns a new Vector which is the sum of the current Vector and the provided direction.
func (v Vector) Add(dir Vector) Vector {
	return Vector{
		X: v.X + dir.X,
		Y: v.Y + dir.Y,
	}
}

func (v Vector) Subtract(dir Vector) Vector {
	return Vector{
		X: v.X - dir.X,
		Y: v.Y - dir.Y,
	}
}

// Equals checks if two Vectors have the same X and Y coordinates.
func (v Vector) Equals(other Vector) bool {
	return v.X == other.X && v.Y == other.Y
}

// isOpposite checks if two directions are directly opposite (180-degree turn).
func IsOpposite(d1, d2 Vector) bool {
	return d1.X == -d2.X && d1.Y == -d2.Y
}

func (v Vector) TurnRight() Vector {
	return Vector{
		X: -v.Y,
		Y: v.X,
	}
}
func (v Vector) TurnLeft() Vector {
	return Vector{
		X: v.Y,
		Y: -v.X,
	}
}
