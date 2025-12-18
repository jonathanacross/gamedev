package enemy

import (
	"math"
	"math/rand"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
	"time"
)

type SpikeTurtleEnemyState int

const (
	SpikeTurtleMoving  SpikeTurtleEnemyState = iota
	SpikeTurtleFlipped SpikeTurtleEnemyState = iota
	SpikeTurtleHurt
	SpikeTurtleDying

	// Movement constants
	SpikeTurtleMoveSpeed         float64 = 0.8
	SpikeTurtleVelocityFrequency float64 = 0.05

	SpikeTurtleFlipDuration time.Duration = 4000 * time.Millisecond
)

type SpikeTurtleEnemy struct {
	character.BaseCharacter
	spriteSheet *core.SpriteSheet
	animations  map[SpikeTurtleEnemyState]*core.Animation

	// AI
	state              SpikeTurtleEnemyState
	direction          core.Direction
	moveStartLocation  core.Location
	moveTargetLocation core.Location
	moveTimeCounter    float64
	flipTimer          *core.Timer
}

func NewSpikeTurtleEnemy(startLoc core.Location) *SpikeTurtleEnemy {

	ssColumns := 4
	ssRows := 4
	animations := map[SpikeTurtleEnemyState]*core.Animation{
		SpikeTurtleMoving:  core.NewAnimation([]int{0, 1}, 10, true),
		SpikeTurtleFlipped: core.NewAnimation([]int{4, 5}, 10, true),
		SpikeTurtleHurt:    core.NewAnimation([]int{8, 10, 9, 11}, 10, false),
		SpikeTurtleDying:   core.NewAnimation([]int{12, 13, 14, 15}, 10, false),
	}

	animations[SpikeTurtleMoving].SetRandomFrame()

	spriteSheet := core.NewSpriteSheet(16, 32, ssColumns, ssRows)
	hitbox := core.Rect{
		Left:   -8,
		Top:    -8,
		Right:  8,
		Bottom: 7,
	}

	return &SpikeTurtleEnemy{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: startLoc,
					DrawOffset: core.Location{
						X: 8,
						Y: 11,
					},
					SrcRect: spriteSheet.Rect(0),
					Image:   assets.SpikeTurleSpritesImage,
				},
				PushBoxOffset: hitbox,
			},
			Health:          3,
			MaxHealth:       3,
			Experience:      3,
			Dead:            false,
			KnockbackFrames: 0,
		},
		spriteSheet:        spriteSheet,
		animations:         animations,
		state:              SpikeTurtleMoving,
		direction:          core.Left,
		moveStartLocation:  startLoc,
		moveTargetLocation: startLoc,
		flipTimer:          core.NewTimer(SpikeTurtleFlipDuration),
	}
}

func (c *SpikeTurtleEnemy) HandleHit(ds *core.DamageSource) {
	if c.Dead || c.state == SpikeTurtleDying {
		return
	}
	// Turtles are immune to stun damage
	if ds.Type == core.DamageTypeStun {
		return
	}

	if c.state == SpikeTurtleMoving {
		// If the turtle is right side up, flip it, only for impact damage
		if ds.Type == core.DamageTypeImpact {
			c.state = SpikeTurtleFlipped
			c.animations[SpikeTurtleFlipped].Reset()
		}
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), c.Location(), core.KnockbackForce)
		c.ApplyKnockback(force, core.KnockbackDuration)
	} else if c.state == SpikeTurtleFlipped {
		if ds.Type != core.DamageTypeStun {
			c.TakeDamage(ds.Damage)
			force := core.CalculateKnockbackForce(ds.HitBox.Center(), c.Location(), core.KnockbackForce)
			c.ApplyKnockback(force, core.KnockbackDuration)
		}
	}
}

func (c *SpikeTurtleEnemy) ApplyKnockback(force core.Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)
}

func (c *SpikeTurtleEnemy) ApplyStun(duration int) {
	// Can't be stunned
}

