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

// resolveCollision takes a Rect (tile/other object) and resolves the collision.
// It uses the velocity pointer to determine collision direction and adjusts the character's X/Y position.
// It also sets the velocity to 0.0 if collision occurred.
func (c *BaseCharacter) resolveCollision(otherRect Rect, axis CollisionAxis, v *float64) {
	characterRect := c.GetPushBox()

	if !characterRect.Intersects(otherRect) {
		return
	}

	if axis == AxisX {
		overlap := 0.0
		// Check the direction of movement using the velocity pointer
		if *v > 0 { // moving right
			overlap = characterRect.Right - otherRect.Left
		} else if *v < 0 { // moving left
			overlap = characterRect.Left - otherRect.Right
		}

		if math.Abs(overlap) > 0 {
			c.X -= overlap
			*v = 0.0 // Velocity consumed by collision
		}
	} else if axis == AxisY {
		overlap := 0.0
		// Check the direction of movement using the velocity pointer
		if *v > 0 { // moving down
			overlap = characterRect.Bottom - otherRect.Top
		} else if *v < 0 { // moving up
			overlap = characterRect.Top - otherRect.Bottom
		}

		if math.Abs(overlap) > 0 {
			c.Y -= overlap
			*v = 0.0 // Velocity consumed by collision
		}
	}
}

// HandleTileCollisions moves the character by the velocity in the pointer 'v' and resolves collisions.
// The caller must provide a pointer to the velocity variable they wish to modify (Vx/Vy or KnockbackVx/KnockbackVy).
func (c *BaseCharacter) HandleTileCollisions(level *Level, axis CollisionAxis, v *float64) {
	// Apply movement for the current axis
	if axis == AxisX {
		c.X += *v
	} else if axis == AxisY {
		c.Y += *v
	}

	// Check and resolve collisions with solid tiles
	characterHitBox := c.GetPushBox()

	// Determine the range of tiles to check
	// TileSize is assumed to be a globally accessible constant (from main.go)
	minX := int(math.Floor(characterHitBox.Left/TileSize)) - 1
	maxX := int(math.Floor(characterHitBox.Right/TileSize)) + 1
	minY := int(math.Floor(characterHitBox.Top/TileSize)) - 1
	maxY := int(math.Floor(characterHitBox.Bottom/TileSize)) + 1

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			tile := level.GetTile(x, y)
			if tile != nil && tile.solid {
				c.resolveCollision(tile.GetPushBox(), axis, v)
			}
		}
	}
}

// UpdateKnockback must be called by the concrete character's Update() method.
func (c *BaseCharacter) UpdateKnockback(level *Level) bool {
	if !c.IsKnockedBack() {
		return false
	}

	c.KnockbackFrames--

	// Apply knockback velocity. Since this is in BaseCharacter,
	// we use the X/Y fields and assume the concrete type handles collision if needed.
	// For now, we apply movement and let the concrete type handle tile collision
	// if it needs it (like the Player does).
	c.X += c.KnockbackVx
	c.Y += c.KnockbackVy

	c.HandleTileCollisions(level, AxisX, &c.KnockbackVx)
	c.HandleTileCollisions(level, AxisY, &c.KnockbackVy)

	return true // Knockback was active this frame
}
