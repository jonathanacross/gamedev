package character

import (
	"math"
	"roguerpg/core"
)

// BaseCharacter holds the common state and logic for any combat-capable entity.
type BaseCharacter struct {
	core.BasePhysical

	Health     int
	MaxHealth  int
	Experience int

	// Knockback fields
	KnockbackVx     float64
	KnockbackVy     float64
	KnockbackFrames int

	// Stun fields
	StunFrames int

	// Internal state tracking for death/dying (specific implementation is handled by Player/Enemy)
	IsDeadFlag bool
	Dead       bool
}

func (c *BaseCharacter) GetHurtBox() core.Rect {
	return c.GetPushBox() // Use the inherited push box (from BasePhysical)
}

func (c *BaseCharacter) ApplyKnockback(force core.Vector, duration int) {
	// A generic BaseCharacter can only be knocked back if it's not logically dead.
	if c.Dead {
		return
	}
	c.KnockbackVx = force.X
	c.KnockbackVy = force.Y
	c.KnockbackFrames = duration
}

func (c *BaseCharacter) IsKnockedBack() bool {
	return c.KnockbackFrames > 0
}

func (c *BaseCharacter) IsDead() bool {
	return c.Dead
}

func (c *BaseCharacter) GetHealth() int {
	return c.Health
}

func (c *BaseCharacter) GetMaxHealth() int {
	return c.MaxHealth
}

func (c *BaseCharacter) ApplyStun(duration int) {
	if c.Dead {
		return
	}
	c.StunFrames = duration
}

func (c *BaseCharacter) IsStunned() bool {
	return c.StunFrames > 0
}

func (c *BaseCharacter) UpdateStun() bool {
	if c.StunFrames > 0 {
		c.StunFrames--
		return true
	}
	return false
}

