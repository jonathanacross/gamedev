package player

import (
	"time"

	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walking
	AttackingSword
	AttackingShield
	AttackingBow
	AttackingWand
	Hurt
	Dying
	Dead
)

const (
	PlayerSpeed  = 2.0
	BombCooldown = 750 * time.Millisecond
	BowCooldown  = 300 * time.Millisecond
	WandCooldown = 2000 * time.Millisecond
)

type PlayerAnimationKey struct {
	State     PlayerState
	Direction core.Direction
}

type Player struct {
	character.BaseCharacter
	images         map[PlayerState]*ebiten.Image
	spriteSheet    *core.SpriteSheet
	animations     map[PlayerState]map[core.Direction]*core.Animation
	state          PlayerState
	direction      core.Direction
	Velocity       core.Vector
	attackHitboxes map[core.Direction]map[int]core.DamageSourceConfig // Defines hitboxes for specific animation frames

	bombCooldownTimer *core.Timer
	bowCooldownTimer  *core.Timer
	wandCooldownTimer *core.Timer
	hasBoomerang      bool

	primaryWeapon   core.WeaponType
	secondaryWeapon core.WeaponType
	shieldHeld      bool
}

func NewPlayer() *Player {
	// TODO: these are attack boxes for the sword; need to add attack boxes for the shield
	attackHitboxes := make(map[core.Direction]map[int]core.DamageSourceConfig)

	baseDmg := 1

	// Set up Hitboxes for specific frames of the attack animation.
	// Downward swing
	attackHitboxes[core.Down] = map[int]core.DamageSourceConfig{
		1: {HitBox: core.Rect{Left: -16, Top: 0, Right: 16, Bottom: 22}, Damage: baseDmg},
		2: {HitBox: core.Rect{Left: -16, Top: 0, Right: 16, Bottom: 22}, Damage: baseDmg},
	}
	// Left swing
	attackHitboxes[core.Left] = map[int]core.DamageSourceConfig{
		1: {HitBox: core.Rect{Left: -22, Top: -24, Right: 0, Bottom: 8}, Damage: baseDmg},
		2: {HitBox: core.Rect{Left: -22, Top: -24, Right: 0, Bottom: 8}, Damage: baseDmg},
	}
	// Right swing
	attackHitboxes[core.Right] = map[int]core.DamageSourceConfig{
		1: {HitBox: core.Rect{Left: 0, Top: -24, Right: 22, Bottom: 8}, Damage: baseDmg},
		2: {HitBox: core.Rect{Left: 0, Top: -24, Right: 22, Bottom: 8}, Damage: baseDmg},
	}
	// Upward swing
	attackHitboxes[core.Up] = map[int]core.DamageSourceConfig{
		1: {HitBox: core.Rect{Left: -16, Top: -22, Right: 16, Bottom: 0}, Damage: baseDmg},
		2: {HitBox: core.Rect{Left: -16, Top: -22, Right: 16, Bottom: 0}, Damage: baseDmg},
	}

	ssColumns := 8
	ssRows := 6
	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := map[PlayerState]map[core.Direction]*core.Animation{
		Idle:            core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 10, true),
		Walking:         core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 10, true),
		AttackingSword:  core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false),
		AttackingShield: core.NewDirectionAnimationMap([]int{1}, 0, directionOffsets, 6, true),
		AttackingBow:    core.NewDirectionAnimationMap([]int{0, 1}, 0, directionOffsets, 6, false),
		AttackingWand:   core.NewDirectionAnimationMap([]int{0, 1, 2}, 0, directionOffsets, 4, false),
		Hurt:            core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false),
		Dying:           core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 8, false),
		Dead:            core.NewDirectionAnimationMap([]int{7}, 0, directionOffsets, 8, false),
	}

	charImages := map[PlayerState]*ebiten.Image{
		Idle:            assets.PlayerIdleSpritesImage,
		Walking:         assets.PlayerWalkSpritesImage,
		AttackingSword:  assets.PlayerAttackSwordSpritesImage,
		AttackingShield: assets.PlayerAttackShieldSpritesImage,
		AttackingBow:    assets.PlayerAttackBowSpritesImage,
		AttackingWand:   assets.PlayerAttackWandSpritesImage,
		Hurt:            assets.PlayerHurtSpritesImage,
		Dying:           assets.PlayerDeathSpritesImage,
		Dead:            assets.PlayerDeathSpritesImage,
	}

	spriteSheet := core.NewSpriteSheet(48, 64, ssColumns, ssRows)
	pushBox := core.Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	return &Player{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: core.Location{
						X: 0,
						Y: 0,
					},
					DrawOffset: core.Location{
						X: 25,
						Y: 38,
					},
					SrcRect: spriteSheet.Rect(0),
				},
				PushBoxOffset: pushBox,
			},
			Health:          8,
			MaxHealth:       8,
			Experience:      0,
			KnockbackFrames: 0,
		},
		images:            charImages,
		spriteSheet:       spriteSheet,
		animations:        animations,
		state:             Idle,
		direction:         core.Down,
		attackHitboxes:    attackHitboxes,
		bombCooldownTimer: core.NewTimer(BombCooldown),
		bowCooldownTimer:  core.NewTimer(BowCooldown),
		wandCooldownTimer: core.NewTimer(WandCooldown),
		hasBoomerang:      true,
		primaryWeapon:     core.WeaponSword,
		secondaryWeapon:   core.WeaponBoomerang,
	}
}

