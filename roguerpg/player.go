package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walking
	Attacking
	Dying
)

type PlayerDirection int

const (
	Left PlayerDirection = iota
	Right
	Up
	Down
)

type CollisionAxis int

const (
	AxisX CollisionAxis = iota
	AxisY
)

type Player struct {
	BaseSprite
	images      map[PlayerState]*ebiten.Image
	spriteSheet *SpriteSheet
	animations  map[PlayerState]map[PlayerDirection]*Animation
	state       PlayerState
	direction   PlayerDirection
	Vx          float64
	Vy          float64
}

func NewPlayer() *Player {
	animations := map[PlayerState]map[PlayerDirection]*Animation{
		Idle: {
			Left:  NewAnimation([]int{6, 7, 8, 9, 10, 11}, 15, true),
			Right: NewAnimation([]int{30, 31, 32, 33, 34, 35}, 15, true),
			Up:    NewAnimation([]int{18, 19, 20, 21, 22, 23}, 15, true),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5}, 15, true),
		},
		Walking: {
			Left:  NewAnimation([]int{6, 7, 8, 9, 10, 11}, 15, true),
			Right: NewAnimation([]int{30, 31, 32, 33, 34, 35}, 15, true),
			Up:    NewAnimation([]int{18, 19, 20, 21, 22, 23}, 15, true),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5}, 15, true),
		},
		Attacking: {
			Left:  NewAnimation([]int{6, 7, 8, 9}, 5, false),
			Right: NewAnimation([]int{30, 31, 32, 33}, 5, false),
			Up:    NewAnimation([]int{18, 19, 20, 21}, 5, false),
			Down:  NewAnimation([]int{0, 1, 2, 3}, 5, false),
		},
		Dying: {
			Left:  NewAnimation([]int{6, 7, 8, 9, 10, 11}, 15, false),
			Right: NewAnimation([]int{30, 31, 32, 33, 34, 35}, 15, false),
			Up:    NewAnimation([]int{18, 19, 20, 21, 22, 23}, 15, false),
			Down:  NewAnimation([]int{0, 1, 2, 3, 4, 5}, 15, false),
		},
	}

	charImages := map[PlayerState]*ebiten.Image{
		Idle:      PlayerIdleSpritesImage,
		Walking:   PlayerWalkSpritesImage,
		Attacking: PlayerAttackSwordSpritesImage,
		Dying:     PlayerDeathSpritesImage,
	}

	spriteSheet := NewSpriteSheet(48, 64, 6, 6)
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
		BaseSprite: BaseSprite{
			Location: Location{
				X: 0,
				Y: 0,
			},
			drawOffset: Location{
				X: 25,
				Y: 38,
			},
			srcRect:    spriteSheet.Rect(0),
			hitbox:     hitbox,
			debugImage: createDebugRectImage(hitbox),
		},
		images:      charImages,
		spriteSheet: spriteSheet,
		animations:  animations,
		state:       Idle,
		direction:   Down,
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

		// Reset the animation for the *new* state
		if anim := p.GetCurrentAnimation(); anim != nil {
			anim.Reset()
		}
	}
}

func (c *Player) Update(level *Level) {
	animation := c.GetCurrentAnimation()
	if animation == nil {
		return
	}

	if c.state == Attacking && animation.IsFinished() {
		c.TransitionState(Idle)
		animation.Reset()
		animation = c.GetCurrentAnimation() // Re-fetch animation after state change
	}

	animation.Update()

	c.image = c.images[c.state]
	c.srcRect = c.spriteSheet.Rect(animation.Frame())

	c.X += c.Vx
	c.HandleTileCollisions(level, AxisX)
	c.Y += c.Vy
	c.HandleTileCollisions(level, AxisY)
}

func (p *Player) HandleUserInput() {
	// If currently attacking, stop movement and skip input processing
	// The Update function is now responsible for transitioning *out* of Attacking.
	if p.state == Attacking {
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
		// Reset movement since attacking is generally stationary
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

func (p *Player) resolveCollision(otherRect Rect, axis CollisionAxis) {
	// TOOD: change HitBox to PushBox
	playerRect := p.HitBox()

	if !playerRect.Intersects(otherRect) {
		return
	}

	if axis == AxisX {
		overlap := 0.0
		if p.Vx > 0 { // moving right
			overlap = playerRect.Right - otherRect.Left
		} else if p.Vx < 0 { // moving left
			overlap = playerRect.Left - otherRect.Right
		}

		if math.Abs(overlap) > 0 {
			p.X -= overlap
			p.Vx = 0.0
		}
	} else if axis == AxisY {
		overlap := 0.0
		if p.Vy > 0 { // moving down
			overlap = playerRect.Bottom - otherRect.Top
		} else if p.Vy < 0 { // moving up
			overlap = playerRect.Top - otherRect.Bottom
		}

		if math.Abs(overlap) > 0 {
			p.Y -= overlap
			p.Vy = 0.0
		}
	}
}

func (p *Player) HandleTileCollisions(level *Level, axis CollisionAxis) {
	// Only check the tiles near the player to improve performance
	playerHitBox := p.HitBox()
	minX := int(math.Floor(playerHitBox.Left/TileSize)) - 1
	maxX := int(math.Floor(playerHitBox.Right/TileSize)) + 1
	minY := int(math.Floor(playerHitBox.Top/TileSize)) - 1
	maxY := int(math.Floor(playerHitBox.Bottom/TileSize)) + 1

	for _, row := range level.Tiles {
		for _, tile := range row {
			if !tile.solid {
				continue
			}

			// Skip over tiles that are not near the player
			tileX := int(tile.X / TileSize)
			tileY := int(tile.Y / TileSize)
			if tileX < minX || tileX > maxX || tileY < minY || tileY > maxY {
				continue
			}

			p.resolveCollision(tile.HitBox(), axis)
		}
	}
}
