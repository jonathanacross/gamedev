package main

import (
	"math"
	"math/rand"
)

type GhostEnemyState int

const (
	GhostIdle GhostEnemyState = iota
	GhostHurt
	GhostDying
	GhostAttacking
	GhostMoving

	// Movement constants
	GhostWalkSpeed  float64 = 0.6
	GhostRunSpeed   float64 = 1.2
	GhostTurnFrames int     = 60 * 3 // every 3 seconds
)

type GhostEnemy struct {
	BaseCharacter
	spriteSheet    *SpriteSheet
	animations     map[GhostEnemyState]map[Direction]*Animation
	attackHitboxes map[Direction]map[int]DamageSourceConfig

	// AI
	state           GhostEnemyState
	velocity        Vector
	direction       Direction
	turnTimeCounter int
}

func NewGhostEnemy(startLoc Location) *GhostEnemy {
	ssColumns := 12
	ssRows := 20
	directionOffsets := map[Direction]int{
		Down:  0 * ssColumns,
		Up:    1 * ssColumns,
		Left:  2 * ssColumns,
		Right: 3 * ssColumns,
	}
	framesPerState := ssColumns * len(directionOffsets)

	animations := map[GhostEnemyState]map[Direction]*Animation{
		GhostIdle:      NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0*framesPerState, directionOffsets, 10, true),
		GhostMoving:    NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5}, 1*framesPerState, directionOffsets, 10, true),
		GhostAttacking: NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 2*framesPerState, directionOffsets, 10, true),
		GhostHurt:      NewDirectionAnimationMap([]int{0, 1, 2, 3}, 3*framesPerState, directionOffsets, 10, false),
		GhostDying:     NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 4*framesPerState, directionOffsets, 10, false),
	}

	animations[GhostMoving][Up].SetRandomFrame()
	animations[GhostMoving][Down].SetRandomFrame()
	animations[GhostMoving][Left].SetRandomFrame()
	animations[GhostMoving][Right].SetRandomFrame()

	spriteSheet := NewSpriteSheet(64, 64, ssColumns, ssRows)
	hitbox := Rect{
		Left:   -7,
		Top:    -7,
		Right:  7,
		Bottom: 7,
	}
	// Define a simple attack hitbox that's only active on the 2nd and 3rd frames (index 1 and 2 in the short animation array)
	attackHitboxes := make(map[Direction]map[int]DamageSourceConfig)
	baseDmg := 1

	attackHitboxes[Up] = map[int]DamageSourceConfig{
		2: {HitBox: Rect{Left: 3, Top: -24, Right: 13, Bottom: 4}, Damage: baseDmg},
		3: {HitBox: Rect{Left: 3, Top: -27, Right: 13, Bottom: 4}, Damage: baseDmg},
		4: {HitBox: Rect{Left: 3, Top: -27, Right: 13, Bottom: 4}, Damage: baseDmg},
	}
	attackHitboxes[Down] = map[int]DamageSourceConfig{
		2: {HitBox: Rect{Left: -13, Top: -4, Right: -3, Bottom: 24}, Damage: baseDmg},
		3: {HitBox: Rect{Left: -13, Top: -4, Right: -3, Bottom: 27}, Damage: baseDmg},
		4: {HitBox: Rect{Left: -13, Top: -4, Right: -3, Bottom: 27}, Damage: baseDmg},
	}
	attackHitboxes[Left] = map[int]DamageSourceConfig{
		2: {HitBox: Rect{Left: -25, Top: 0, Right: 6, Bottom: 10}, Damage: baseDmg},
		3: {HitBox: Rect{Left: -28, Top: 0, Right: 6, Bottom: 10}, Damage: baseDmg},
		4: {HitBox: Rect{Left: -28, Top: 0, Right: 6, Bottom: 10}, Damage: baseDmg},
	}
	attackHitboxes[Right] = map[int]DamageSourceConfig{
		2: {HitBox: Rect{Left: -6, Top: 0, Right: 25, Bottom: 10}, Damage: baseDmg},
		3: {HitBox: Rect{Left: -6, Top: 0, Right: 28, Bottom: 10}, Damage: baseDmg},
		4: {HitBox: Rect{Left: -6, Top: 0, Right: 28, Bottom: 10}, Damage: baseDmg},
	}

	return &GhostEnemy{
		BaseCharacter: BaseCharacter{
			BasePhysical: BasePhysical{
				BaseSprite: BaseSprite{
					Location: startLoc,
					drawOffset: Location{
						X: 32,
						Y: 32,
					},
					srcRect: spriteSheet.Rect(0),
					image:   GhostSpritesImage,
				},
				pushBoxOffset: hitbox,
			},
			Health:          5,
			MaxHealth:       5,
			isDead:          false,
			KnockbackFrames: 0,
		},
		spriteSheet:     spriteSheet,
		animations:      animations,
		attackHitboxes:  attackHitboxes,
		state:           GhostMoving,
		direction:       Left,
		velocity:        getRandomVelocity(),
		turnTimeCounter: rand.Intn(GhostTurnFrames),
	}
}

func (c *GhostEnemy) ApplyKnockback(force Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != GhostDying {
		c.state = GhostHurt
		c.animations[GhostHurt][c.direction].Reset()
	}
}

func (c *GhostEnemy) IsKnockedBack() bool {
	return c.KnockbackFrames > 0
}

func (c *GhostEnemy) TakeDamage(damage int) {
	if c.isDead || c.state == GhostDying || c.state == GhostHurt {
		return
	}

	// TODO: consider a state transition like Player that handles animation reset
	c.state = GhostHurt
	c.animations[GhostHurt][c.direction].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = GhostDying
	}
}

