package main

import (
	"math"
	"math/rand"
)

type GoblinEnemyState int

const (
	GoblinIdle GoblinEnemyState = iota
	GoblinHurt
	GoblinDying
	GoblinAttacking
	GoblinRunning
	GoblinWalking

	// Movement constants
	GoblinMoveSpeed         float64 = 0.8
	GoblinVelocityFrequency float64 = 0.05
)

type GoblinEnemy struct {
	BaseCharacter
	spriteSheet *SpriteSheet
	animations  map[GoblinEnemyState]map[Direction]*Animation

	// AI
	state              GoblinEnemyState
	direction          Direction
	moveStartLocation  Location
	moveTargetLocation Location
	moveTimeCounter    float64
}

func NewDirectionAnimationMap(frames []int, stateStartFrame int, directionOffsets map[Direction]int, speed int, looping bool) map[Direction]*Animation {
	animations := map[Direction]*Animation{}
	for direction, dirOffset := range directionOffsets {
		offsetFrames := make([]int, len(frames))
		for i := range frames {
			offsetFrames[i] = frames[i] + dirOffset + stateStartFrame
		}
		animations[direction] = NewAnimation(offsetFrames, speed, looping)
	}
	return animations
}

func NewGoblinEnemy(startLoc Location) *GoblinEnemy {

	directionOffsets := map[Direction]int{
		Down:  0,
		Up:    8,
		Left:  16,
		Right: 24,
	}

	animations := map[GoblinEnemyState]map[Direction]*Animation{
		GoblinIdle:      NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*32, directionOffsets, 10, true),
		GoblinWalking:   NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 1*32, directionOffsets, 10, true),
		GoblinRunning:   NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 2*32, directionOffsets, 10, true),
		GoblinAttacking: NewDirectionAnimationMap([]int{0, 1, 2, 3, 4}, 3*32, directionOffsets, 10, false),
		GoblinHurt:      NewDirectionAnimationMap([]int{0, 1, 2, 3}, 6*32, directionOffsets, 10, false),
		GoblinDying:     NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 7*32, directionOffsets, 10, false),
	}

	animations[GoblinIdle][Up].SetRandomFrame()
	animations[GoblinIdle][Down].SetRandomFrame()
	animations[GoblinIdle][Left].SetRandomFrame()
	animations[GoblinIdle][Right].SetRandomFrame()

	spriteSheet := NewSpriteSheet(64, 64, 8, 32)
	hitbox := Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	return &GoblinEnemy{
		BaseCharacter: BaseCharacter{
			BasePhysical: BasePhysical{
				BaseSprite: BaseSprite{
					Location: startLoc,
					drawOffset: Location{
						X: 32,
						Y: 32,
					},
					srcRect: spriteSheet.Rect(0),
					image:   GoblinSpritesImage,
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
		state:              GoblinIdle,
		direction:          Left,
		moveStartLocation:  startLoc,
		moveTargetLocation: startLoc,
	}
}

func (c *GoblinEnemy) ApplyKnockback(force Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != GoblinDying {
		c.state = GoblinHurt
		c.animations[GoblinHurt][c.direction].Reset()
	}
}

func (c *GoblinEnemy) IsKnockedBack() bool {
	return c.KnockbackFrames > 0
}

func (c *GoblinEnemy) TakeDamage(damage int) {
	if c.isDead || c.state == GoblinDying || c.state == GoblinHurt {
		return
	}

	// TODO: consider a state transition like Player that handles animation reset
	c.state = GoblinHurt
	c.animations[GoblinHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = GoblinDying
	}
}

func GetDirection(dirVector Vector) Direction {
	angle := math.Atan2(dirVector.Y, dirVector.X)
	if angle >= -math.Pi/4 && angle < math.Pi/4 {
		return Right
	} else if angle >= math.Pi/4 && angle < 3*math.Pi/4 {
		return Down
	} else if angle >= -3*math.Pi/4 && angle < -math.Pi/4 {
		return Up
	} else {
		return Left
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *GoblinEnemy) findNewTargetTile(level *Level) bool {
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

			dirVector := Vector(c.moveTargetLocation).Minus(Vector(c.Location()))
			c.direction = GetDirection(dirVector)
			return true
		}
	}

	// No adjacent open tile found
	return false
}

func (c *GoblinEnemy) isNearPlayer(playerLoc Location) bool {
	dist := Vector(playerLoc).Minus(Vector(c.Location()))
	return dist.Length() <= TileSize
}

// updateMovement handles the movement logic and target finding.
func (c *GoblinEnemy) updateMovement(level *Level) {
	c.moveTimeCounter += GoblinVelocityFrequency
	if c.moveTimeCounter > 2*math.Pi {
		c.moveTimeCounter -= 2 * math.Pi
	}

	// Calculate Speed
	speedMultiplier := (math.Sin(c.moveTimeCounter) + 1.0) * 0.5
	currentSpeed := speedMultiplier * GoblinMoveSpeed

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
		velocity.X = c.HandleTileCollisions(level, AxisX, velocity.X)
		velocity.Y = c.HandleTileCollisions(level, AxisY, velocity.Y)

		// Update Direction for Animation
		c.direction = GetDirection(velocity)
		// if velocity.X < 0 {
		// 	c.direction = Left
		// } else if velocity.X > 0 {
		// 	c.direction = Right
		// }
	}
}

func (c *GoblinEnemy) Update(level *Level, player *Player) UpdateResult {
	var actions []Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.srcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow GoblinHurt animation to finish during knockback
		if c.state == GoblinHurt && c.animations[GoblinHurt][c.direction].IsFinished() {
			c.state = GoblinIdle
		}

		if c.state != GoblinDying {
			return UpdateResult{Actions: actions}
		}
	}

	// Handle Dying state which overrides all AI
	if c.state == GoblinDying {
		if c.animations[GoblinDying][c.direction].IsFinished() {
			c.isDead = true
		}
		return UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == GoblinHurt {
		if !c.animations[GoblinHurt][c.direction].IsFinished() {
			return UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to flying state
		c.state = GoblinIdle
	}

	// Core AI Logic
	c.updateMovement(level)

	if c.isNearPlayer(player.Location()) && c.state == GoblinIdle {
		ds := NewDamageSource(TagEnemy, c.GetHurtBox(), 1)
		actions = append(actions, Action{
			Type:         ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return UpdateResult{Actions: actions}
}

func (c *GoblinEnemy) CanRemove() bool {
	// TODO: consider if we want to merge this with isdead.
	return false
}
