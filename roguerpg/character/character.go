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
	return c.GetPushBox()
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
			if level.IsTileSolid(x, y) {
				tileRect := core.Rect{
					Left:   float64(x * core.TileSize),
					Top:    float64(y * core.TileSize),
					Right:  float64((x + 1) * core.TileSize),
					Bottom: float64((y + 1) * core.TileSize),
				}

				t := c.calculateCollisionTime(characterRect, tileRect, axis, v)

				// A collision occurs if t is between -tolerance and 1.0
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

func (c *BaseCharacter) calculateCollisionTime(characterRect, tileRect core.Rect, axis core.CollisionAxis, v float64) float64 {
	t := 1.0 // Default to no collision (or collision > 1.0)

	if axis == core.AxisX {
		if !characterRect.IntersectsY(tileRect) {
			return 1.0
		}
		if v > 0 { // moving right
			t = (tileRect.Left - characterRect.Right) / v
		} else if v < 0 { // moving left
			t = (tileRect.Right - characterRect.Left) / v
		}
	} else if axis == core.AxisY {
		if !characterRect.IntersectsX(tileRect) {
			return 1.0
		}
		if v > 0 { // moving down
			t = (tileRect.Top - characterRect.Bottom) / v
		} else if v < 0 { // moving up
			t = (tileRect.Bottom - characterRect.Top) / v
		}
	}

	return t
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