func (c *GhostEnemy) shouldAttackPlayer(player *Player) bool {
	// Player should be near the goblin, and aligned either horizontally or vertically
	dist := Vector(player.Location()).Minus(Vector(c.Location()))
	alignedHorizontally := math.Abs(dist.Y) < TileSize/2
	alignedVertically := math.Abs(dist.X) < TileSize/2
	return dist.Length() <= TileSize*5 && (alignedHorizontally || alignedVertically)
}

func (c *GhostEnemy) getAttackVector(player *Player) Vector {
	dist := Vector(player.Location()).Minus(Vector(c.Location()))
	alignedHorizontally := math.Abs(dist.Y) < TileSize/2
	alignedVertically := math.Abs(dist.X) < TileSize/2
	if alignedHorizontally {
		if dist.X < 0 {
			return Vector{-1, 0}
		} else {
			return Vector{1, 0}
		}
	} else if alignedVertically {
		if dist.Y < 0 {
			return Vector{0, -1}
		} else {
			return Vector{0, 1}
		}
	}
	return Vector{0, 0}
}

// updateMovement handles the movement logic and target finding.
func (c *GhostEnemy) updateWalk(level *Level, player *Player) {
	if c.shouldAttackPlayer(player) {
		c.state = GhostAttacking
		c.animations[GhostAttacking][c.direction].Reset()
		c.velocity = c.getAttackVector(player).Scale(GhostRunSpeed)
		return
	}

	c.turnTimeCounter++
	if c.turnTimeCounter >= GhostTurnFrames {
		c.turnTimeCounter = 0
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	// Apply movement and handle collisions
	// Use the BaseCharacter's collision handler, which modifies c.X/c.Y
	c.velocity = c.velocity.Normalize().Scale(GhostWalkSpeed)
	var v Vector
	v.X = c.HandleTileCollisions(level, AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, AxisY, c.velocity.Y)
	if v.Length() < 0.01 {
		// Turn immediately if we run into something; don't reset the turn counter, though
		c.velocity = c.velocity.Rotate(math.Pi / 2)
	}

	c.direction = VectorToDirection(v)
}

func (c *GhostEnemy) isNearPlayer(playerLoc Location) bool {
	dist := Vector(playerLoc).Minus(Vector(c.Location()))
	return dist.Length() <= TileSize
}

func (c *GhostEnemy) getActiveDamageSources(player *Player) []*DamageSource {
	sources := []*DamageSource{}

	if c.state == GhostAttacking {
		anim := c.animations[GhostAttacking][c.direction]
		if anim == nil {
			return nil
		}

		// Check if we have an attack config for the current direction and animation frame index
		animIndex := anim.frameIndex
		if dirConfigs, ok := c.attackHitboxes[c.direction]; ok {
			if config, ok := dirConfigs[animIndex]; ok {
				worldHitbox := config.HitBox.Offset(c.X, c.Y)
				sources = append(sources, NewDamageSource(TagEnemy, worldHitbox, config.Damage))
			}
		}
	}

	// Add goblin pushbox when near player as well.
	if c.isNearPlayer(player.Location()) && (c.state == GhostAttacking || c.state == GhostMoving) {
		sources = append(sources, NewDamageSource(TagEnemy, c.GetHurtBox(), 1))
	}

	return sources
}

func (c *GhostEnemy) updateAttack(level *Level, player *Player) {
	if (c.velocity.X < 0 && c.X <= player.X-TileSize) ||
		(c.velocity.X > 0 && c.X >= player.X+TileSize) ||
		(c.velocity.Y < 0 && c.Y <= player.Y-TileSize) ||
		(c.velocity.Y > 0 && c.Y >= player.Y+TileSize) {
		// Passed the player
		c.state = GhostMoving
		return
	}

	var v Vector
	v.X = c.HandleTileCollisions(level, AxisX, c.velocity.X)
	v.Y = c.HandleTileCollisions(level, AxisY, c.velocity.Y)
	c.direction = VectorToDirection(v)
}

func (c *GhostEnemy) Update(level *Level, player *Player) UpdateResult {
	var actions []Action

	// Animation Update
	// The current animation is determined by the previous frame's state/direction.
	currentAnim := c.animations[c.state][c.direction]
	currentAnim.Update()
	c.srcRect = c.spriteSheet.Rect(currentAnim.Frame())

	// Priority State Checks (Dying, Knockback/Hurt)
	if c.UpdateKnockback(level) {
		// Allow GhostHurt animation to finish during knockback
		if c.state == GhostHurt && c.animations[GhostHurt][c.direction].IsFinished() {
			c.state = GhostMoving
		}

		if c.state != GhostDying {
			return UpdateResult{Actions: actions}
		}
	}

	// Handle Dying state which overrides all AI
	if c.state == GhostDying {
		if c.animations[GhostDying][c.direction].IsFinished() {
			c.isDead = true
		}
		return UpdateResult{Actions: actions}
	}

	// Handle Hurt state which overrides AI while animation plays
	if c.state == GhostHurt {
		if !c.animations[GhostHurt][c.direction].IsFinished() {
			return UpdateResult{Actions: actions} // Wait for hurt anim to finish
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
		actions = append(actions, Action{
			Type:         ActionCreateDamageSource,
			DamageSource: ds,
		})
	}

	return UpdateResult{Actions: actions}
}

func (c *GhostEnemy) CanRemove() bool {
	// TODO: consider if we want to merge this with isdead.
	return false
}
