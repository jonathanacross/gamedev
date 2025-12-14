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
	GhostAttacking
	GhostMoving
	GhostStunned

	// Movement constants
	GhostMaxFloatSpeed float64 = 0.8
	GhostAttackSpeed   float64 = 0.8
	GhostXFrequency    float64 = 0.02 // radians per frame
	GhostYFrequency    float64 = 0.03 // radians per frame
)

type GhostEnemy struct {
	character.BaseCharacter
	spriteSheet    *core.SpriteSheet
	animations     map[GhostEnemyState]map[core.Direction]*core.Animation
	attackHitboxes map[core.Direction]map[int]core.DamageSourceConfig

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
		GhostIdle:      core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*framesPerState, directionOffsets, 10, true),
		GhostMoving:    core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 1*framesPerState, directionOffsets, 10, true),
		GhostAttacking: core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 2*framesPerState, directionOffsets, 8, true),
		GhostHurt:      core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, false),
		GhostDying:     core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 4*framesPerState, directionOffsets, 8, false),
		GhostStunned:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, true),
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
	attackHitboxes := make(map[core.Direction]map[int]core.DamageSourceConfig)
	baseDmg := 1

	attackHitboxes[core.Up] = map[int]core.DamageSourceConfig{
		3: {HitBox: core.Rect{Left: -9, Top: -26, Right: 9, Bottom: -3}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -10, Top: -30, Right: 10, Bottom: -6}, Damage: baseDmg},
		5: {HitBox: core.Rect{Left: -11, Top: -11, Right: 10, Bottom: 13}, Damage: baseDmg},
		6: {HitBox: core.Rect{Left: -11, Top: -9, Right: 10, Bottom: 13}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -11, Top: -11, Right: 10, Bottom: 11}, Damage: baseDmg},
	}
	attackHitboxes[core.Down] = map[int]core.DamageSourceConfig{
		3: {HitBox: core.Rect{Left: -9, Top: -26, Right: 9, Bottom: -3}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -10, Top: -30, Right: 10, Bottom: -6}, Damage: baseDmg},
		5: {HitBox: core.Rect{Left: -11, Top: 5, Right: 10, Bottom: 20}, Damage: baseDmg},
		6: {HitBox: core.Rect{Left: -11, Top: -2, Right: 10, Bottom: 20}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -11, Top: -4, Right: 10, Bottom: 18}, Damage: baseDmg},
	}
	attackHitboxes[core.Left] = map[int]core.DamageSourceConfig{
		3: {HitBox: core.Rect{Left: -9, Top: -26, Right: 9, Bottom: -3}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -4, Top: -30, Right: 16, Bottom: -5}, Damage: baseDmg},
		5: {HitBox: core.Rect{Left: -24, Top: -6, Right: -3, Bottom: 16}, Damage: baseDmg},
		6: {HitBox: core.Rect{Left: -24, Top: -6, Right: -3, Bottom: 16}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -24, Top: -8, Right: -3, Bottom: 14}, Damage: baseDmg},
	}
	attackHitboxes[core.Right] = map[int]core.DamageSourceConfig{
		3: {HitBox: core.Rect{Left: -9, Top: -26, Right: 9, Bottom: -3}, Damage: baseDmg},
		4: {HitBox: core.Rect{Left: -16, Top: -30, Right: 4, Bottom: -5}, Damage: baseDmg},
		5: {HitBox: core.Rect{Left: 3, Top: -6, Right: 24, Bottom: 16}, Damage: baseDmg},
		6: {HitBox: core.Rect{Left: 3, Top: -6, Right: 24, Bottom: 16}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: 3, Top: -8, Right: 24, Bottom: 14}, Damage: baseDmg},
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
		spriteSheet:    spriteSheet,
		animations:     animations,
		attackHitboxes: attackHitboxes,
		state:          GhostMoving,
		direction:      core.Left,
		velocity:       core.Vector{0, 0}, // Init from random below?
		// Logic: velocity = getRandomVelocity() in original.
		// Wait, GhostEnemy used getRandomVelocity in NewGhostEnemy?
		// Original: velocity: getRandomVelocity().
		// But getRandomVelocity was defined in GoblinEnemy file or ghost?
		// Ah, getRandomVelocity was in Goblin. Ghost uses it?
		// If it was in same package main, yes.
		// Now they are in package enemy. `getRandomVelocity` is unexported in goblin.go.
		// So I must duplicate it or export it or put in utils.
		// I will duplicate it for now or assume it's there.
		movementTimer: rand.Float64() * 100,
	}
}

