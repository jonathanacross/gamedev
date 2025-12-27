package enemy

import (
	"math"
	"math/rand"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
	"time"
)

type GolemEnemyState int

const (
	GolemIdle GolemEnemyState = iota
	GolemHurt
	GolemDying
	GolemAttacking
	GolemWalking
	GolemStunned

	// Movement constants
	GolemWalkSpeed        float64       = 0.6
	GolemTurnFrames       int           = 60 * 3 // every 3 seconds
	GolemThrowBoulderTime time.Duration = 1 * time.Second
)

type GolemEnemy struct {
	character.BaseCharacter
	spriteSheet    *core.SpriteSheet
	animations     map[GolemEnemyState]map[core.Direction]*core.Animation
	attackHitboxes map[core.Direction]map[int]core.DamageSourceConfig

	// AI
	state             GolemEnemyState
	velocity          core.Vector
	direction         core.Direction
	turnTimeCounter   int
	boulderThrowTimer *core.Timer
}

// Already defined in goblin
// func getRandomVelocity() core.Vector {
// 	velocities := []core.Vector{
// 		{X: 1, Y: 0},
// 		{X: -1, Y: 0},
// 		{X: 0, Y: 1},
// 		{X: 0, Y: -1},
// 	}
// 	return velocities[rand.Intn(len(velocities))]
// }

func NewGolemEnemy(startLoc core.Location) *GolemEnemy {
	ssColumns := 9
	ssRows := 32
	directionOffsets := map[core.Direction]int{
		core.Down:  0 * ssColumns,
		core.Up:    1 * ssColumns,
		core.Left:  2 * ssColumns,
		core.Right: 3 * ssColumns,
	}
	framesPerState := ssColumns * len(directionOffsets)

	animations := map[GolemEnemyState]map[core.Direction]*core.Animation{
		GolemIdle:      core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*framesPerState, directionOffsets, 10, true),
		GolemWalking:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 1*framesPerState, directionOffsets, 10, true),
		GolemAttacking: core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7, 8}, 2*framesPerState, directionOffsets, 10, true),
		GolemHurt:      core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, false),
		GolemDying:     core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7, 8}, 4*framesPerState, directionOffsets, 10, false),
		GolemStunned:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, true),
	}

	animations[GolemWalking][core.Up].SetRandomFrame()
	animations[GolemWalking][core.Down].SetRandomFrame()
	animations[GolemWalking][core.Left].SetRandomFrame()
	animations[GolemWalking][core.Right].SetRandomFrame()

	spriteSheet := core.NewSpriteSheet(128, 128, ssColumns, ssRows)
	hitbox := core.Rect{
		Left:   -20,
		Top:    -27,
		Right:  20,
		Bottom: 13,
	}
	// Define a simple attack hitbox that's only active on the 2nd and 3rd frames (index 1 and 2 in the short animation array)
	attackHitboxes := make(map[core.Direction]map[int]core.DamageSourceConfig)
	baseDmg := 1

	attackHitboxes[core.Up] = map[int]core.DamageSourceConfig{
		6: {HitBox: core.Rect{Left: -39, Top: -11, Right: 39, Bottom: 27}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -39, Top: -11, Right: 39, Bottom: 27}, Damage: baseDmg},
	}
	attackHitboxes[core.Down] = map[int]core.DamageSourceConfig{
		6: {HitBox: core.Rect{Left: -39, Top: -11, Right: 39, Bottom: 27}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -39, Top: -11, Right: 39, Bottom: 27}, Damage: baseDmg},
	}
	attackHitboxes[core.Left] = map[int]core.DamageSourceConfig{
		6: {HitBox: core.Rect{Left: -27, Top: -9, Right: 8, Bottom: 29}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -27, Top: -9, Right: 8, Bottom: 29}, Damage: baseDmg},
	}
	attackHitboxes[core.Right] = map[int]core.DamageSourceConfig{
		6: {HitBox: core.Rect{Left: -8, Top: -9, Right: 27, Bottom: 29}, Damage: baseDmg},
		7: {HitBox: core.Rect{Left: -8, Top: -9, Right: 27, Bottom: 29}, Damage: baseDmg},
	}

	return &GolemEnemy{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: startLoc,
					DrawOffset: core.Location{
						X: 64,
						Y: 64,
					},
					SrcRect: spriteSheet.Rect(0),
					Image:   assets.GolemSpritesImage,
				},
				PushBoxOffset: hitbox,
			},
			Health:          5,
			MaxHealth:       5,
			Experience:      6,
			Dead:            false,
			KnockbackFrames: 0,
		},
		spriteSheet:       spriteSheet,
		animations:        animations,
		attackHitboxes:    attackHitboxes,
		state:             GolemWalking,
		direction:         core.Left,
		velocity:          getRandomVelocity(),
		turnTimeCounter:   rand.Intn(GolemTurnFrames),
		boulderThrowTimer: core.NewTimer(GolemThrowBoulderTime),
	}
}