func (p *Player) GetCurrentAnimation() *core.Animation {
	animationSet, exists := p.animations[p.state]
	if !exists {
		return nil
	}
	animation, exists := animationSet[p.direction]
	if !exists {
		return nil
	}
	return animation
}

func (p *Player) Move(moveVector core.Vector) {
	// Don't change direction or velocity if attacking
	if p.state != AttackingSword && p.state != AttackingShield && p.state != AttackingBow {
		// Determine facing direction from the move vector.
		// Prioritize vertical direction for diagonal movement.
		if moveVector.Y < 0 {
			p.direction = core.Up
		} else if moveVector.Y > 0 {
			p.direction = core.Down
		} else if moveVector.X < 0 {
			p.direction = core.Left
		} else if moveVector.X > 0 {
			p.direction = core.Right
		}

		// Set velocity based on the move vector, normalized and scaled.
		p.Velocity = moveVector.Normalize().Scale(PlayerSpeed)
		p.TransitionState(Walking)
	}
}

func (p *Player) StopMoving() {
	// If currently Walking, transition to Idle
	if p.state == Walking {
		p.TransitionState(Idle)
	}
}

func (p *Player) AttackSword() {
	// Only start attack if Idle or Walking
	if p.state == Idle || p.state == Walking {
		p.TransitionState(AttackingSword)
	}
}

func (p *Player) AttackShield() {
	p.shieldHeld = true
	// Only start attack if Idle or Walking
	if p.state == Idle || p.state == Walking {
		p.TransitionState(AttackingShield)
	}
}

func (p *Player) UseBomb() *core.Action {
	// Optional: Block item use while attacking
	// if p.state == AttackingSword || p.state == AttackingShield {
	// 	return nil
	// }

	if p.bombCooldownTimer.IsReady() {
		p.bombCooldownTimer.Reset()

		// Calculate the bomb's spawn location based on player's direction
		loc := core.Vector(p.Location()).Plus(core.DirectionToVector(p.direction).Scale(float64(core.TileSize)))

		return &core.Action{
			Type:      core.ActionDropBomb,
			Location:  core.Location(loc),
			Direction: core.DirectionToVector(p.direction),
		}
	}
	return nil
}

func (p *Player) UseBow() {
	// Only start attack if Idle or Walking
	if (p.state == Idle || p.state == Walking) && p.bowCooldownTimer.IsReady() {
		p.TransitionState(AttackingBow)
	}
}

func (p *Player) AimWand() {
	// Only start attack if Idle or Walking
	if (p.state == Idle || p.state == Walking) && p.wandCooldownTimer.IsReady() {
		p.TransitionState(AttackingWand)
	}
}

func (p *Player) ShootStar(level core.Level) *core.Action {
	if p.wandCooldownTimer.IsReady() {
		p.wandCooldownTimer.Reset()

		// Calculate the bomb's spawn location based on player's direction
		starStartOffset := core.Vector{X: 0, Y: -4}
		loc := core.Vector(p.Location()).Plus(starStartOffset).Plus(core.DirectionToVector(p.direction).Scale(float64(core.TileSize)))

		var target core.Character
		if level != nil {
			target = level.FindNearestEnemy(p.Location(), 6*float64(core.TileSize))
		}

		return &core.Action{
			Type:      core.ActionCreateStar,
			Location:  core.Location(loc),
			Direction: core.DirectionToVector(p.direction),
			Target:    target,
		}
	}
	return nil
}

func (p *Player) ShootArrow() *core.Action {
	if p.bowCooldownTimer.IsReady() {
		p.bowCooldownTimer.Reset()

		// Calculate the bow spawn location based on player's direction.
		// since the player holds the bow a few pixels above the origin, we need an offet.
		arrowStartOffset := core.Vector{X: 0, Y: -6}
		loc := core.Vector(p.Location()).Plus(arrowStartOffset).Plus(core.DirectionToVector(p.direction).Scale(float64(core.TileSize)))

		return &core.Action{
			Type:      core.ActionShootArrow,
			Location:  core.Location(loc),
			Direction: core.DirectionToVector(p.direction),
		}
	}
	return nil
}

