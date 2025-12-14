package enemy

import (
	"math"
	"math/rand"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
)

type GoblinEnemyState int

const (
	GoblinIdle GoblinEnemyState = iota
	GoblinHurt
	GoblinDying
	GoblinAttacking
	GoblinRunning
	GoblinWalking
	GoblinStunned

	// Movement constants
	GoblinWalkSpeed  float64 = 0.6
	GoblinRunSpeed   float64 = 1.2
	GoblinTurnFrames int     = 60 * 3 // every 3 seconds
)

type GoblinEnemy struct {
	character.BaseCharacter
	spriteSheet    *core.SpriteSheet
	animations     map[GoblinEnemyState]map[core.Direction]*core.Animation
	attackHitboxes map[core.Direction]map[int]core.DamageSourceConfig

	// AI
	state           GoblinEnemyState
	velocity        core.Vector
	direction       core.Direction
	turnTimeCounter int
}

func getRandomVelocity() core.Vector {
	velocities := []core.Vector{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}
	return velocities[rand.Intn(len(velocities))]
}

func NewGoblinEnemy(startLoc core.Location) *GoblinEnemy {
	ssColumns := 8
	ssRows := 32
	directionOffsets := map[core.Direction]int{
		core.Down:  0 * ssColumns,
		core.Up:    1 * ssColumns,
		core.Left:  2 * ssColumns,
		core.Right: 3 * ssColumns,
	}
	framesPerState := ssColumns * len(directionOffsets)

	animations := map[GoblinEnemyState]map[core.Direction]*core.Animation{
		GoblinIdle:      core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*framesPerState, directionOffsets, 10, true),
		GoblinWalking:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 1*framesPerState, directionOffsets, 10, true),
		GoblinRunning:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 2*framesPerState, directionOffsets, 10, true),
		GoblinAttacking: core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4}, 3*framesPerState, directionOffsets, 10, true),
		GoblinHurt:      core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 6*framesPerState, directionOffsets, 10, false),
		GoblinDying:     core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 7*framesPerState, directionOffsets, 10, false),
		GoblinStunned:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 6*framesPerState, directionOffsets, 10, true),
	}

	animations[GoblinWalking][core.Up].SetRandomFrame()
	animations[GoblinWalking][core.Down].SetRandomFrame()
	animations[GoblinWalking][core.Left].SetRandomFrame()
	animations[GoblinWalking][core.Right].SetRandomFrame()

	spriteSheet := core.NewSpriteSheet(64, 64, ssColumns, ssRows)
	hitbox := core.Rect{
		Left:   -7,
		Top:    -7,
		Right:  7,
		Bottom: 7,
	}
	// Define a simple attack hitbox that's only active on the 2nd and 3rd frames (index 1 and 2 in the short animation array)
	attackHitboxes := make(map[core.Direction]map[int]core.DamageSourceConfig)
	baseDmg := 1

	attackHitboxes[core.Up] = map[int]core.DamageSourceConfig{
		2: {HitBox: core.Rect{Left: 3, Top: -24, Right: 13, Bottom: 4}, Damage: baseDmg},
		3: {HitBox: core.Rect{Left: 3, Top: -27, Right: 13, Bottom: 4}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: 3, Top: -27, Right: 13, Bottom: 4}, Damage: baseDmg},
	}
	attackHitboxes[core.Down] = map[int]core.DamageSourceConfig{
		2: {HitBox: core.Rect{Left: -13, Top: -4, Right: -3, Bottom: 24}, Damage: baseDmg},
		3: {HitBox: core.Rect{Left: -13, Top: -4, Right: -3, Bottom: 27}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -13, Top: -4, Right: -3, Bottom: 27}, Damage: baseDmg},
	}
	attackHitboxes[core.Left] = map[int]core.DamageSourceConfig{
		2: {HitBox: core.Rect{Left: -25, Top: 0, Right: 6, Bottom: 10}, Damage: baseDmg},
		3: {HitBox: core.Rect{Left: -28, Top: 0, Right: 6, Bottom: 10}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -28, Top: 0, Right: 6, Bottom: 10}, Damage: baseDmg},
	}
	attackHitboxes[core.Right] = map[int]core.DamageSourceConfig{
		2: {HitBox: core.Rect{Left: -6, Top: 0, Right: 25, Bottom: 10}, Damage: baseDmg},
		3: {HitBox: core.Rect{Left: -6, Top: 0, Right: 28, Bottom: 10}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -6, Top: 0, Right: 28, Bottom: 10}, Damage: baseDmg},
	}

	return &GoblinEnemy{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: startLoc,
					DrawOffset: core.Location{
						X: 32,
						Y: 32,
					},
					SrcRect: spriteSheet.Rect(0),
					Image:   assets.GoblinSpritesImage,
				},
				PushBoxOffset: hitbox,
			},
			Health:          5,
			MaxHealth:       5,
			Experience:      6,
			Dead:            false,
			KnockbackFrames: 0,
		},
		spriteSheet:     spriteSheet,
		animations:      animations,
		attackHitboxes:  attackHitboxes,
		state:           GoblinWalking,
		direction:       core.Left,
		velocity:        getRandomVelocity(),
		turnTimeCounter: rand.Intn(GoblinTurnFrames),
	}
}