func (c *SpikeTurtleEnemy) TakeDamage(damage int) {
	if c.Dead || c.state == SpikeTurtleDying || c.state == SpikeTurtleHurt {
		return
	}

	c.state = SpikeTurtleHurt
	c.animations[SpikeTurtleHurt].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = SpikeTurtleDying
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *SpikeTurtleEnemy) findNewTargetTile(level core.Level) bool {
	// Get current tile coordinates
	tx, ty := level.WorldToTile(c.Location())

	// Define a set of nearby squares to check
	nearbySquares := []core.Point{}
	radius := 3
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			nearbySquares = append(nearbySquares, core.Point{X: dx, Y: dy})
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
				c.direction = core.Left
			} else if c.moveTargetLocation.X > c.Location().X {
				c.direction = core.Right
			}
			return true
		}
	}

	// No adjacent open tile found
	return false
}

func (c *SpikeTurtleEnemy) isNearPlayer(playerLoc core.Location) bool {
	dist := core.Vector(playerLoc).Minus(core.Vector(c.Location()))
	return dist.Length() <= core.TileSize
}

// updateMovement handles the movement logic and target finding.
func (c *SpikeTurtleEnemy) updateMovement(level core.Level) {
	c.moveTimeCounter += SpikeTurtleVelocityFrequency
	if c.moveTimeCounter > 2*math.Pi {
		c.moveTimeCounter -= 2 * math.Pi
	}

	// Calculate Speed
	speedMultiplier := (math.Sin(c.moveTimeCounter) + 1.0) * 0.5
	currentSpeed := speedMultiplier * SpikeTurtleMoveSpeed

	// Calculate Movement Vector
	targetVector := core.Vector(c.moveTargetLocation).Minus(core.Vector(c.Location()))
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
		velocity.X = c.HandleTileCollisions(level, core.AxisX, velocity.X)
		velocity.Y = c.HandleTileCollisions(level, core.AxisY, velocity.Y)

		// If we hit a wall (velocity blocked), pick a new target immediately to avoid getting stuck
		if (math.Abs(originalVelocity.X) > 0.001 && velocity.X == 0) ||
			(math.Abs(originalVelocity.Y) > 0.001 && velocity.Y == 0) {
			c.findNewTargetTile(level)
		}

		// Update Direction for Animation
		if velocity.X < 0 {
			c.direction = core.Left
		} else if velocity.X > 0 {
			c.direction = core.Right
		}
	}
}

func (c *SpikeTurtleEnemy) Update(level core.Level, player core.Player) core.UpdateResult {
	var actions []core.Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state]
	currentAnim.Update()
	c.SrcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow SpikeTurtleHurt animation to finish during knockback
		if c.state == SpikeTurtleHurt && c.animations[SpikeTurtleHurt].IsFinished() {
			c.state = SpikeTurtleMoving
		}

		if c.state != SpikeTurtleDying {
			return core.UpdateResult{Actions: actions}
		}
	}

	// Handle Dying state which overrides all AI
	if c.state == SpikeTurtleDying {
		if c.animations[SpikeTurtleDying].IsFinished() {
			c.Dead = true
			actions = append(actions, core.Action{
				Type:       core.ActionGainXP,
				Experience: c.Experience,
			})
		}
		return core.UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == SpikeTurtleHurt {
		if !c.animations[SpikeTurtleHurt].IsFinished() {
			return core.UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to flying state
		c.state = SpikeTurtleMoving
	}

	if c.state == SpikeTurtleFlipped {
		c.flipTimer.Update()
		if c.flipTimer.IsReady() {
			c.flipTimer.Reset()
			c.state = SpikeTurtleMoving
			c.animations[SpikeTurtleMoving].Reset()
		}
		return core.UpdateResult{Actions: actions}
	}

	// Core AI Logic
	c.updateMovement(level)

	if c.isNearPlayer(player.Location()) && c.state == SpikeTurtleMoving {
		ds := core.NewDamageSource(core.TagEnemy, c.GetHurtBox(), core.DamageTypePhysical, 1)
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return core.UpdateResult{Actions: actions}
}

func (c *SpikeTurtleEnemy) CanRemove() bool {
	return false
}