func (c *GolemEnemy) HandleHit(ds *core.DamageSource) {
	if ds.Type == core.DamageTypeStun {
		c.ApplyStun(ds.Duration)
	} else {
		c.TakeDamage(ds.Damage)
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), c.Location(), core.KnockbackForce)
		c.ApplyKnockback(force, core.KnockbackDuration)
	}
}

func (c *GolemEnemy) ApplyKnockback(force core.Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != GolemDying {
		c.state = GolemHurt
		c.animations[GolemHurt][c.direction].Reset()
	}
}

func (c *GolemEnemy) ApplyStun(duration int) {
	c.BaseCharacter.ApplyStun(duration)
	c.state = GolemStunned
	c.animations[GolemStunned][c.direction].Reset()
}

func (c *GolemEnemy) TakeDamage(damage int) {
	if c.Dead || c.state == GolemDying || c.state == GolemHurt {
		return
	}

	c.state = GolemHurt
	c.animations[GolemHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = GolemDying
	}
}

func (c *GolemEnemy) throwBoulders(player core.Player) *core.Action {
	c.boulderThrowTimer.Update()

	if !c.boulderThrowTimer.IsReady() {
		return nil
	}

	distToPlayer := core.Vector(player.Location()).Minus(core.Vector(c.Location()))
	if distToPlayer.Length() < float64(core.TileSize)*5 {
		c.boulderThrowTimer.Reset()
		return &core.Action{
			Type:      core.ActionThrowBoulder,
			Location:  c.Location(),
			Direction: distToPlayer.Normalize(),
		}
	}

	return nil
}

func (c *GolemEnemy) shouldAttackPlayer(player core.Player) bool {
	// Player should be near the goblin, and aligned either horizontally or vertically
	dist := core.Vector(player.Location()).Minus(core.Vector(c.Location()))
	alignedHorizontally := math.Abs(dist.Y) < float64(core.TileSize)/2
	alignedVertically := math.Abs(dist.X) < float64(core.TileSize)/2
	return dist.Length() <= float64(core.TileSize)*5 && (alignedHorizontally || alignedVertically)
}

