package main

import (
	"math/rand"
)

type BatEnemyState int

const (
	BatIdle BatEnemyState = iota
	BatMoving
	BatAttacking
	BatHurt
	BatDying

	// Movement constants
	BatMoveSpeed     float64 = 0.5
	BatMaxWaitFrames int     = 60 // Max 1 second wait (up to 60 frames)
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

func NewBatEnemy() *BatEnemy {

	animations := map[BatEnemyState]map[Direction]*Animation{
		BatIdle: {
			Left:  NewAnimation([]int{0, 1, 2, 3, 4}, 10, true),
			Right: NewAnimation([]int{6, 7, 8, 9, 10}, 10, true),
		},
		BatMoving: {
			Left:  NewAnimation([]int{0, 1, 2, 3, 4}, 10, true),
			Right: NewAnimation([]int{6, 7, 8, 9, 10}, 10, true),
		},
		BatAttacking: {
			Left:  NewAnimation([]int{12, 13, 14}, 10, true),
			Right: NewAnimation([]int{18, 19, 20}, 10, true),
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

	animations[BatIdle][Left].SetRandomFrame()
	animations[BatIdle][Right].SetRandomFrame()

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
					Location: Location{
						X: 0,
						Y: 0,
					},
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
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       BatIdle,
		direction:   Left,
		waitFrames:  rand.Intn(MaxWaitFrames) + 1,
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

func (c *BatEnemy) isNearPlayer(playerLoc Location) bool {
	dist := Vector(playerLoc).Minus(Vector(c.Location()))
	return dist.Length() <= 24
}

func (c *BatEnemy) Update(level *Level, player *Player) UpdateResult {
	c.animations[c.state][c.direction].Update()
	c.srcRect = c.spriteSheet.Rect(c.animations[c.state][c.direction].Frame())
	var actions []Action

	if c.UpdateKnockback(level) {
		// Ensure the BatHurt animation can finish, even during knockback
		if c.state == BatHurt && c.animations[BatHurt][c.direction].IsFinished() {
			c.state = BatIdle
		}

		if c.state != BatDying {
			return UpdateResult{Actions: actions} // Skip AI and normal movement logic
		}
	}

	switch c.state {
	case BatIdle:
		if c.isNearPlayer(player.Location()) {
			c.state = BatAttacking
			c.animations[BatAttacking][c.direction].Reset()
		} else {
			c.waitFrames--
			if c.waitFrames <= 0 {
				// Wait time is over. Look for a new target tile.
				if c.findNewTargetTile(level) {
					c.state = BatMoving
				} else {
					// Enemy is cornered or blocked. Wait again.
					c.waitFrames = rand.Intn(MaxWaitFrames) + 1
				}
			}
		}

	case BatMoving:
		if c.isNearPlayer(player.Location()) {
			c.state = BatAttacking
			c.animations[BatAttacking][c.direction].Reset()
		} else {
			target := Vector(c.moveTargetLocation).Minus(Vector(c.Location()))

			distance := target.Length()
			if distance <= BatMoveSpeed {
				// We are close enough to snap to the target.
				c.SetLocation(c.moveTargetLocation)

				// Wait for a short time.
				c.state = BatIdle
				c.waitFrames = rand.Intn(MaxWaitFrames) + 1
				return UpdateResult{Actions: actions}
			}

			velocity := target.Normalize().Scale(BatMoveSpeed)
			c.HandleTileCollisions(level, AxisX, velocity.X)
			c.HandleTileCollisions(level, AxisY, velocity.Y)
		}

	case BatAttacking:
		hitBox := c.GetHurtBox()
		if c.animations[BatAttacking][c.direction].Frame() != 0 {
			ds := NewDamageSource(TagEnemy, hitBox, 1)
			actions = append(actions, Action{
				Type:         ActionCreateDamageSource,
				DamageSource: ds,
			})
		}

		if c.animations[BatAttacking][c.direction].IsFinished() {
			c.state = BatIdle
			c.waitFrames = rand.Intn(MaxWaitFrames) + 1
		}

	case BatHurt:
		if c.animations[BatHurt][c.direction].IsFinished() {
			c.state = BatIdle
		}

	case BatDying:
		if c.animations[BatDying][c.direction].IsFinished() {
			c.isDead = true
		}
	}

	return UpdateResult{Actions: actions}
}

func (c *BatEnemy) CanRemove() bool {
	// TODO: consider if we want to merge this with isdead.
	return false
}
