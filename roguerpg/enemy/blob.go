package enemy

import (
	"math/rand"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
)

type BlobEnemyState int

const (
	BlobMoving BlobEnemyState = iota
	BlobAttacking
	BlobHurt
	BlobStunned
	BlobDying

	// Movement constants
	BlobMoveSpeed float64 = 0.5
	MaxWaitFrames int     = 60 // Max 1 second wait (up to 60 frames)
)

type BlobEnemy struct {
	character.BaseCharacter
	spriteSheet *core.SpriteSheet
	animations  map[BlobEnemyState]*core.Animation

	// AI
	state BlobEnemyState

	idling             bool
	moveStartLocation  core.Location
	moveTargetLocation core.Location
	currentFrame       int // Frame counter for the current move or wait action
	waitFrames         int // Total frames to wait when idle
}

func NewBlobEnemy(startLoc core.Location) *BlobEnemy {

	animations := map[BlobEnemyState]*core.Animation{
		BlobMoving:    core.NewAnimation([]int{0, 1, 2}, 20, true),
		BlobAttacking: core.NewAnimation([]int{15, 16, 17, 18}, 10, false),
		BlobHurt:      core.NewAnimation([]int{5, 6, 5, 6}, 10, false),
		BlobDying:     core.NewAnimation([]int{5, 6, 10, 11, 12, 13, 14}, 10, false),
		BlobStunned:   core.NewAnimation([]int{5, 6}, 10, true),
	}
	animations[BlobMoving].SetRandomFrame()

	spriteSheet := core.NewSpriteSheet(32, 32, 5, 3)
	hitbox := core.Rect{Left: -7, Top: -8, Right: 8, Bottom: 5}

	return &BlobEnemy{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: startLoc,
					DrawOffset: core.Location{
						X: 16,
						Y: 16,
					},
					SrcRect: spriteSheet.Rect(0),
					Image:   assets.BlobSpritesImage,
				},
				PushBoxOffset: hitbox,
			},
			Health:          3,
			MaxHealth:       3,
			Experience:      2,
			Dead:            false,
			KnockbackFrames: 0,
		},
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       BlobMoving,
		idling:      true,
		waitFrames:  rand.Intn(MaxWaitFrames) + 1,
	}
}

func (c *BlobEnemy) ApplyKnockback(force core.Vector, duration int) {
	c.BaseCharacter.ApplyKnockback(force, duration)

	if c.IsKnockedBack() && c.state != BlobDying {
		c.state = BlobHurt
		c.animations[BlobHurt].Reset()
	}
}

func (c *BlobEnemy) HandleHit(ds *core.DamageSource) {
	if ds.Type == core.DamageTypeStun {
		c.ApplyStun(ds.Duration)
	} else {
		c.TakeDamage(ds.Damage)
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), c.Location(), core.KnockbackForce)
		c.ApplyKnockback(force, core.KnockbackDuration)
	}
}

func (c *BlobEnemy) ApplyStun(duration int) {
	c.BaseCharacter.ApplyStun(duration)
	c.state = BlobStunned
	c.animations[BlobStunned].Reset()
}

func (c *BlobEnemy) TakeDamage(damage int) {
	if c.Dead || c.state == BlobDying || c.state == BlobHurt {
		return
	}

	c.state = BlobHurt
	c.animations[BlobHurt].Reset()

	c.Health -= damage
	if c.Health <= 0 {
		c.state = BlobDying
	}
}

// findNewTargetTile attempts to find a random, adjacent, non-solid tile.
func (c *BlobEnemy) findNewTargetTile(level core.Level) bool {
	// Get current tile coordinates
	tx, ty := level.WorldToTile(c.Location())

	// Define the 4 cardinal directions for "adjacent square"
	directions := []struct{ dx, dy int }{
		{0, 1},  // Down
		{0, -1}, // Up
		{1, 0},  // Right
		{-1, 0}, // Left
	}

	// Shuffle directions to pick a random one first
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newTx := tx + dir.dx
		newTy := ty + dir.dy

		if !level.IsTileSolid(newTx, newTy) {
			// Found an open tile. Set up the movement.
			c.moveStartLocation = c.Location()
			c.moveTargetLocation = level.TileToWorld(newTx, newTy)
			c.currentFrame = 0 // Reset frame counter for movement
			return true
		}
	}

	// No adjacent open tile found
	return false
}

func (c *BlobEnemy) updateMovement(level core.Level) {
	if c.idling {
		c.waitFrames--
		if c.waitFrames <= 0 {
			// Wait time is over. Look for a new target tile.
			if c.findNewTargetTile(level) {
				c.idling = false
			} else {
				// Enemy is cornered or blocked. Wait again.
				c.waitFrames = rand.Intn(MaxWaitFrames) + 1
			}
		}
	} else {
		target := core.Vector(c.moveTargetLocation).Minus(core.Vector(c.Location()))

		distance := target.Length()
		if distance <= BlobMoveSpeed {
			// We are close enough to snap to the target.
			c.SetLocation(c.moveTargetLocation)

			// Wait for a short time.
			c.idling = true
			c.waitFrames = rand.Intn(MaxWaitFrames) + 1
		}

		velocity := target.Normalize().Scale(BlobMoveSpeed)
		c.HandleTileCollisions(level, core.AxisX, velocity.X)
		c.HandleTileCollisions(level, core.AxisY, velocity.Y)
	}
}

func (c *BlobEnemy) Update(level core.Level, player core.Player) core.UpdateResult {
	c.animations[c.state].Update()
	c.SrcRect = c.spriteSheet.Rect(c.animations[c.state].Frame())
	var actions []core.Action

	if c.UpdateKnockback(level) {
		// Ensure the BlobHurt animation can finish, even during knockback
		if c.state == BlobHurt && c.animations[BlobHurt].IsFinished() {
			c.state = BlobMoving
			c.animations[BlobMoving].Reset()
		}

		if c.state != BlobDying {
			return core.UpdateResult{Actions: actions} // Skip AI and normal movement logic
		}
	}

	if c.state != BlobDying && c.UpdateStun() {
		c.state = BlobStunned
		return core.UpdateResult{Actions: actions}
	} else if c.state == BlobStunned {
		c.state = BlobMoving
	}

	c.updateMovement(level)

	switch c.state {
	case BlobMoving:
		if c.IsNearTo(player, 3) {
			c.state = BlobAttacking
			c.animations[BlobAttacking].Reset()
		}

	case BlobAttacking:
		hitBox := c.GetHurtBox()
		if c.animations[BlobAttacking].Frame() != 0 {
			ds := core.NewDamageSource(core.TagEnemy, hitBox, core.DamageTypePhysical, 1)
			actions = append(actions, core.Action{
				Type:         core.ActionCreateDamageSource,
				DamageSource: ds,
			})
		}

		if c.animations[BlobAttacking].IsFinished() {
			c.state = BlobMoving
		}

	case BlobHurt:
		if c.animations[BlobHurt].IsFinished() {
			c.state = BlobMoving
		}

	case BlobDying:
		if c.animations[BlobDying].IsFinished() {
			c.Dead = true
			actions = append(actions, core.Action{
				Type:       core.ActionGainXP,
				Experience: c.Experience,
			})
			actions = append(actions, maybeDropHeart(c.Location())...)
		}
	}

	return core.UpdateResult{Actions: actions}
}

func (c *BlobEnemy) CanRemove() bool {
	return false
}
