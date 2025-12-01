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
	Hurt
	Dying
	Dead
)

type PlayerDirection int

const (
	Left PlayerDirection = iota
	Right
	Up
	Down
)

const (
	PlayerSpeed  = 2.0
	BombCooldown = 750 * time.Millisecond
)

type PlayerAnimationKey struct {
	State     PlayerState
	Direction PlayerDirection
}

type Player struct {
	BaseCharacter
	images         map[PlayerState]*ebiten.Image
	spriteSheet    *SpriteSheet
	animations     map[PlayerState]map[PlayerDirection]*Animation
	state          PlayerState
	direction      PlayerDirection
	Velocity       Vector
	attackHitboxes map[PlayerDirection]map[int]DamageSourceConfig // Defines hitboxes for specific animation frames

	bombCooldownTimer *Timer
	hasBoomerang      bool

	primaryWeapon   WeaponType
	secondaryWeapon WeaponType
}

func NewPlayer() *Player {
	// Define a simple attack hitbox that's only active on the 2nd and 3rd frames (index 1 and 2 in the short animation array)
	attackHitboxes := make(map[PlayerDirection]map[int]DamageSourceConfig)

	baseDmg := 1

	// Setup Hitboxes for specific frames of the attack animation.
	// The key (int) is the index within the animation array.

	// Downward swing: Attack box is in front (down) of the player
	attackHitboxes[Down] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -8, Top: 10, Right: 8, Bottom: 30}, Damage: baseDmg},   // Frame 1 (sword out)
		2: {HitBox: Rect{Left: -16, Top: 15, Right: 16, Bottom: 35}, Damage: baseDmg}, // Frame 2 (full extension)
	}
	// Left swing: Attack box is to the left
	attackHitboxes[Left] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -30, Top: -8, Right: -10, Bottom: 8}, Damage: baseDmg},
		2: {HitBox: Rect{Left: -35, Top: -16, Right: -15, Bottom: 16}, Damage: baseDmg},
	}
	// Right swing: Attack box is to the right
	attackHitboxes[Right] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: 10, Top: -8, Right: 30, Bottom: 8}, Damage: baseDmg},
		2: {HitBox: Rect{Left: 15, Top: -16, Right: 35, Bottom: 16}, Damage: baseDmg},
	}
	// Upward swing: Attack box is above the player
	attackHitboxes[Up] = map[int]DamageSourceConfig{
		1: {HitBox: Rect{Left: -8, Top: -30, Right: 8, Bottom: -10}, Damage: baseDmg},
		2: {HitBox: Rect{Left: -16, Top: -35, Right: 16, Bottom: -15}, Damage: baseDmg},
	}

	animations := map[PlayerState]map[PlayerDirection]*Animation{
		Idle: {
			Left:  NewAnimation([]int{8, 9, 10, 11, 12, 13, 14, 15}, 10, true),
			Right: NewAnimation([]int{40, 41, 42, 43, 44, 45, 46, 47}, 10, true),
			Up:    NewAnimation([]int{24, 25, 26, 27, 28, 29, 30, 31}, 10, true),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5, 6, 7}, 10, true),
		},
		Walking: {
			Left:  NewAnimation([]int{8, 9, 10, 11, 12, 13, 14, 15}, 10, true),
			Right: NewAnimation([]int{40, 41, 42, 43, 44, 45, 46, 47}, 10, true),
			Up:    NewAnimation([]int{24, 25, 26, 27, 28, 29, 30, 31}, 10, true),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5, 6, 7}, 10, true),
		},
		AttackingSword: {
			Left:  NewAnimation([]int{8, 9, 10, 11}, 6, false),
			Right: NewAnimation([]int{40, 41, 42, 43}, 6, false),
			Up:    NewAnimation([]int{24, 25, 26, 27}, 6, false),
			Down:  NewAnimation([]int{0, 1, 2, 3}, 6, false),
		},
		AttackingShield: {
			Left:  NewAnimation([]int{8, 9}, 6, false),
			Right: NewAnimation([]int{40, 41}, 6, false),
			Up:    NewAnimation([]int{24, 25}, 6, false),
			Down:  NewAnimation([]int{0, 1}, 6, false),
		},
		Hurt: {
			Left:  NewAnimation([]int{8, 9, 10, 11}, 10, false),
			Right: NewAnimation([]int{40, 41, 42, 43}, 10, false),
			Up:    NewAnimation([]int{24, 25, 26, 27}, 10, false),
			Down:  NewAnimation([]int{0, 1, 2, 3}, 10, false),
		},
		Dying: {
			Left:  NewAnimation([]int{8, 9, 10, 11, 12, 13, 14, 15}, 8, false),
			Right: NewAnimation([]int{40, 41, 42, 43, 44, 45, 46, 47}, 8, false),
			Up:    NewAnimation([]int{24, 25, 26, 27, 28, 29, 30, 31}, 8, false),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5, 6, 7}, 8, false),
		},
		Dead: {
			Left:  NewAnimation([]int{15}, 100, true),
			Right: NewAnimation([]int{47}, 100, true),
			Up:    NewAnimation([]int{31}, 100, true),
			Down:  NewAnimation([]int{7}, 100, true),
		},
	}

	charImages := map[PlayerState]*ebiten.Image{
		Idle:            PlayerIdleSpritesImage,
		Walking:         PlayerWalkSpritesImage,
		AttackingSword:  PlayerAttackSwordSpritesImage,
		AttackingShield: PlayerAttackShieldSpritesImage,
		Hurt:            PlayerHurtSpritesImage,
		Dying:           PlayerDeathSpritesImage,
		Dead:            PlayerDeathSpritesImage,
	}

	spriteSheet := NewSpriteSheet(48, 64, 8, 6)
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
			KnockbackFrames: 0,
		},
		images:            charImages,
		spriteSheet:       spriteSheet,
		animations:        animations,
		state:             Idle,
		direction:         Down,
		attackHitboxes:    attackHitboxes,
		bombCooldownTimer: NewTimer(BombCooldown),
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
	if p.state != AttackingSword && p.state != AttackingShield {
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
		loc := Vector(p.Location()).Plus(p.getDirectionVector().Scale(float64(TileSize)))

		return &Action{
			Type:      ActionDropBomb,
			Location:  Location(loc),
			Direction: p.getDirectionVector(),
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
		loc := Vector(p.Location()).Plus(p.getDirectionVector().Scale(float64(TileSize)))

		return &Action{
			Type:      ActionThrowBoomerang,
			Location:  Location(loc),
			Direction: p.getDirectionVector(),
		}
	}
	return nil
}

