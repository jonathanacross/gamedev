package main

import (
	"time"

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
	Direction Direction
}

type Player struct {
	BaseCharacter
	images         map[PlayerState]*ebiten.Image
	spriteSheet    *SpriteSheet
	animations     map[PlayerState]map[Direction]*Animation
	state          PlayerState
	direction      Direction
	Velocity       Vector
	attackHitboxes map[Direction]map[int]DamageSourceConfig // Defines hitboxes for specific animation frames

	bombCooldownTimer *Timer
	bowCooldownTimer  *Timer
	wandCooldownTimer *Timer
	hasBoomerang      bool

	primaryWeapon   WeaponType
	secondaryWeapon WeaponType
}

func NewPlayer() *Player {
	// Define a simple attack hitbox that's only active on the 2nd and 3rd frames (index 1 and 2 in the short animation array)
	attackHitboxes := make(map[Direction]map[int]DamageSourceConfig)

	baseDmg := 1

	// Setup Hitboxes for specific frames of the attack animation.
	// The key (int) is the index within the animation array.

	// Downward swing
	attackHitboxes[Down] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -16, Top: 0, Right: 16, Bottom: 22}, Damage: baseDmg},
		2: {HitBox: Rect{Left: -16, Top: 0, Right: 16, Bottom: 22}, Damage: baseDmg},
	}
	// Left swing
	attackHitboxes[Left] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -22, Top: -24, Right: 0, Bottom: 8}, Damage: baseDmg},
		2: {HitBox: Rect{Left: -22, Top: -24, Right: 0, Bottom: 8}, Damage: baseDmg},
	}
	// Right swing
	attackHitboxes[Right] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: 0, Top: -24, Right: 22, Bottom: 8}, Damage: baseDmg},
		2: {HitBox: Rect{Left: 0, Top: -24, Right: 22, Bottom: 8}, Damage: baseDmg},
	}
	// Upward swing
	attackHitboxes[Up] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -16, Top: -22, Right: 16, Bottom: 0}, Damage: baseDmg},
		2: {HitBox: Rect{Left: -16, Top: -22, Right: 16, Bottom: 0}, Damage: baseDmg},
	}

	ssColumns := 8
	ssRows := 6
	directionOffsets := map[Direction]int{
		Left:  8,
		Right: 40,
		Up:    24,
		Down:  0,
	}

	animations := map[PlayerState]map[Direction]*Animation{
		Idle:            NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 10, true),
		Walking:         NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 10, true),
		AttackingSword:  NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false),
		AttackingShield: NewDirectionAnimationMap([]int{0, 1}, 0, directionOffsets, 6, false),
		AttackingBow:    NewDirectionAnimationMap([]int{0, 1}, 0, directionOffsets, 10, false),
		Hurt:            NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false),
		Dying:           NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 8, false),
		Dead:            NewDirectionAnimationMap([]int{7}, 0, directionOffsets, 8, false),
	}

	charImages := map[PlayerState]*ebiten.Image{
		Idle:            PlayerIdleSpritesImage,
		Walking:         PlayerWalkSpritesImage,
		AttackingSword:  PlayerAttackSwordSpritesImage,
		AttackingShield: PlayerAttackShieldSpritesImage,
		AttackingBow:    PlayerAttackBowSpritesImage,
		AttackingWand:   PlayerIdleSpritesImage, // TODO: add player sprite sheet for wand attack
		Hurt:            PlayerHurtSpritesImage,
		Dying:           PlayerDeathSpritesImage,
		Dead:            PlayerDeathSpritesImage,
	}

	spriteSheet := NewSpriteSheet(48, 64, ssColumns, ssRows)
	pushBox := Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	return &Player{
		BaseCharacter: BaseCharacter{
			BasePhysical: BasePhysical{
				BaseSprite: BaseSprite{
					Location: Location{
						X: 0,
						Y: 0,
					},
					drawOffset: Location{
						X: 25,
						Y: 38,
					},
					srcRect: spriteSheet.Rect(0),
				},
				pushBoxOffset: pushBox,
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
		direction:         Down,
		attackHitboxes:    attackHitboxes,
		bombCooldownTimer: NewTimer(BombCooldown),
		bowCooldownTimer:  NewTimer(BowCooldown),
		wandCooldownTimer: NewTimer(WandCooldown),
		hasBoomerang:      true,
		primaryWeapon:     WeaponSword,
		secondaryWeapon:   WeaponBoomerang,
	}
}

func (p *Player) GetCurrentAnimation() *Animation {
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

func (p *Player) Move(moveVector Vector) {
	// Don't change direction or velocity if attacking
	if p.state != AttackingSword && p.state != AttackingShield && p.state != AttackingBow {
		// Determine facing direction from the move vector.
		// Prioritize vertical direction for diagonal movement.
		if moveVector.Y < 0 {
			p.direction = Up
		} else if moveVector.Y > 0 {
			p.direction = Down
		} else if moveVector.X < 0 {
			p.direction = Left
		} else if moveVector.X > 0 {
			p.direction = Right
		}

		// Set velocity based on the move vector, normalized and scaled.
		p.Velocity = moveVector.Normalize().Scale(PlayerSpeed)
	}
	p.TransitionState(Walking)
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
	// Only start attack if Idle or Walking
	if p.state == Idle || p.state == Walking {
		p.TransitionState(AttackingShield)
	}
}

func (p *Player) UseBomb() *Action {
	// Optional: Block item use while attacking
	// if p.state == AttackingSword || p.state == AttackingShield {
	// 	return nil
	// }

	if p.bombCooldownTimer.IsReady() {
		p.bombCooldownTimer.Reset()

		// Calculate the bomb's spawn location based on player's direction
		loc := Vector(p.Location()).Plus(DirectionToVector(p.direction).Scale(float64(TileSize)))

		return &Action{
			Type:      ActionDropBomb,
			Location:  Location(loc),
			Direction: DirectionToVector(p.direction),
		}
	}
	return nil
}

func (p *Player) UseBow() {
	// Only start attack if Idle or Walking
	if p.state == Idle || p.state == Walking && p.bowCooldownTimer.IsReady() {
		p.TransitionState(AttackingBow)
	}
}

func (p *Player) UseWand(level *Level) *Action {
	if p.wandCooldownTimer.IsReady() {
		p.wandCooldownTimer.Reset()

		// Calculate the bomb's spawn location based on player's direction
		loc := Vector(p.Location()).Plus(DirectionToVector(p.direction).Scale(float64(TileSize)))

		var target Character
		if level != nil {
			target = level.FindNearestEnemy(p.Location(), 6*TileSize)
		}

		return &Action{
			Type:      ActionCreateStar,
			Location:  Location(loc),
			Direction: DirectionToVector(p.direction),
			Target:    target,
		}
	}
	return nil
}

func (p *Player) ShootArrow() *Action {
	if p.bowCooldownTimer.IsReady() {
		p.bowCooldownTimer.Reset()

		// Calculate the bow spawn location based on player's direction
		loc := Vector(p.Location()).Plus(DirectionToVector(p.direction).Scale(float64(TileSize)))

		return &Action{
			Type:      ActionShootArrow,
			Location:  Location(loc),
			Direction: DirectionToVector(p.direction),
		}
	}
	return nil
}

func (p *Player) UseBoomerang() *Action {
	// Optional: Block item use while attacking
	// if p.state == AttackingSword || p.state == AttackingShield {
	// 	return nil
	// }

	if p.hasBoomerang {
		p.hasBoomerang = false

		// Calculate boomerang spawn location
		loc := Vector(p.Location()).Plus(DirectionToVector(p.direction).Scale(float64(TileSize)))

		return &Action{
			Type:      ActionThrowBoomerang,
			Location:  Location(loc),
			Direction: DirectionToVector(p.direction),
		}
	}
	return nil
}

func (p *Player) PrimaryAttack(level *Level) *Action {
	switch p.primaryWeapon {
	case WeaponSword:
		p.AttackSword()
		return nil
	case WeaponBomb:
		return p.UseBomb()
	case WeaponBoomerang:
		return p.UseBoomerang()
	case WeaponShield:
		p.AttackShield()
		return nil
	case WeaponBow:
		p.UseBow()
		return nil
	case WeaponWand:
		return p.UseWand(level)
	default:
		return nil
	}
}

func (p *Player) SecondaryAttack(level *Level) *Action {
	switch p.secondaryWeapon {
	case WeaponSword:
		p.AttackSword()
		return nil
	case WeaponBomb:
		return p.UseBomb()
	case WeaponBoomerang:
		return p.UseBoomerang()
	case WeaponShield:
		p.AttackShield()
		return nil
	case WeaponBow:
		p.UseBow()
		return nil
	case WeaponWand:
		return p.UseWand(level)
	default:
		return nil
	}
}

func (p *Player) SwitchWeapon(weapon WeaponType) {
	p.secondaryWeapon = weapon
}

func (p *Player) IsActive() bool {
	// Input should be blocked if the player is Hurt, Dying, or Dead.
	// We allow input only when Idle, Walking, or Attacking (to allow chaining/canceling).
	return p.state != Hurt && p.state != Dying && p.state != Dead
}

// handleState runs the logic for the Player's current state and determines
// the next state, direction, and velocity (Vx/Vy).
func (p *Player) handleState() []Action {
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
		p.Velocity = Vector{}
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
	if p.state == AttackingShield && animation.IsFinished() {
		p.TransitionState(Idle)
	}
	if p.state == AttackingBow && animation.IsFinished() {
		// Shoot arrow
		if action := p.ShootArrow(); action != nil {
			return []Action{*action}
		}
		p.TransitionState(Idle)
	}

	// Crucially, we no longer call handleMovementInput() here.
	// Movement is now controlled by external command methods (e.g., p.Move()).

	// The Player's internal state handles movement physics based on its state:
	switch p.state {
	case Idle:
		// Idle means no velocity from movement input
		p.Velocity = Vector{}
	case Walking:
		// The Move() command has already set the velocity for this frame.
		// If StopMoving() is not called, the player will continue with this velocity.

	case AttackingSword, AttackingShield, Hurt, Dying, Dead:
		// In these states, zero out user-controlled movement. Knockback is handled separately.
		p.Velocity = Vector{}
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

	p.BaseCharacter.isDead = (newState == Dead)
}

// GetActiveDamageSource returns the current attack's DamageSource if the player is attacking
// and the current frame has an active hitbox, otherwise returns nil.
func (p *Player) getActiveDamageSource() *DamageSource {
	if p.state != AttackingSword && p.state != AttackingShield {
		return nil
	}

	anim := p.GetCurrentAnimation()
	if anim == nil {
		return nil
	}

	// The current frame index within the *animation slice*
	animIndex := anim.frameIndex

	// Check if we have an attack config for the current direction and animation frame index
	if dirConfigs, ok := p.attackHitboxes[p.direction]; ok {
		if config, ok := dirConfigs[animIndex]; ok {
			// Found an active hitbox config! Create the world-space DamageSource.
			worldHitbox := config.HitBox.Offset(p.X, p.Y)

			return NewDamageSource(TagPlayer, worldHitbox, config.Damage)
		}
	}

	return nil
}

func (p *Player) TakeDamage(damage int) {
	if p.state == Dead || p.state == Dying || p.state == Hurt {
		return
	}

	p.TransitionState(Hurt)
	p.Health -= damage
}

func (p *Player) ApplyKnockback(force Vector, duration int) {
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

func (p *Player) Update(level *Level, _ *Player) UpdateResult {
	// Handle Knockback Physics.
	p.UpdateKnockback(level)

	// Handle all state transitions.
	stateActions := p.handleState()

	p.Velocity.X = p.HandleTileCollisions(level, AxisX, p.Velocity.X)
	p.Velocity.Y = p.HandleTileCollisions(level, AxisY, p.Velocity.Y)

	// Update visuals
	animation := p.GetCurrentAnimation()
	animation.Update()
	p.image = p.images[p.state]
	p.srcRect = p.spriteSheet.Rect(animation.Frame())

	p.bombCooldownTimer.Update()
	p.bowCooldownTimer.Update()
	p.wandCooldownTimer.Update()

	// Return any actions back to the game
	actions := []Action{}
	if len(stateActions) > 0 {
		actions = append(actions, stateActions...)
	}
	if ds := p.getActiveDamageSource(); ds != nil {
		actions = append(actions, Action{
			Type:         ActionCreateDamageSource,
			Location:     p.Location(),
			Direction:    DirectionToVector(p.direction),
			DamageSource: ds,
		})
	}
	return UpdateResult{Actions: actions}
}

func (p *Player) CanRemove() bool {
	return false
}
