package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walking
	Attacking
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

type Player struct {
	BaseCharacter
	images         map[PlayerState]*ebiten.Image
	spriteSheet    *SpriteSheet
	animations     map[PlayerState]map[PlayerDirection]*Animation
	state          PlayerState
	direction      PlayerDirection
	Vx             float64
	Vy             float64
	attackHitboxes map[PlayerDirection]map[int]DamageSourceConfig
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
		Attacking: {
			Left:  NewAnimation([]int{8, 9, 10, 11}, 6, false),
			Right: NewAnimation([]int{40, 41, 42, 43}, 6, false),
			Up:    NewAnimation([]int{24, 25, 26, 27}, 6, false),
			Down:  NewAnimation([]int{0, 1, 2, 3}, 6, false),
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
		Idle:      PlayerIdleSpritesImage,
		Walking:   PlayerWalkSpritesImage,
		Attacking: PlayerAttackSwordSpritesImage,
		Hurt:      PlayerHurtSpritesImage,
		Dying:     PlayerDeathSpritesImage,
		Dead:      PlayerDeathSpritesImage,
	}

	spriteSheet := NewSpriteSheet(48, 64, 8, 6)
	// TODO: this is really the "pushbox" for the player;
	// need to make a separate hurtbox for the player, and hitboxes
	// for attacks/weapons.
	hitbox := Rect{
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
				pushBoxOffset: hitbox,
			},
			Health:          8,
			MaxHealth:       8,
			KnockbackFrames: 0,
		},
		images:         charImages,
		spriteSheet:    spriteSheet,
		animations:     animations,
		state:          Idle,
		direction:      Down,
		attackHitboxes: attackHitboxes,
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

func (p *Player) TransitionState(newState PlayerState) {
	if p.state != newState {
		p.state = newState

		if anim := p.GetCurrentAnimation(); anim != nil {
			anim.Reset()
		}
	}
}

// GetActiveDamageSource returns the current attack's DamageSource if the player is attacking
// and the current frame has an active hitbox, otherwise returns nil.
func (p *Player) GetActiveDamageSource() *DamageSource {
	if p.state != Attacking {
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

func (c *Player) TakeDamage(damage int) {
	if c.state == Dead || c.state == Dying || c.state == Hurt {
		return
	}

	c.TransitionState(Hurt)

	c.Health -= damage
	if c.Health < 0 {
		c.TransitionState(Dying)
	}
}

func (p *Player) ApplyKnockback(force Vector, duration int) {
	p.BaseCharacter.ApplyKnockback(force, duration)
	if p.IsKnockedBack() {
		p.TransitionState(Hurt)
	}
}

func (c *Player) Update(level *Level) {
	animation := c.GetCurrentAnimation()
	if animation == nil {
		return
	}

	if c.UpdateKnockback(level) {
		animation.Update()
		c.image = c.images[c.state]
		c.srcRect = c.spriteSheet.Rect(animation.Frame())
		return
	}

	if c.state == Hurt && animation.IsFinished() {
		c.TransitionState(Idle)
		animation = c.GetCurrentAnimation()
	}

	if c.state == Attacking && animation.IsFinished() {
		c.TransitionState(Idle)
		animation = c.GetCurrentAnimation()
	}

	if c.state == Dying && animation.IsFinished() {
		c.TransitionState(Dead)
		animation = c.GetCurrentAnimation()
		return
	}

	animation.Update()

	c.image = c.images[c.state]
	c.srcRect = c.spriteSheet.Rect(animation.Frame())

	c.HandleTileCollisions(level, AxisX, &c.Vx)
	c.HandleTileCollisions(level, AxisY, &c.Vy)
}

func (p *Player) HandleUserInput() {
	// If currently in middle of blocking animation, break out
	if p.state == Attacking || p.state == Hurt || p.state == Dying || p.state == Dead || p.IsKnockedBack() {
		p.Vx = 0
		p.Vy = 0
		return
	}

	moveDir := Vector{X: 0, Y: 0}

	// Handle Movement
	isMoving := false

	// Vertical movement
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		moveDir.Y = -1
		p.direction = Up
		isMoving = true
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		moveDir.Y = 1
		p.direction = Down
		isMoving = true
	}

	// Horizontal movement (updates state if no vertical was pressed,
	// otherwise just changes direction/velocity)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveDir.X = -1
		p.direction = Left
		isMoving = true
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveDir.X = 1
		p.direction = Right
		isMoving = true
	}

	if isMoving {
		p.TransitionState(Walking)
	} else {
		p.TransitionState(Idle)
	}

	// Handle Attack
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// If attack is pressed, it takes precedence over Walking/Idle
		p.TransitionState(Attacking)
		moveDir = Vector{X: 0, Y: 0}
	}

	// Calculate Velocity
	if p.state == Walking {
		walkSpeed := 2.0
		moveDir = moveDir.Normalize().Scale(walkSpeed)
	}
	p.Vx = moveDir.X
	p.Vy = moveDir.Y
}