// Duplicate helper locally
func getRandomVelocityGhost() core.Vector {
	velocities := []core.Vector{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}
	return velocities[rand.Intn(len(velocities))]
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

func (c *GhostEnemy) shouldAttackPlayer(player core.Player) bool {
	// Player should be near the ghost, and aligned either horizontally or vertically
	dist := core.Vector(player.Location()).Minus(core.Vector(c.Location()))
	alignedHorizontally := math.Abs(dist.Y) < float64(core.TileSize)/2
	alignedVertically := math.Abs(dist.X) < float64(core.TileSize)/2
	return dist.Length() <= float64(core.TileSize)*5 && (alignedHorizontally || alignedVertically)
}

func (c *GhostEnemy) getAttackVector(player core.Player) core.Vector {
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
func (c *GhostEnemy) updateWalk(level core.Level, player core.Player) {
	if c.shouldAttackPlayer(player) {
		c.state = GhostAttacking
		c.animations[GhostAttacking][c.direction].Reset()
		c.velocity = c.getAttackVector(player).Scale(GhostAttackSpeed)
		return
	}

	c.movementTimer++

	// Lissajous pattern:
	// velocity.x = maxGhostVelocity * cos(xFrequency * t)
	// velocity.y = maxGhostVelocity * sin(yFrequency * t)
	vx := GhostMaxFloatSpeed * math.Cos(GhostXFrequency*c.movementTimer)
	vy := GhostMaxFloatSpeed * math.Sin(GhostYFrequency*c.movementTimer)
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

func (c *GhostEnemy) isNearPlayer(playerLoc core.Location) bool {
	dist := core.Vector(playerLoc).Minus(core.Vector(c.Location()))
	return dist.Length() <= core.TileSize
}

func (c *GhostEnemy) getActiveDamageSources(player core.Player) []*core.DamageSource {
	sources := []*core.DamageSource{}

	if c.state == GhostAttacking {
		anim := c.animations[GhostAttacking][c.direction]
		if anim == nil {
			return nil
		}

		// Check if we have an attack config for the current direction and animation frame index
		// Note: the ghost does not hurt the player when it is near them except during attack
		animIndex := anim.CurrentFrameIndex()

		if dirConfigs, ok := c.attackHitboxes[c.direction]; ok {
			if config, ok := dirConfigs[animIndex]; ok {
				worldHitbox := config.HitBox.Offset(c.Loc.X, c.Loc.Y)
				sources = append(sources, core.NewDamageSource(core.TagEnemy, worldHitbox, config.Damage))
			}
		}
	}

	return sources
}

func (c *GhostEnemy) updateAttack(level core.Level, player core.Player) {
	if (c.velocity.X < 0 && c.Loc.X <= player.Location().X-core.TileSize) ||
		(c.velocity.X > 0 && c.Loc.X >= player.Location().X+core.TileSize) ||
		(c.velocity.Y < 0 && c.Loc.Y <= player.Location().Y-core.TileSize) ||
		(c.velocity.Y > 0 && c.Loc.Y >= player.Location().Y+core.TileSize) {
		// Passed the player
		c.state = GhostMoving
		return
	}

	var v core.Vector
	v.X = c.HandleTileCollisions(level, core.AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, core.AxisY, c.velocity.Y)
	c.direction = core.VectorToDirection(v)
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
	switch c.state {
	case GhostMoving:
		c.updateWalk(level, player)
	case GhostAttacking:
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

func (c *GhostEnemy) CanRemove() bool {
	return false
}