// CheckAndApplyMovement performs the movement for the given velocity v and returns true
// if a collision occurred.  It modifies c.X or c.Y depending on axis.
func (c *BaseCharacter) CheckAndApplyMovement(level core.Level, axis core.CollisionAxis, v float64) bool {
	if v == 0.0 {
		return false
	}

	characterRect := c.GetPushBox()
	hitT := 1.0 // Fraction of movement completed before collision (0.0 to 1.0)

	// Define a small tolerance for floating point errors in the collision check
	const collisionTolerance float64 = 0.001

	// Determine the range of tiles to check.
	minX := int(math.Floor(characterRect.Left/core.TileSize)) - 1
	maxX := int(math.Floor(characterRect.Right/core.TileSize)) + 1
	minY := int(math.Floor(characterRect.Top/core.TileSize)) - 1
	maxY := int(math.Floor(characterRect.Bottom/core.TileSize)) + 1

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// level.GetTile returns Tile (struct).
			// In core/interfaces.go I commented: "Region of code where I define Level... Tile is a struct in core".
			// However, Level interface: IsTileSolid(tileX, tileY int) bool.
			// CheckAndApplyMovement used: tile := level.GetTile(x, y); if tile != nil && tile.solid ...
			// The Level interface in core (interface.go) defined: IsTileSolid(x, y int) bool.
			// It ALSO defined GetTile(x, y int) Tile? No, I need to check interface definition.

			// I defined `GetTile(x, y int) Tile` in `core/interfaces.go`.
			// But `Tile` is a struct, checking against nil is tricky if it returns value.
			// It should return *Tile.

			// Wait, Level interface needs to work for this code.
			// The code uses `tile.solid` and `tile.GetPushBox()`.
			// So `GetTile` should return `*core.Tile`.

			// Wait, I should verify `core/interfaces.go` definition of `Level`.
			// I'll proceed assuming I can fix interface if needed.
			// For now, I'll use IsTileSolid from interface if possible?
			// But I need `tile.GetPushBox()`.
			// `core.Tile` has `GetPushBox()`.
			// So `GetTile` must return `*core.Tile`.

			// If `GetTile` returns nil, it means no tile?
			// `level.GetTile` in original code returned `*Tile`.
			// So I will assume `level.GetTile(x,y)` returns `*core.Tile`.

			tileSolid := level.IsTileSolid(x, y)
			if tileSolid {
				// We need the tile's rect.
				// Tile rect is derived from position.
				// Original: tile.GetPushBox().
				// Construct a tile at that position to get its box?
				// Or use `level.TileToWorld`.
				// Tile always fills the grid square?
				// `core/sprites.go` NewTile: PushBox is (0,0) to (16,16).
				// `Tile` uses `PushBoxOffset`.
				// `tile.GetPushBox()` uses `TileToWorld` location?

				// If I don't have access to the Tile object, can I calculate the box?
				// Yes, generic tile at (x,y)
				// location := level.TileToWorld(x,y) -> gives center.
				// Wait, TileToWorld gives center.
				// NewTile location is what? Center? TopLeft?
				// Original `BuildLevel`: location x*TileSize, y*TileSize (TopLeft?)
				// `NewTile` sets Location.
				// `core` Tile: `GetPushBox` = Box + Loc.
				// If I know it's a solid tile, I can calculate its box without the Tile object.
				// Box is (x*TileSize, y*TileSize, (x+1)*TileSize, (y+1)*TileSize).

				// Let's rely on calculation instead of getting the Tile object if possible,
				// to reduce dependency on Tile implementation details if we want abstract Level.
				// But Level interface is in core, Tile is in core.
				// I'll stick to logical calculation for collision.

				tileRect := core.Rect{
					Left:   float64(x * core.TileSize),
					Top:    float64(y * core.TileSize),
					Right:  float64((x + 1) * core.TileSize),
					Bottom: float64((y + 1) * core.TileSize),
				}
				t := 1.0

				// Calculate collision time 't' (Swept AABB logic)
				if axis == core.AxisX {
					if !characterRect.IntersectsY(tileRect) {
						continue // Skip this tile, no Y-overlap
					}
					if v > 0 { // moving right (Right edge hits Left edge)
						t = (tileRect.Left - characterRect.Right) / v
					} else if v < 0 { // moving left (Left edge hits Right edge)
						t = (tileRect.Right - characterRect.Left) / v
					}
				} else if axis == core.AxisY {
					if !characterRect.IntersectsX(tileRect) {
						continue // Skip this tile, no X-overlap
					}
					if v > 0 { // moving down (Bottom edge hits Top edge)
						t = (tileRect.Top - characterRect.Bottom) / v
					} else if v < 0 { // moving up (Top edge hits Bottom edge)
						t = (tileRect.Bottom - characterRect.Top) / v
					}
				}

				// A collision occurs if t is between -tolerance and 1.0
				// -tolerance ensures we detect collisions even if the boxes are slightly overlapping
				// due to previous floating point math.
				if t >= -collisionTolerance && t < 1.0 {
					hitT = math.Min(hitT, t)
				}
			}
		}
	}

	// Apply the movement up to the point of impact (minus epsilon)
	const separationEpsilon float64 = 0.0001
	moveFraction := hitT // Start with the fraction of movement allowed

	// Only apply the separation epsilon if a collision was detected
	if hitT < 1.0 {
		moveFraction = math.Max(0.0, hitT-separationEpsilon)
	}

	// Apply movement to position.
	moveDistance := v * moveFraction
	if axis == core.AxisX {
		c.Loc.X += moveDistance // Using Loc field from BaseSprite (exported in core)
	} else if axis == core.AxisY {
		c.Loc.Y += moveDistance
	}

	// Return true if a collision was detected before the full movement (hitT < 1.0)
	return hitT < 1.0
}

// ResolveTileCollision applies the default response (stopping) to a velocity vector.
// This function can be overridden or extended for different behaviors (e.g., bounce).
func (c *BaseCharacter) ResolveTileCollision(axis core.CollisionAxis, v float64) float64 {
	// Default response: stop movement along this axis
	return 0.0
}

// HandleTileCollisions performs collision checks, moves the character,
// and returns the resolved velocity for that axis.
func (c *BaseCharacter) HandleTileCollisions(level core.Level, axis core.CollisionAxis, v float64) float64 {
	// Move the character and check if a collision occurred.
	// CheckAndApplyMovement uses the velocity v to determine distance,
	// and moves the character's position (c.X/c.Y).
	hit := c.CheckAndApplyMovement(level, axis, v)

	// If a collision occurred, apply the collision response.
	if hit {
		return c.ResolveTileCollision(axis, v)
	}

	// If no collision, the full velocity is kept for the next frame.
	return v
}

// UpdateKnockback must be called by the concrete character's Update() method.
func (c *BaseCharacter) UpdateKnockback(level core.Level) bool {
	if !c.IsKnockedBack() {
		return false
	}

	c.KnockbackFrames--

	// Update character position with knockback velocity.
	c.HandleTileCollisions(level, core.AxisX, c.KnockbackVx)
	c.HandleTileCollisions(level, core.AxisY, c.KnockbackVy)

	return true // Knockback was active this frame
}
