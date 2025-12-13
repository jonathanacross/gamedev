package main

import (
	"math"
	"math/rand"
)

type BatEnemyState int

const (
	BatFlying BatEnemyState = iota
	BatHurt
	BatDying
	BatStunned

	// Movement constants
	BatMoveSpeed         float64 = 0.8
	BatVelocityFrequency float64 = 0.05
)

type BatEnemy struct {
	BaseCharacter
	spriteSheet *SpriteSheet
	animations  map[BatEnemyState]map[Direction]*Animation

	// AI
	state              BatEnemyState
	direction          Direction
	moveStartLocation  Location
	moveTargetLocation Location
	moveTimeCounter    float64
}

func NewBatEnemy(startLoc Location) *BatEnemy {

	ssColumns := 8
	ssRows := 8
	directionOffsets := map[Direction]int{
		Right: 0,
		Left:  ssColumns,
	}
	framesPerState := ssColumns * len(directionOffsets)
	animations := map[BatEnemyState]map[Direction]*Animation{
		BatFlying:  NewDirectionAnimationMap([]int{0, 1, 2, 3, 4}, 0*framesPerState, directionOffsets, 10, true),
		BatHurt:    NewDirectionAnimationMap([]int{0, 1, 0, 1}, 2*framesPerState, directionOffsets, 10, false),
		BatDying:   NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 3*framesPerState, directionOffsets, 6, false),
		BatStunned: NewDirectionAnimationMap([]int{0, 1}, 2*framesPerState, directionOffsets, 10, true),
	}

	animations[BatFlying][Left].SetRandomFrame()
	animations[BatFlying][Right].SetRandomFrame()

	spriteSheet := NewSpriteSheet(16, 16, ssColumns, ssRows)
	hitbox := Rect{
		Left:   -9,
		Top:    -4,
		Right:  9,
		Bottom: 4,
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
			Experience:      3,
			isDead:          false,
			KnockbackFrames: 0,
		},
		spriteSheet:        spriteSheet,
		animations:         animations,
		state:              BatFlying,
		direction:          Left,
		moveStartLocation:  startLoc,
		moveTargetLocation: startLoc,
	}
}

func (c *BatEnemy) ApplyKnockback(force Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != BatDying {
		c.state = BatHurt
		c.animations[BatHurt][c.direction].Reset()
	}
}

func (c *BatEnemy) ApplyStun(duration int) {
	c.BaseCharacter.ApplyStun(duration)
	c.state = BatStunned
	c.animations[BatStunned][c.direction].Reset()
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

// updateMovement handles the movement logic and target finding.
func (c *BatEnemy) updateMovement(level *Level) {
	c.moveTimeCounter += BatVelocityFrequency
	if c.moveTimeCounter > 2*math.Pi {
		c.moveTimeCounter -= 2 * math.Pi
	}

	// Calculate Speed
	speedMultiplier := (math.Sin(c.moveTimeCounter) + 1.0) * 0.5
	currentSpeed := speedMultiplier * BatMoveSpeed

	// Calculate Movement Vector
	targetVector := Vector(c.moveTargetLocation).Minus(Vector(c.Location()))
	distance := targetVector.Length()

	if distance <= currentSpeed {
		// Arrival at Target: find a new place to go
		c.findNewTargetTile(level)
		c.moveTimeCounter = 0.0

	} else {
		// Continue Movement
		velocity := targetVector.Normalize().Scale(currentSpeed)

		// Apply movement and handle collisions
		// We use the BaseCharacter's collision handler, which modifies c.X/c.Y
		originalVelocity := velocity
		velocity.X = c.HandleTileCollisions(level, AxisX, velocity.X)
		velocity.Y = c.HandleTileCollisions(level, AxisY, velocity.Y)

		// If we hit a wall (velocity blocked), pick a new target immediately to avoid getting stuck
		if (math.Abs(originalVelocity.X) > 0.001 && velocity.X == 0) ||
			(math.Abs(originalVelocity.Y) > 0.001 && velocity.Y == 0) {
			c.findNewTargetTile(level)
		}

		// Update Direction for Animation
		if velocity.X < 0 {
			c.direction = Left
		} else if velocity.X > 0 {
			c.direction = Right
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

	if c.state != BatDying && c.UpdateStun() {
		c.state = BatStunned
		// While stunned, we just play the animation and do nothing else
		// The animation loop is set to true for BatStunned, so it will loop until stun logic ends
		return UpdateResult{Actions: actions}
	} else if c.state == BatStunned {
		// Stun finished, return to flying
		c.state = BatFlying
	}

	// Handle Dying state which overrides all AI
	if c.state == BatDying {
		if c.animations[BatDying][c.direction].IsFinished() {
			c.isDead = true
			actions = append(actions, Action{
				Type:       ActionGainXP,
				Experience: c.Experience,
			})
		}
		return UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == BatHurt {
		if !c.animations[BatHurt][c.direction].IsFinished() {
			return UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to flying state
		c.state = BatFlying
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