func (p *Player) UseBoomerang() *core.Action {
	// Optional: Block item use while attacking
	// if p.state == AttackingSword || p.state == AttackingShield {
	// 	return nil
	// }

	if p.hasBoomerang {
		p.hasBoomerang = false

		// Calculate boomerang spawn location
		loc := core.Vector(p.Location()).Plus(core.DirectionToVector(p.direction).Scale(float64(core.TileSize)))

		return &core.Action{
			Type:      core.ActionThrowBoomerang,
			Location:  core.Location(loc),
			Direction: core.DirectionToVector(p.direction),
		}
	}
	return nil
}

func (p *Player) PrimaryAttack() *core.Action {
	switch p.primaryWeapon {
	case core.WeaponSword:
		p.AttackSword()
		return nil
	case core.WeaponBomb:
		return p.UseBomb()
	case core.WeaponBoomerang:
		return p.UseBoomerang()
	case core.WeaponShield:
		p.AttackShield()
		return nil
	case core.WeaponBow:
		p.UseBow()
		return nil
	case core.WeaponWand:
		p.AimWand()
		return nil
	default:
		return nil
	}
}

func (p *Player) SecondaryAttack() *core.Action {
	switch p.secondaryWeapon {
	case core.WeaponSword:
		p.AttackSword()
		return nil
	case core.WeaponBomb:
		return p.UseBomb()
	case core.WeaponBoomerang:
		return p.UseBoomerang()
	case core.WeaponShield:
		p.AttackShield()
		return nil
	case core.WeaponBow:
		p.UseBow()
		return nil
	case core.WeaponWand:
		p.AimWand()
		return nil
	default:
		return nil
	}
}

func (p *Player) SwitchWeapon(weapon core.WeaponType) {
	p.secondaryWeapon = weapon
}

func (p *Player) GetPrimaryWeapon() core.WeaponType {
	return p.primaryWeapon
}

func (p *Player) GetSecondaryWeapon() core.WeaponType {
	return p.secondaryWeapon
}

func (p *Player) GetWeaponProgress(weapon core.WeaponType) float64 {
	switch weapon {
	case core.WeaponSword:
		return 1.0
	case core.WeaponBomb:
		return p.bombCooldownTimer.GetProgress()
	case core.WeaponBoomerang:
		if p.hasBoomerang {
			return 1.0
		} else {
			return 0.0
		}
	case core.WeaponShield:
		return 1.0
	case core.WeaponBow:
		return p.bowCooldownTimer.GetProgress()
	case core.WeaponWand:
		return p.wandCooldownTimer.GetProgress()
	default:
		return 0.0
	}
}

func (p *Player) IsActive() bool {
	// Input should be blocked if the player is Hurt, Dying, or Dead.
	// We allow input only when Idle, Walking, or Attacking (to allow chaining/canceling).
	return p.state != Hurt && p.state != Dying && p.state != Dead
}

// handleState runs the logic for the Player's current state and determines
// the next state, direction, and velocity (Vx/Vy).
func (p *Player) handleState(level core.Level) []core.Action {
	animation := p.GetCurrentAnimation()
	if animation == nil {
		p.state = Idle // Should never happen
		return nil
	}

	if p.Health <= 0 && p.state != Dying && p.state != Dead {
		p.TransitionState(Dying)
		return nil
	}

	// Check for terminal state completion.
	if p.state == Dying && animation.IsFinished() {
		p.TransitionState(Dead)
		return nil
	}

	// Once Dead, stop all logic.
	if p.state == Dead {
		p.Velocity = core.Vector{}
		return nil
	}

	// If knockback just finished, transition out of Hurt.
	if p.state == Hurt && !p.IsKnockedBack() && animation.IsFinished() {
		p.TransitionState(Idle)
	}

	// If attack just finished, transition out of Attacking.
	if p.state == AttackingSword && animation.IsFinished() {
		p.TransitionState(Idle)
	}
	if p.state == AttackingShield {
		if !p.shieldHeld {
			p.TransitionState(Idle)
		}
	}
	if p.state == AttackingBow && animation.IsFinished() {
		// Shoot arrow
		if action := p.ShootArrow(); action != nil {
			return []core.Action{*action}
		}
		p.TransitionState(Idle)
	}
	if p.state == AttackingWand && animation.IsFinished() {
		// Shoot wand
		if action := p.ShootStar(level); action != nil {
			return []core.Action{*action}
		}
		p.TransitionState(Idle)
	}

	// The Player's internal state handles movement physics based on its state:
	switch p.state {
	case Idle:
		// Idle means no velocity from movement input
		p.Velocity = core.Vector{}
	case Walking:
		// The Move() command has already set the velocity for this frame.
		// If StopMoving() is not called, the player will continue with this velocity.

	case AttackingSword, AttackingShield, Hurt, Dying, Dead:
		// In these states, zero out user-controlled movement. Knockback is handled separately.
		p.Velocity = core.Vector{}
	}
	return nil
}

