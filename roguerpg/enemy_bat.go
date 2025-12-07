package main

import (
	"math/rand"
)

type BatEnemyState int

const (
	BatFlying BatEnemyState = iota
	BatHurt
	BatDying

	// Movement constants
	BatMoveSpeed     float64 = 0.8
	BatMaxWaitFrames int     = 90 // Max 1 second wait (up to 60 frames)
)

type BatAnimationKey struct {
	State     BatEnemyState
	Direction Direction
}

type BatEnemy struct {
	BaseCharacter
	spriteSheet *SpriteSheet
	animations  map[BatEnemyState]map[Direction]*Animation

	// AI
	state     BatEnemyState
	direction Direction

	moveStartLocation  Location
	moveTargetLocation Location
	currentFrame       int // Frame counter for the current move or wait action
	waitFrames         int // Total frames to wait when idle
}

func NewBatEnemy(startLoc Location) *BatEnemy {

	animations := map[BatEnemyState]map[Direction]*Animation{
		BatFlying: {
			Left:  NewAnimation([]int{0, 1, 2, 3, 4}, 10, true),
			Right: NewAnimation([]int{6, 7, 8, 9, 10}, 10, true),
		},
		BatHurt: {
			Left:  NewAnimation([]int{24, 25, 24, 25}, 10, false),
			Right: NewAnimation([]int{30, 31, 30, 31}, 10, false),
		},
		BatDying: {
			Left:  NewAnimation([]int{24, 25, 36, 37, 38, 39, 40, 41}, 6, false),
			Right: NewAnimation([]int{30, 31, 42, 43, 44, 45, 46, 47}, 6, false),
		},
	}

	animations[BatFlying][Left].SetRandomFrame()
	animations[BatFlying][Right].SetRandomFrame()

	spriteSheet := NewSpriteSheet(16, 16, 6, 8)
	hitbox := Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	return &BatEnemy{
		BaseCharacter: BaseCharacter{
			BasePhysical: BasePhysical{
				BaseSprite: BaseSprite{
					Location: startLoc,
					drawOffset: Location{
						X: 8,
						Y: 8,
					},
					srcRect: spriteSheet.Rect(0),
					image:   BatSpritesImage,
				},
				pushBoxOffset: hitbox,
			},
			Health:          3,
			MaxHealth:       3,
			isDead:          false,
			KnockbackFrames: 0,
		},
		spriteSheet:        spriteSheet,
		animations:         animations,
		state:              BatFlying,
		direction:          Left,
		moveStartLocation:  startLoc,
		moveTargetLocation: startLoc,
		waitFrames:         rand.Intn(BatMaxWaitFrames) + 1,
	}
}

func (c *BatEnemy) ApplyKnockback(force Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != BatDying {
		c.state = BatHurt
		c.animations[BatHurt][c.direction].Reset()
	}
}

func (c *BatEnemy) IsKnockedBack() bool {
	return c.KnockbackFrames > 0
}

func (c *BatEnemy) TakeDamage(damage int) {
	if c.isDead || c.state == BatDying || c.state == BatHurt {
		return
	}

	// TODO: consider a state transition like Player that handles animation reset
	c.state = BatHurt
	c.animations[BatHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = BatDying
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *BatEnemy) findNewTargetTile(level *Level) bool {
	// Get current tile coordinates
	tx, ty := level.WorldToTile(c.Location())

	// Define a set of nearby squares to check
	nearbySquares := []Point{}
	radius := 3
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			nearbySquares = append(nearbySquares, Point{dx, dy})
		}
	}

	// Shuffle nearbySquares to pick a random one first
	rand.Shuffle(len(nearbySquares), func(i, j int) {
		nearbySquares[i], nearbySquares[j] = nearbySquares[j], nearbySquares[i]
	})

	for _, offset := range nearbySquares {
		newTx := tx + offset.X
		newTy := ty + offset.Y

		if !level.IsTileSolid(newTx, newTy) {
			// Found an open tile. Set up the movement.
			c.moveStartLocation = c.Location()
			c.moveTargetLocation = level.TileToWorld(newTx, newTy)
			c.currentFrame = 0 // Reset frame counter for movement

			if c.moveTargetLocation.X < c.Location().X {
				c.direction = Left
			} else if c.moveTargetLocation.X > c.Location().X {
				c.direction = Right
			}
			return true
		}
	}

	// No adjacent open tile found
	return false
}

