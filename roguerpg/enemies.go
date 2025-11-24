package main

import (
	"math/rand"
)

type BlobEnemyState int

const (
	BlobIdle BlobEnemyState = iota
	BlobMoving
	BlobAttacking
	BlobHurt
	BlobDying

	// Movement constants
	BlobMoveSpeed float64 = 0.5
	MaxWaitFrames int     = 60 // Max 1 second wait (up to 60 frames)
)

type BlobEnemy struct {
	BaseCharacter
	spriteSheet *SpriteSheet
	animations  map[BlobEnemyState]*Animation

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
		BaseCharacter: BaseCharacter{
			BasePhysical: BasePhysical{
				BaseSprite: BaseSprite{
					Location: Location{
						X: 0,
						Y: 0,
					},
					drawOffset: Location{
						X: 8,
						Y: 8,
					},
					srcRect: spriteSheet.Rect(0),
					image:   BlobSpritesImage,
				},
				pushBoxOffset: hitbox,
			},
			Health:          3,
			MaxHealth:       3,
			isDead:          false,
			KnockbackFrames: 0,
		},
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       BlobIdle,
		waitFrames:  rand.Intn(MaxWaitFrames) + 1,
	}
}

func (c *BlobEnemy) ApplyKnockback(force Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != BlobDying {
		c.state = BlobHurt
		c.animations[BlobHurt].Reset()
	}
}

func (c *BlobEnemy) IsKnockedBack() bool {
	return c.KnockbackFrames > 0
}

func (c *BlobEnemy) TakeDamage(damage int) {
	if c.isDead || c.state == BlobDying || c.state == BlobHurt {
		return
	}

	// TODO: consider a state transition like Player that handles animation reset
	c.state = BlobHurt
	c.animations[BlobHurt].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = BlobDying
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *BlobEnemy) findNewTargetTile(level *Level) bool {
	// Get current tile coordinates
	tx, ty := level.WorldToTile(c.Location())

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
			c.moveStartLocation = c.Location()
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

	if c.UpdateKnockback(level) {
		// Ensure the BlobHurt animation can finish, even during knockback
		if c.state == BlobHurt && c.animations[BlobHurt].IsFinished() {
			c.state = BlobIdle
		}

		if c.state != BlobDying {
			return // Skip AI and normal movement logic
		}
	}

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
		target := Vector{
			X: c.moveTargetLocation.X - c.X,
			Y: c.moveTargetLocation.Y - c.Y,
		}

		distance := target.Length()
		if distance <= BlobMoveSpeed {
			// We are close enough to snap to the target.
			c.SetLocation(c.moveTargetLocation)

			// Wait for a short time.
			c.state = BlobIdle
			c.waitFrames = rand.Intn(MaxWaitFrames) + 1
			return
		}

		velocity := target.Normalize().Scale(BlobMoveSpeed)
		c.HandleTileCollisions(level, AxisX, velocity.X)
		c.HandleTileCollisions(level, AxisY, velocity.Y)

	case BlobAttacking:
		// For now, immediately return to idle/exploring state
		c.state = BlobIdle

	case BlobHurt:
		if c.animations[BlobHurt].IsFinished() {
			c.state = BlobIdle
		}

	case BlobDying:
		if c.animations[BlobDying].IsFinished() {
			c.isDead = true
		}
	}
}