func (p *Player) TransitionState(newState PlayerState) {
	if p.state == newState {
		return
	}

	p.state = newState

	if anim := p.GetCurrentAnimation(); anim != nil {
		anim.Reset()
	}

	p.BaseCharacter.Dead = (newState == Dead)
}

// GetActiveDamageSource returns the current attack's DamageSource if the player is attacking
// and the current frame has an active hitbox, otherwise returns nil.
func (p *Player) getActiveDamageSource() *core.DamageSource {
	if p.state != AttackingSword && p.state != AttackingShield {
		return nil
	}

	anim := p.GetCurrentAnimation()
	if anim == nil {
		return nil
	}

	// The current frame index within the *animation slice*
	animIndex := anim.CurrentFrameIndex()

	// Check if we have an attack config for the current direction and animation frame index
	if dirConfigs, ok := p.attackHitboxes[p.direction]; ok {
		if config, ok := dirConfigs[animIndex]; ok {
			// Found an active hitbox config! Create the world-space DamageSource.
			worldHitbox := config.HitBox.Offset(p.Loc.X, p.Loc.Y)

			switch p.state {
			case AttackingSword:
				return core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypePhysical, config.Damage)
			case AttackingShield:
				return core.NewDamageSource(core.TagPlayer, worldHitbox, core.DamageTypeStun, 0)
			default:
				// Other types of attacks handled by projectiles
				return nil
			}
		}
	}

	return nil
}

func (p *Player) HandleHit(ds *core.DamageSource) {
	if ds.Type == core.DamageTypeStun {
		p.ApplyStun(ds.Duration)
	} else {
		p.TakeDamage(ds.Damage)
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), p.Location(), core.KnockbackForce)
		p.ApplyKnockback(force, core.KnockbackDuration)
	}
}

func (p *Player) TakeDamage(damage int) {
	if p.state == Dead || p.state == Dying || p.state == Hurt {
		return
	}

	// Shield blocks damage
	if p.state == AttackingShield {
		return
	}

	p.TransitionState(Hurt)
	p.Health -= damage
}

func (p *Player) ApplyKnockback(force core.Vector, duration int) {
	if p.state == AttackingShield {
		return
	}

	p.BaseCharacter.ApplyKnockback(force, duration)

	if p.state == Dying || p.state == Dead {
		return
	}

	if p.IsKnockedBack() {
		p.TransitionState(Hurt)
	}
}

func (p *Player) ReturnBoomerang() {
	p.hasBoomerang = true
}

// Add methods to satisfy Player interface
func (p *Player) AddExperience(amount int) {
	p.Experience += amount
}

func (p *Player) GetExperience() int {
	return p.Experience
}

func (p *Player) Update(level core.Level, _ core.Player) core.UpdateResult {
	// Handle Knockback Physics.
	p.UpdateKnockback(level)

	// Handle all state transitions.
	stateActions := p.handleState(level)

	p.Velocity.X = p.HandleTileCollisions(level, core.AxisX, p.Velocity.X)
	p.Velocity.Y = p.HandleTileCollisions(level, core.AxisY, p.Velocity.Y)

	// Update visuals
	animation := p.GetCurrentAnimation()
	animation.Update()
	p.Image = p.images[p.state]
	p.SrcRect = p.spriteSheet.Rect(animation.Frame())

	p.bombCooldownTimer.Update()
	p.bowCooldownTimer.Update()
	p.wandCooldownTimer.Update()

	// Return any actions back to the game
	actions := []core.Action{}
	if len(stateActions) > 0 {
		actions = append(actions, stateActions...)
	}
	if ds := p.getActiveDamageSource(); ds != nil {
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			Location:     p.Location(),
			Direction:    core.DirectionToVector(p.direction),
			DamageSource: ds,
		})
	}
	// Reset shieldHeld at the end of the frame so it defaults to false
	// unless set True by input in the next frame.
	p.shieldHeld = false

	return core.UpdateResult{Actions: actions}
}

func (p *Player) CanRemove() bool {
	return false
}