func (c *GoblinEnemy) ApplyKnockback(force core.Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != GoblinDying {
		c.state = GoblinHurt
		c.animations[GoblinHurt][c.direction].Reset()
	}
}

func (c *GoblinEnemy) ApplyStun(duration int) {
	c.BaseCharacter.ApplyStun(duration)
	c.state = GoblinStunned
	c.animations[GoblinStunned][c.direction].Reset()
}

func (c *GoblinEnemy) TakeDamage(damage int) {
	if c.Dead || c.state == GoblinDying || c.state == GoblinHurt {
		return
	}

	c.state = GoblinHurt
	c.animations[GoblinHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = GoblinDying
	}
}

func (c *GoblinEnemy) shouldAttackPlayer(player core.Player) bool {
	// Player should be near the goblin, and aligned either horizontally or vertically
	dist := core.Vector(player.Location()).Minus(core.Vector(c.Location()))
	alignedHorizontally := math.Abs(dist.Y) < float64(core.TileSize)/2
	alignedVertically := math.Abs(dist.X) < float64(core.TileSize)/2
	return dist.Length() <= float64(core.TileSize)*5 && (alignedHorizontally || alignedVertically)
}

func (c *GoblinEnemy) getAttackVector(player core.Player) core.Vector {
	dist := core.Vector(player.Location()).Minus(core.Vector(c.Location()))
	alignedHorizontally := math.Abs(dist.Y) < float64(core.TileSize)/2
	alignedVertically := math.Abs(dist.X) < float64(core.TileSize)/2
	if alignedHorizontally {
		if dist.X < 0 {
			return core.Vector{X: -1, Y: 0}
		} else {
			return core.Vector{X: 1, Y: 0}
		}
	} else if alignedVertically {
		if dist.Y < 0 {
			return core.Vector{X: 0, Y: -1}
		} else {
			return core.Vector{X: 0, Y: 1}
		}
	}
	return core.Vector{X: 0, Y: 0}
}

