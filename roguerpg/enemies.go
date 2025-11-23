package main

import (
	"math/rand"
)

type BlobEnemyState int

const (
	BlobIdle   BlobEnemyState = iota // Waiting phase
	BlobMoving                       // Moving phase (renamed from BlobWalking)
	BlobAttacking
	BlobHurt
	BlobDying

	// Movement constants
	MoveDurationFrames = 60 // Fixed 1 second movement (60 frames at 60 FPS)
	MaxWaitFrames      = 60 // Max 1 second wait (up to 60 frames)
)

type BlobEnemy struct {
	BaseSprite
	spriteSheet *SpriteSheet
	animations  map[BlobEnemyState]*Animation

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

	animations := map[BlobEnemyState]*Animation{
		BlobIdle:      NewAnimation([]int{0, 1, 2}, 20, true),
		BlobMoving:    NewAnimation([]int{0, 1, 2}, 20, true),
		BlobAttacking: NewAnimation([]int{0, 1, 2}, 20, true),
		BlobHurt:      NewAnimation([]int{5, 6, 5, 6}, 10, false),
		BlobDying:     NewAnimation([]int{5, 6, 10, 11, 12, 13, 14}, 10, false),
	}
	animations[BlobIdle].SetRandomFrame()

	spriteSheet := NewSpriteSheet(16, 16, 5, 3)
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
		animations:  animations,
		state:       BlobIdle,
		Health:      3,
		IsDead:      false,
		waitFrames:  rand.Intn(MaxWaitFrames) + 1,
	}
}

func (c *BlobEnemy) TakeDamage(damage int) {
	if c.IsDead || c.state == BlobDying || c.state == BlobHurt {
		return
	}

	c.state = BlobHurt

	c.Health -= damage
	if c.Health <= 0 {
		c.state = BlobDying
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
	c.animations[c.state].Update()
	c.srcRect = c.spriteSheet.Rect(c.animations[c.state].Frame())

	// Future extension: Check for Attacking proximity here first (Step 1)
	// if c.IsNearPlayer(level.Player) {
	//     c.state = BlobAttacking
	// }

	switch c.state {
	case BlobIdle:
		c.waitFrames--

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
		c.currentFrame++

		if c.currentFrame >= MoveDurationFrames {
			c.Location = c.moveTargetLocation // Snap to final position
			c.state = BlobIdle
			// Wait for a random time (up to 1 second)
			c.waitFrames = rand.Intn(MaxWaitFrames) + 1
			return
		}

		// Calculate interpolation factor (t: 0.0 -> 1.0)
		t := float64(c.currentFrame) / float64(MoveDurationFrames)
		dx := c.moveTargetLocation.X - c.moveStartLocation.X
		dy := c.moveTargetLocation.Y - c.moveStartLocation.Y
		c.X = c.moveStartLocation.X + dx*t
		c.Y = c.moveStartLocation.Y + dy*t

	case BlobAttacking:
		// For now, immediately return to idle/exploring state
		c.state = BlobIdle

	case BlobHurt:
		if c.animations[BlobHurt].IsFinished() {
			c.state = BlobIdle
		}

	case BlobDying:
		if c.animations[BlobDying].IsFinished() {
			c.IsDead = true
		}
	}
}