func (c *GolemEnemy) getAttackVector(player core.Player) core.Vector {
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
func (c *GolemEnemy) updateWalk(level core.Level, player core.Player) {
	if c.shouldAttackPlayer(player) {
		c.state = GolemAttacking
		c.animations[GolemAttacking][c.direction].Reset()
		c.velocity = c.getAttackVector(player).Scale(GolemWalkSpeed)
		return
	}

	c.turnTimeCounter++
	if c.turnTimeCounter >= GolemTurnFrames {
		c.turnTimeCounter = 0
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	// Apply movement and handle collisions
	// Use the BaseCharacter's collision handler, which modifies c.X/c.Y
	c.velocity = c.velocity.Normalize().Scale(GolemWalkSpeed)
	var v core.Vector
	v.X = c.HandleTileCollisions(level, core.AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, core.AxisY, c.velocity.Y)
	if v.Length() < 0.01 {
		// Turn immediately if we run into something; don't reset the turn counter, though
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	c.direction = core.VectorToDirection(v)
}

func (c *GolemEnemy) isNearPlayer(playerLoc core.Location) bool {
	dist := core.Vector(playerLoc).Minus(core.Vector(c.Location()))
	return dist.Length() <= core.TileSize
}

func (c *GolemEnemy) getActiveDamageSources(player core.Player) []*core.DamageSource {
	sources := []*core.DamageSource{}

	if c.state == GolemAttacking {
		anim := c.animations[GolemAttacking][c.direction]
		if anim == nil {
			return nil
		}

		// Check if we have an attack config for the current direction and animation frame index
		animIndex := anim.CurrentFrameIndex()
		if dirConfigs, ok := c.attackHitboxes[c.direction]; ok {
			if config, ok := dirConfigs[animIndex]; ok {
				worldHitbox := config.HitBox.Offset(c.Loc.X, c.Loc.Y)
				sources = append(sources, core.NewDamageSource(core.TagEnemy, worldHitbox, core.DamageTypePhysical, config.Damage))
			}
		}
	}

	// Add goblin pushbox when near player as well.
	if c.isNearPlayer(player.Location()) && (c.state == GolemAttacking || c.state == GolemWalking) {
		sources = append(sources, core.NewDamageSource(core.TagEnemy, c.GetHurtBox(), core.DamageTypePhysical, 1))
	}

	return sources
}

func (c *GolemEnemy) updateAttack(level core.Level, player core.Player) {
	if (c.velocity.X < 0 && c.Loc.X <= player.Location().X-float64(core.TileSize)) ||
		(c.velocity.X > 0 && c.Loc.X >= player.Location().X+float64(core.TileSize)) ||
		(c.velocity.Y < 0 && c.Loc.Y <= player.Location().Y-float64(core.TileSize)) ||
		(c.velocity.Y > 0 && c.Loc.Y >= player.Location().Y+float64(core.TileSize)) {
		// Passed the player
		c.state = GolemWalking
		return
	}

	var v core.Vector
	v.X = c.HandleTileCollisions(level, core.AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, core.AxisY, c.velocity.Y)
	c.direction = core.VectorToDirection(v)
}

func (c *GolemEnemy) Update(level core.Level, player core.Player) core.UpdateResult {
	var actions []core.Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.SrcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow GolemHurt animation to finish during knockback
		if c.state == GolemHurt && c.animations[GolemHurt][c.direction].IsFinished() {
			c.state = GolemWalking
		}

		if c.state != GolemDying {
			return core.UpdateResult{Actions: actions}
		}
	}

	if c.state != GolemDying && c.UpdateStun() {
		c.state = GolemStunned
		return core.UpdateResult{Actions: actions}
	} else if c.state == GolemStunned {
		c.state = GolemWalking
	}

	// Handle Dying state which overrides all AI
	if c.state == GolemDying {
		if c.animations[GolemDying][c.direction].IsFinished() {
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
	if c.state == GolemHurt {
		if !c.animations[GolemHurt][c.direction].IsFinished() {
			return core.UpdateResult{Actions: actions} // Wait for hurt anim to finish
		}

		// Hurt animation finished, return to walking
		c.state = GolemWalking
		c.animations[GolemWalking][c.direction].Reset()
	}

	// Core AI Logic
	switch c.state {
	case GolemWalking:
		c.updateWalk(level, player)
	case GolemAttacking:
		c.updateAttack(level, player)
	}

	action := c.throwBoulders(player)
	if action != nil {
		actions = append(actions, *action)
	}

	for _, ds := range c.getActiveDamageSources(player) {
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return core.UpdateResult{Actions: actions}
}

func (c *GolemEnemy) CanRemove() bool {
	return false
}