func (p *Player) PrimaryAttack() *Action {
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
	default:
		return nil
	}
}

func (p *Player) SecondaryAttack() *Action {
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
func (p *Player) handleState() {
	animation := p.GetCurrentAnimation()
	if animation == nil {
		p.state = Idle // Should never happen
		return
	}

	if p.Health <= 0 && p.state != Dying && p.state != Dead {
		p.TransitionState(Dying)
		return
	}

	// Check for terminal state completion.
	if p.state == Dying && animation.IsFinished() {
		p.TransitionState(Dead)
		return
	}

	// Once Dead, stop all logic.
	if p.state == Dead {
		p.Velocity = Vector{}
		return
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

// TODO: move this, and Direction, to utility function
func (p *Player) getDirectionVector() Vector {
	switch p.direction {
	case Left:
		return Vector{X: -1, Y: 0}
	case Right:
		return Vector{X: 1, Y: 0}
	case Up:
		return Vector{X: 0, Y: -1}
	case Down:
		return Vector{X: 0, Y: 1}
	}
	return Vector{X: 0, Y: 0}
}

func (p *Player) Update(level *Level, _ *Player) UpdateResult {
	// Handle Knockback Physics.
	p.UpdateKnockback(level)

	// Handle all state transitions.
	p.handleState()

	p.Velocity.X = p.HandleTileCollisions(level, AxisX, p.Velocity.X)
	p.Velocity.Y = p.HandleTileCollisions(level, AxisY, p.Velocity.Y)

	// Update visuals
	animation := p.GetCurrentAnimation()
	animation.Update()
	p.image = p.images[p.state]
	p.srcRect = p.spriteSheet.Rect(animation.Frame())

	p.bombCooldownTimer.Update()

	// Return any actions back to the game
	actions := []Action{}
	if ds := p.getActiveDamageSource(); ds != nil {
		actions = append(actions, Action{
			Type:         ActionCreateDamageSource,
			Location:     p.Location(),
			Direction:    p.getDirectionVector(),
			DamageSource: ds,
		})
	}
	return UpdateResult{Actions: actions}
}

func (p *Player) CanRemove() bool {
	return false
}
