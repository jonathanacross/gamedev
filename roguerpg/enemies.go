package main

import (
	"math/rand"
)

type BlobEnemyState int

const (
	BlobIdle   BlobEnemyState = iota // Waiting phase
	BlobMoving                       // Moving phase (renamed from BlobWalking)
	BlobAttacking
	BlobDying

	// Movement constants
	MoveDurationFrames = 60 // Fixed 1 second movement (60 frames at 60 FPS)
	MaxWaitFrames      = 60 // Max 1 second wait (up to 60 frames)
)

type BlobEnemy struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animation   *Animation
	Vx          float64
	Vy          float64

	Health int
	IsDead bool

	// AI
	state BlobEnemyState

	moveStartLocation  Location
	moveTargetLocation Location
	currentFrame       int // Frame counter for the current move or wait action
	waitFrames         int // Total frames to wait when idle

}

func NewBlobEnemy() *BlobEnemy {
	animation := NewAnimation([]int{0, 1, 2}, 20, true)
	animation.SetRandomFrame()
	spriteSheet := NewSpriteSheet(16, 16, 3, 1)
	hitbox := Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	return &BlobEnemy{
		BaseSprite: BaseSprite{
			Location: Location{
				X: 0,
				Y: 0,
			},
			drawOffset: Location{
				X: 8,
				Y: 8,
			},
			srcRect:    spriteSheet.Rect(0),
			image:      BlobSpritesImage,
			hitbox:     hitbox,
			debugImage: createDebugRectImage(hitbox),
		},
		spriteSheet: spriteSheet,
		animation:   animation,
		state:       BlobIdle,
		Health:      3, // Set initial health
		IsDead:      false,
		waitFrames:  rand.Intn(MaxWaitFrames) + 1,
	}
}

func (c *BlobEnemy) TakeDamage(damage int) {
	if c.IsDead || c.state == BlobDying {
		return
	}
	c.Health -= damage
	if c.Health <= 0 {
		c.state = BlobDying
		c.IsDead = true
		// TODO: transition to death animation, or start fading/removal process
		// For now, just stop movement immediately.
		c.Vx = 0
		c.Vy = 0
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *BlobEnemy) findNewTargetTile(level *Level) bool {
	// 1. Get current tile coordinates
	tx, ty := level.WorldToTile(c.Location)

	// Define the 4 cardinal directions for "adjacent square"
	directions := []struct{ dx, dy int }{
		{0, 1},  // Down
		{0, -1}, // Up
		{1, 0},  // Right
		{-1, 0}, // Left
	}

	// Shuffle directions to pick a random one first
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newTx := tx + dir.dx
		newTy := ty + dir.dy

		if !level.IsTileSolid(newTx, newTy) {
			// Found an open tile. Set up the movement.
			c.moveStartLocation = c.Location
			c.moveTargetLocation = level.TileToWorld(newTx, newTy)
			c.currentFrame = 0 // Reset frame counter for movement
			return true
		}
	}

	// No adjacent open tile found
	return false
}

func (c *BlobEnemy) Update(level *Level) {
	c.animation.Update()
	c.srcRect = c.spriteSheet.Rect(c.animation.Frame())

	// Future extension: Check for Attacking proximity here first (Step 1)
	// if c.IsNearPlayer(level.Player) {
	//     c.state = BlobAttacking
	// }

	switch c.state {
	case BlobIdle:
		// 1. Wait
		c.waitFrames--
		c.Vx = 0
		c.Vy = 0

		if c.waitFrames <= 0 {
			// Wait time is over. Look for a new target tile.
			if c.findNewTargetTile(level) {
				c.state = BlobMoving
			} else {
				// Enemy is cornered or blocked. Wait again.
				c.waitFrames = rand.Intn(MaxWaitFrames) + 1
			}
		}

	case BlobMoving:
		// 2. Move (Smoothly move from current to target over 60 frames)
		c.currentFrame++

		if c.currentFrame >= MoveDurationFrames {
			// Movement finished
			c.Location = c.moveTargetLocation // Snap to final position
			c.state = BlobIdle
			// 3. Wait for a random time (up to 1 second)
			c.waitFrames = rand.Intn(MaxWaitFrames) + 1
			return
		}

		// Calculate interpolation factor (t: 0.0 -> 1.0)
		t := float64(c.currentFrame) / float64(MoveDurationFrames)

		// Linear interpolation (Lerp) for smooth movement
		dx := c.moveTargetLocation.X - c.moveStartLocation.X
		dy := c.moveTargetLocation.Y - c.moveStartLocation.Y

		c.X = c.moveStartLocation.X + dx*t
		c.Y = c.moveStartLocation.Y + dy*t

		// Since we are setting position directly via Lerp, Vx/Vy are zero
		c.Vx = 0
		c.Vy = 0

	case BlobAttacking:
		// Placeholder for future Attacking logic
		c.Vx = 0
		c.Vy = 0
		// For now, immediately return to idle/exploring state
		c.state = BlobIdle
	case BlobDying:
		// TODO: Implement dying logic here
		c.Vx = 0
		c.Vy = 0
	}

	// if c.state == BlobWalking {
	// 	c.X += c.Vx
	// 	c.HandleTileCollisions(level, AxisX)
	// 	c.Y += c.Vy
	// 	c.HandleTileCollisions(level, AxisY)
	// }

	// c.X += c.Vx
	// c.HandleTileCollisions(level, AxisX)
	// c.Y += c.Vy
	// c.HandleTileCollisions(level, AxisY)
}
