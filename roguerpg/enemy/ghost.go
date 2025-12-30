package enemy

import (
	"math"
	"math/rand"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
)

type GhostEnemyState int

const (
	GhostIdle GhostEnemyState = iota
	GhostHurt
	GhostDying
	GhostMoving
	GhostStunned

	// Movement constants
	GhostMaxFloatSpeed float64 = 0.8
	GhostAttractSpeed  float64 = 0.7
	GhostXFrequency    float64 = 0.02 // radians per frame
	GhostYFrequency    float64 = 0.03 // radians per frame
)

type GhostEnemy struct {
	character.BaseCharacter
	spriteSheet *core.SpriteSheet
	animations  map[GhostEnemyState]map[core.Direction]*core.Animation

	// AI
	state         GhostEnemyState
	velocity      core.Vector
	direction     core.Direction
	movementTimer float64 // time t for Lissajous pattern
}

func NewGhostEnemy(startLoc core.Location) *GhostEnemy {
	ssColumns := 12
	ssRows := 20
	directionOffsets := map[core.Direction]int{
		core.Down:  0 * ssColumns,
		core.Up:    1 * ssColumns,
		core.Left:  2 * ssColumns,
		core.Right: 3 * ssColumns,
	}
	framesPerState := ssColumns * len(directionOffsets)

	animations := map[GhostEnemyState]map[core.Direction]*core.Animation{
		GhostIdle:    core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*framesPerState, directionOffsets, 10, true),
		GhostMoving:  core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 1*framesPerState, directionOffsets, 10, true),
		GhostHurt:    core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, false),
		GhostDying:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 4*framesPerState, directionOffsets, 8, false),
		GhostStunned: core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, true),
	}

	animations[GhostMoving][core.Up].SetRandomFrame()
	animations[GhostMoving][core.Down].SetRandomFrame()
	animations[GhostMoving][core.Left].SetRandomFrame()
	animations[GhostMoving][core.Right].SetRandomFrame()

	spriteSheet := core.NewSpriteSheet(64, 64, ssColumns, ssRows)
	hitbox := core.Rect{
		Left:   -7,
		Top:    -21,
		Right:  7,
		Bottom: 0,
	}

	return &GhostEnemy{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: startLoc,
					DrawOffset: core.Location{
						X: 32,
						Y: 32,
					},
					SrcRect: spriteSheet.Rect(0),
					Image:   assets.GhostSpritesImage,
				},
				PushBoxOffset: hitbox,
			},
			Health:          5,
			MaxHealth:       5,
			Experience:      5,
			Dead:            false,
			KnockbackFrames: 0,
		},
		spriteSheet:   spriteSheet,
		animations:    animations,
		state:         GhostMoving,
		direction:     core.Left,
		velocity:      core.Vector{X: 0, Y: 0},
		movementTimer: rand.Float64() * 100,
	}
}

func (c *GhostEnemy) HandleHit(ds *core.DamageSource) {
	if ds.Type == core.DamageTypeStun {
		c.ApplyStun(ds.Duration)
	} else {
		c.TakeDamage(ds.Damage)
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), c.Location(), core.KnockbackForce)
		c.ApplyKnockback(force, core.KnockbackDuration)
	}
}

func (c *GhostEnemy) ApplyKnockback(force core.Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != GhostDying {
		c.state = GhostHurt
		c.animations[GhostHurt][c.direction].Reset()
	}
}

func (c *GhostEnemy) ApplyStun(duration int) {
	c.BaseCharacter.ApplyStun(duration)
	c.state = GhostStunned
	c.animations[GhostStunned][c.direction].Reset()
}

func (c *GhostEnemy) TakeDamage(damage int) {
	if c.Dead || c.state == GhostDying || c.state == GhostHurt {
		return
	}

	c.state = GhostHurt
	c.animations[GhostHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = GhostDying
	}
}

// updateMovement handles the movement logic and target finding.
func (c *GhostEnemy) updateWalk(level core.Level, player core.Player) {

	c.movementTimer++

	// Lissajous pattern:
	// velocity.x = maxGhostVelocity * cos(xFrequency * t)
	// velocity.y = maxGhostVelocity * sin(yFrequency * t)
	vx := GhostMaxFloatSpeed * math.Cos(GhostXFrequency*c.movementTimer)
	vy := GhostMaxFloatSpeed * math.Sin(GhostYFrequency*c.movementTimer)
	playerAttraction := 0.0
	if c.IsNearTo(player, 7) {
		playerAttraction = 0.5
	}
	if c.IsNearTo(player, 5) {
		playerAttraction = 1.0
	}
	dirToPlayer := core.Vector(player.Location()).Minus(core.Vector(c.Location()))
	dirToPlayer = dirToPlayer.Normalize().Scale(GhostAttractSpeed * playerAttraction)
	vx += dirToPlayer.X
	vy += dirToPlayer.Y
	c.velocity = core.Vector{X: vx, Y: vy}

	// Apply movement and handle collisions
	var v core.Vector
	v.X = c.HandleTileCollisions(level, core.AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, core.AxisY, c.velocity.Y)

	// Update direction based on actual movement if significant,
	// otherwise keep facing the same way or face velocity?
	// The velocity determines direction.
	if c.velocity.Length() > 0.1 {
		c.direction = core.VectorToDirection(c.velocity)
	}
}

func (c *GhostEnemy) getActiveDamageSources(player core.Player) []*core.DamageSource {
	sources := []*core.DamageSource{}

	// Add source for hitting the ghost
	if c.IsNearTo(player, 1) {
		ds := core.NewDamageSource(core.TagEnemy, c.GetHurtBox(), core.DamageTypePhysical, 1)
		return []*core.DamageSource{ds}
	}

	return sources
}

func (c *GhostEnemy) Update(level core.Level, player core.Player) core.UpdateResult {
	var actions []core.Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.SrcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow GhostHurt animation to finish during knockback
		if c.state == GhostHurt && c.animations[GhostHurt][c.direction].IsFinished() {
			c.state = GhostMoving
		}

		if c.state != GhostDying {
			return core.UpdateResult{Actions: actions}
		}
	}

	if c.state != GhostDying && c.UpdateStun() {
		c.state = GhostStunned
		return core.UpdateResult{Actions: actions}
	} else if c.state == GhostStunned {
		c.state = GhostMoving
	}

	// Handle Dying state which overrides all AI
	if c.state == GhostDying {
		if c.animations[GhostDying][c.direction].IsFinished() {
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
	if c.state == GhostHurt {
		if !c.animations[GhostHurt][c.direction].IsFinished() {
			return core.UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to walking
		c.state = GhostMoving
		c.animations[GhostMoving][c.direction].Reset()
	}

	// Core AI Logic
	if c.state == GhostMoving {
		c.updateWalk(level, player)
	}

	for _, ds := range c.getActiveDamageSources(player) {
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return core.UpdateResult{Actions: actions}
}

func (c *GhostEnemy) CanRemove() bool {
	return false
}