func (c *BatEnemy) isNearPlayer(playerLoc Location) bool {
	dist := Vector(playerLoc).Minus(Vector(c.Location()))
	return dist.Length() <= TileSize
}

// updateActionState handles the transition between BatAttacking and BatFlying
// func (c *BatEnemy) updateActionState(playerLoc Location) {
// 	isPlayerNear := c.isNearPlayer(playerLoc)

// 	if isPlayerNear {
// 		// Priority 1: Start Attack if near player and not already attacking
// 		if c.state != BatAttacking {
// 			c.state = BatAttacking
// 			c.animations[BatAttacking][c.direction].Reset()
// 		}
// 	} else if c.state == BatAttacking && c.animations[BatAttacking][c.direction].IsFinished() {
// 		// Priority 2: Finish Attack animation, return to flying
// 		c.state = BatFlying
// 	} else if c.state != BatAttacking {
// 		// Priority 3: Maintain flying state if not near player and not currently finishing an attack
// 		c.state = BatFlying
// 	}
// 	// Note: If c.state is BatAttacking but the animation isn't finished, we stay in BatAttacking.
// 	// If c.state is BatHurt/BatDying, this function is skipped by the caller.
// }

// updateMovement handles the movement logic and target finding.
func (c *BatEnemy) updateMovement(level *Level) {
	// Check if movement is currently active
	targetVector := Vector(c.moveTargetLocation).Minus(Vector(c.Location()))
	isMovingActive := targetVector.Length() > BatMoveSpeed

	if isMovingActive {
		// Execute movement
		distance := targetVector.Length()

		if distance <= BatMoveSpeed {
			// We are close enough to snap to the target.
			c.SetLocation(c.moveTargetLocation)

			// --- FIX: Invalidate the target to end the movement loop ---
			c.moveTargetLocation = c.Location()
			// -----------------------------------------------------------

			// Movement finished, start waiting.
			c.waitFrames = rand.Intn(BatMaxWaitFrames) + 1
		} else {
			// Continue movement
			velocity := targetVector.Normalize().Scale(BatMoveSpeed)
			c.HandleTileCollisions(level, AxisX, velocity.X)
			c.HandleTileCollisions(level, AxisY, velocity.Y)
		}
	} else {
		// Idle/Wait logic
		c.waitFrames--
		if c.waitFrames <= 0 {
			// Wait time is over. Look for a new target tile.
			if !c.findNewTargetTile(level) {
				// Enemy is cornered or blocked. Wait again.
				c.waitFrames = rand.Intn(BatMaxWaitFrames) + 1
			}
			// If a new target is found, isMovingActive will be true on the next frame.
		}
	}
}

func (c *BatEnemy) Update(level *Level, player *Player) UpdateResult {
	var actions []Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.srcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow BatHurt animation to finish during knockback
		if c.state == BatHurt && c.animations[BatHurt][c.direction].IsFinished() {
			c.state = BatFlying
		}

		if c.state != BatDying {
			return UpdateResult{Actions: actions}
		}
	}

	// Handle Dying state which overrides all AI
	if c.state == BatDying {
		if c.animations[BatDying][c.direction].IsFinished() {
			c.isDead = true
		}
		return UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == BatHurt {
		if !c.animations[BatHurt][c.direction].IsFinished() {
			return UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}
	}

	// Core AI Logic
	c.updateMovement(level)

	if c.isNearPlayer(player.Location()) && c.state == BatFlying {
		ds := NewDamageSource(TagEnemy, c.GetHurtBox(), 1)
		actions = append(actions, Action{
			Type:         ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return UpdateResult{Actions: actions}
}

func (c *BatEnemy) CanRemove() bool {
	// TODO: consider if we want to merge this with isdead.
	return false
}
