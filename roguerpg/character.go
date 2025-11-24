// character.go
package main

import "math"

// BaseCharacter holds the common state and logic for any combat-capable entity.
type BaseCharacter struct {
	BasePhysical

	Health    int
	MaxHealth int

	// Knockback fields
	KnockbackVx     float64
	KnockbackVy     float64
	KnockbackFrames int

	// Internal state tracking for death/dying (specific implementation is handled by Player/Enemy)
	isDead bool
}

func (c *BaseCharacter) GetHurtBox() Rect {
	return c.GetPushBox() // Use the inherited push box (from BasePhysical)
}

func (c *BaseCharacter) ApplyKnockback(force Vector, duration int) {
	// A generic BaseCharacter can only be knocked back if it's not logically dead.
	if c.isDead {
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
	return c.isDead
}

// CheckAndApplyMovement performs the movement for the given velocity v and returns true
// if a collision occurred.  It modifies c.X or c.Y depending on axis.
func (c *BaseCharacter) CheckAndApplyMovement(level *Level, axis CollisionAxis, v float64) bool {
	if v == 0.0 {
		return false
	}

	characterRect := c.GetPushBox()
	hitT := 1.0 // Fraction of movement completed before collision (0.0 to 1.0)

	// Define a small tolerance for floating point errors in the collision check
	const collisionTolerance float64 = 0.001

	// Determine the range of tiles to check.
	minX := int(math.Floor(characterRect.Left/TileSize)) - 1
	maxX := int(math.Floor(characterRect.Right/TileSize)) + 1
	minY := int(math.Floor(characterRect.Top/TileSize)) - 1
	maxY := int(math.Floor(characterRect.Bottom/TileSize)) + 1

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			tile := level.GetTile(x, y)
			if tile != nil && tile.solid {
				tileRect := tile.GetPushBox()
				t := 1.0

				// Calculate collision time 't' (Swept AABB logic)
				if axis == AxisX {
					if !characterRect.IntersectsY(tileRect) {
						continue // Skip this tile, no Y-overlap
					}
					if v > 0 { // moving right (Right edge hits Left edge)
						t = (tileRect.Left - characterRect.Right) / v
					} else if v < 0 { // moving left (Left edge hits Right edge)
						t = (tileRect.Right - characterRect.Left) / v
					}
				} else if axis == AxisY {
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
	if axis == AxisX {
		c.X += moveDistance
	} else if axis == AxisY {
		c.Y += moveDistance
	}

	// Return true if a collision was detected before the full movement (hitT < 1.0)
	return hitT < 1.0
}

// ResolveTileCollision applies the default response (stopping) to a velocity vector.
// This function can be overridden or extended for different behaviors (e.g., bounce).
func (c *BaseCharacter) ResolveTileCollision(axis CollisionAxis, v float64) float64 {
	// Default response: stop movement along this axis
	return 0.0
}

// HandleTileCollisions performs collision checks, moves the character,
// and returns the resolved velocity for that axis.
func (c *BaseCharacter) HandleTileCollisions(level *Level, axis CollisionAxis, v float64) float64 {
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
func (c *BaseCharacter) UpdateKnockback(level *Level) bool {
	if !c.IsKnockedBack() {
		return false
	}

	c.KnockbackFrames--

	// Update character position with knockback velocity.
	c.HandleTileCollisions(level, AxisX, c.KnockbackVx)
	c.HandleTileCollisions(level, AxisY, c.KnockbackVy)

	return true // Knockback was active this frame
}
