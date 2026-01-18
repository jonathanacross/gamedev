package enemy

import (
	"math"
	"math/rand"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
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

	BatHp  int = 4
	BatExp int = 4
)

type BatEnemy struct {
	character.BaseCharacter
	spriteSheet *core.SpriteSheet
	animations  map[BatEnemyState]map[core.Direction]*core.Animation

	// AI
	state              BatEnemyState
	direction          core.Direction
	moveStartLocation  core.Location
	moveTargetLocation core.Location
	moveTimeCounter    float64
}

func NewBatEnemy(startLoc core.Location) *BatEnemy {

	ssColumns := 8
	ssRows := 8
	directionOffsets := map[core.Direction]int{
		core.Right: 0,
		core.Left:  ssColumns,
	}
	framesPerState := ssColumns * len(directionOffsets)
	animations := map[BatEnemyState]map[core.Direction]*core.Animation{
		BatFlying:  core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4}, 0*framesPerState, directionOffsets, 10, true),
		BatHurt:    core.NewDirectionAnimationMap([]int{0, 1, 0, 1}, 2*framesPerState, directionOffsets, 10, false),
		BatDying:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 3*framesPerState, directionOffsets, 6, false),
		BatStunned: core.NewDirectionAnimationMap([]int{0, 1}, 2*framesPerState, directionOffsets, 10, true),
	}

	animations[BatFlying][core.Left].SetRandomFrame()
	animations[BatFlying][core.Right].SetRandomFrame()

	spriteSheet := core.NewSpriteSheet(16, 16, ssColumns, ssRows)
	hitbox := core.Rect{
		Left:   -9,
		Top:    -4,
		Right:  9,
		Bottom: 4,
	}

	return &BatEnemy{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: startLoc,
					DrawOffset: core.Location{
						X: 8,
						Y: 8,
					},
					SrcRect: spriteSheet.Rect(0),
					Image:   assets.BatSpritesImage,
				},
				PushBoxOffset: hitbox,
			},
			Health:          BatHp,
			MaxHealth:       BatHp,
			Experience:      BatExp,
			Dead:            false,
			KnockbackFrames: 0,
		},
		spriteSheet:        spriteSheet,
		animations:         animations,
		state:              BatFlying,
		direction:          core.Left,
		moveStartLocation:  startLoc,
		moveTargetLocation: startLoc,
	}
}

func (c *BatEnemy) ApplyKnockback(force core.Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != BatDying {
		c.state = BatHurt
		c.animations[BatHurt][c.direction].Reset()
	}
}

func (c *BatEnemy) HandleHit(ds *core.DamageSource) {
	if ds.Type == core.DamageTypeStun {
		c.ApplyStun(ds.Duration)
	} else {
		c.TakeDamage(ds.Damage)
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), c.Location(), core.KnockbackForce)
		c.ApplyKnockback(force, core.KnockbackDuration)
	}
}

func (c *BatEnemy) ApplyStun(duration int) {
	c.BaseCharacter.ApplyStun(duration)
	c.state = BatStunned
	c.animations[BatStunned][c.direction].Reset()
}

func (c *BatEnemy) TakeDamage(damage int) {
	if c.Dead || c.state == BatDying || c.state == BatHurt {
		return
	}

	c.state = BatHurt
	c.animations[BatHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = BatDying
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *BatEnemy) findNewTargetTile(level core.Level) bool {
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

// updateMovement handles the movement logic and target finding.
func (c *BatEnemy) updateMovement(level core.Level) {
	c.moveTimeCounter += BatVelocityFrequency
	if c.moveTimeCounter > 2*math.Pi {
		c.moveTimeCounter -= 2 * math.Pi
	}

	// Calculate Speed
	speedMultiplier := (math.Sin(c.moveTimeCounter) + 1.0) * 0.5
	currentSpeed := speedMultiplier * BatMoveSpeed

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

func (c *BatEnemy) Update(level core.Level, player core.Player) core.UpdateResult {
	var actions []core.Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.SrcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow BatHurt animation to finish during knockback
		if c.state == BatHurt && c.animations[BatHurt][c.direction].IsFinished() {
			c.state = BatFlying
		}

		if c.state != BatDying {
			return core.UpdateResult{Actions: actions}
		}
	}

	if c.state != BatDying && c.UpdateStun() {
		c.state = BatStunned
		// While stunned, we just play the animation and do nothing else
		// The animation loop is set to true for BatStunned, so it will loop until stun logic ends
		return core.UpdateResult{Actions: actions}
	} else if c.state == BatStunned {
		// Stun finished, return to flying
		c.state = BatFlying
	}

	// Handle Dying state which overrides all AI
	if c.state == BatDying {
		if c.animations[BatDying][c.direction].IsFinished() {
			c.Dead = true
			actions = append(actions, core.Action{
				Type:       core.ActionGainXP,
				Experience: c.Experience,
			})
			actions = append(actions, maybeDropHeart(c.Location())...)
		}
		return core.UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == BatHurt {
		if !c.animations[BatHurt][c.direction].IsFinished() {
			return core.UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to flying state
		c.state = BatFlying
	}

	// Core AI Logic
	c.updateMovement(level)

	if c.IsNearTo(player, 3) && c.state == BatFlying {
		ds := core.NewDamageSource(core.TagEnemy, c.GetHurtBox(), core.DamageTypePhysical, 1)
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return core.UpdateResult{Actions: actions}
}

func (c *BatEnemy) CanRemove() bool {
	return false
}