// updateMovement handles the movement logic and target finding.
func (c *GoblinEnemy) updateWalk(level core.Level, player core.Player) {
	if c.shouldAttackPlayer(player) {
		c.state = GoblinAttacking
		c.animations[GoblinAttacking][c.direction].Reset()
		c.velocity = c.getAttackVector(player).Scale(GoblinRunSpeed)
		return
	}

	c.turnTimeCounter++
	if c.turnTimeCounter >= GoblinTurnFrames {
		c.turnTimeCounter = 0
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	// Apply movement and handle collisions
	// Use the BaseCharacter's collision handler, which modifies c.X/c.Y
	c.velocity = c.velocity.Normalize().Scale(GoblinWalkSpeed)
	var v core.Vector
	v.X = c.HandleTileCollisions(level, core.AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, core.AxisY, c.velocity.Y)
	if v.Length() < 0.01 {
		// Turn immediately if we run into something; don't reset the turn counter, though
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	c.direction = core.VectorToDirection(v)
}

func (c *GoblinEnemy) isNearPlayer(playerLoc core.Location) bool {
	dist := core.Vector(playerLoc).Minus(core.Vector(c.Location()))
	return dist.Length() <= core.TileSize
}

func (c *GoblinEnemy) getActiveDamageSources(player core.Player) []*core.DamageSource {
	sources := []*core.DamageSource{}

	if c.state == GoblinAttacking {
		anim := c.animations[GoblinAttacking][c.direction]
		if anim == nil {
			return nil
		}

		// Check if we have an attack config for the current direction and animation frame index
		animIndex := anim.CurrentFrameIndex()
		if dirConfigs, ok := c.attackHitboxes[c.direction]; ok {
			if config, ok := dirConfigs[animIndex]; ok {
				worldHitbox := config.HitBox.Offset(c.Loc.X, c.Loc.Y)
				sources = append(sources, core.NewDamageSource(core.TagEnemy, worldHitbox, config.Damage))
			}
		}
	}

	// Add goblin pushbox when near player as well.
	if c.isNearPlayer(player.Location()) && (c.state == GoblinAttacking || c.state == GoblinWalking) {
		sources = append(sources, core.NewDamageSource(core.TagEnemy, c.GetHurtBox(), 1))
	}

	return sources
}

func (c *GoblinEnemy) updateAttack(level core.Level, player core.Player) {
	if (c.velocity.X < 0 && c.Loc.X <= player.Location().X-float64(core.TileSize)) ||
		(c.velocity.X > 0 && c.Loc.X >= player.Location().X+float64(core.TileSize)) ||
		(c.velocity.Y < 0 && c.Loc.Y <= player.Location().Y-float64(core.TileSize)) ||
		(c.velocity.Y > 0 && c.Loc.Y >= player.Location().Y+float64(core.TileSize)) {
		// Passed the player
		c.state = GoblinWalking
		return
	}

	var v core.Vector
	v.X = c.HandleTileCollisions(level, core.AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, core.AxisY, c.velocity.Y)
	c.direction = core.VectorToDirection(v)
}

func (c *GoblinEnemy) Update(level core.Level, player core.Player) core.UpdateResult {
	var actions []core.Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.SrcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow GoblinHurt animation to finish during knockback
		if c.state == GoblinHurt && c.animations[GoblinHurt][c.direction].IsFinished() {
			c.state = GoblinWalking
		}

		if c.state != GoblinDying {
			return core.UpdateResult{Actions: actions}
		}
	}

	if c.state != GoblinDying && c.UpdateStun() {
		c.state = GoblinStunned
		return core.UpdateResult{Actions: actions}
	} else if c.state == GoblinStunned {
		c.state = GoblinWalking
	}

	// Handle Dying state which overrides all AI
	if c.state == GoblinDying {
		if c.animations[GoblinDying][c.direction].IsFinished() {
			c.Dead = true
			actions = append(actions, core.Action{
				Type:       core.ActionGainXP,
				Experience: c.Experience,
			})
		}
		return core.UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == GoblinHurt {
		if !c.animations[GoblinHurt][c.direction].IsFinished() {
			return core.UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to walking
		c.state = GoblinWalking
		c.animations[GoblinWalking][c.direction].Reset()
	}

	// Core AI Logic
	switch c.state {
	case GoblinWalking:
		c.updateWalk(level, player)
	case GoblinAttacking:
		c.updateAttack(level, player)
	}

	for _, ds := range c.getActiveDamageSources(player) {
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return core.UpdateResult{Actions: actions}
}

func (c *GoblinEnemy) CanRemove() bool {
	return false
}
