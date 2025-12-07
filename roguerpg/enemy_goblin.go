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
	GoblinWalkSpeed  float64 = 0.6
	GoblinTurnFrames int     = 60 * 3 // every 3 seconds
)

type GoblinEnemy struct {
	BaseCharacter
	spriteSheet *SpriteSheet
	animations  map[GoblinEnemyState]map[Direction]*Animation

	// AI
	state    GoblinEnemyState
	velocity Vector

	direction       Direction
	turnTimeCounter int
}

func getRandomVelocity() Vector {
	velocities := []Vector{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}
	return velocities[rand.Intn(len(velocities))]
}

func NewGoblinEnemy(startLoc Location) *GoblinEnemy {
	ssColumns := 8
	ssRows := 32
	directionOffsets := map[Direction]int{
		Down:  0,
		Up:    8,
		Left:  16,
		Right: 24,
	}
	framesPerState := ssColumns * len(directionOffsets)

	animations := map[GoblinEnemyState]map[Direction]*Animation{
		GoblinIdle:      NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*framesPerState, directionOffsets, 10, true),
		GoblinWalking:   NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 1*framesPerState, directionOffsets, 10, true),
		GoblinRunning:   NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 2*framesPerState, directionOffsets, 10, true),
		GoblinAttacking: NewDirectionAnimationMap([]int{0, 1, 2, 3, 4}, 3*framesPerState, directionOffsets, 10, false),
		GoblinHurt:      NewDirectionAnimationMap([]int{0, 1, 2, 3}, 6*framesPerState, directionOffsets, 10, false),
		GoblinDying:     NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 7*framesPerState, directionOffsets, 10, false),
	}

	animations[GoblinWalking][Up].SetRandomFrame()
	animations[GoblinWalking][Down].SetRandomFrame()
	animations[GoblinWalking][Left].SetRandomFrame()
	animations[GoblinWalking][Right].SetRandomFrame()

	spriteSheet := NewSpriteSheet(64, 64, ssColumns, ssRows)
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
		spriteSheet:     spriteSheet,
		animations:      animations,
		state:           GoblinWalking,
		direction:       Left,
		velocity:        getRandomVelocity(),
		turnTimeCounter: rand.Intn(GoblinTurnFrames),
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

func (c *GoblinEnemy) isNearPlayer(playerLoc Location) bool {
	dist := Vector(playerLoc).Minus(Vector(c.Location()))
	return dist.Length() <= TileSize
}

// updateMovement handles the movement logic and target finding.
func (c *GoblinEnemy) updateMovement(level *Level) {
	c.turnTimeCounter++
	if c.turnTimeCounter >= GoblinTurnFrames {
		c.turnTimeCounter = 0
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	// Apply movement and handle collisions
	// Use the BaseCharacter's collision handler, which modifies c.X/c.Y
	var v Vector
	v.X = c.HandleTileCollisions(level, AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, AxisY, c.velocity.Y)
	if v.Length() < 0.01 {
		// Turn immediately if we run into something; don't reset the turn counter, though
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	c.direction = GetDirection(v)
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
			c.state = GoblinWalking
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

		// Hurt animation finished, return to walking
		c.state = GoblinWalking
	}

	// Core AI Logic
	c.updateMovement(level)

	if c.isNearPlayer(player.Location()) && c.state == GoblinWalking {
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
