package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// CollisionAxis defines the axis of a collision.
type CollisionAxis int

const (
	AxisX CollisionAxis = iota
	AxisY
)

type PlayerState int

const (
	Walking = iota
	Idle
)

type Player struct {
	FlippableSprite
	spriteSheet *GridTileSet
	animations  map[PlayerState]*Animation

	// Player state
	Vx         float64
	Vy         float64
	onGround   bool
	facingLeft bool
	state      PlayerState

	numCrystals int
	numDeaths   int
}

func NewPlayer() *Player {
	spriteSheet := NewGridTileSet(TileSize, TileSize, 4, 1)

	return &Player{
		FlippableSprite: FlippableSprite{
			BaseSprite: BaseSprite{
				Location: Location{
					X: 0,
					Y: 0,
				},
				image:   PlayerSprite,
				srcRect: spriteSheet.Rect(0),
				hitbox: Rect{
					left:   3,
					top:    5,
					right:  13,
					bottom: 16,
				},
			},
		},
		spriteSheet: spriteSheet,
		animations: map[PlayerState]*Animation{
			Walking: NewAnimation(0, 3, 10),
			Idle:    NewAnimation(1, 1, 100),
		},
		Vx:          0.0,
		Vy:          0.0,
		onGround:    false,
		facingLeft:  false,
		state:       Idle,
		numCrystals: 0,
	}
}

func (p *Player) Reset() {
	p.numCrystals = 0
	p.numDeaths = 0
}

func (p *Player) Draw(screen *ebiten.Image, debug bool) {
	currSpriteFrame := p.animations[p.state].Frame()
	p.srcRect = p.spriteSheet.Rect(currSpriteFrame)
	p.flipHoriz = p.facingLeft
	p.FlippableSprite.Draw(screen, debug)
}

// HandleUserInput is a cleaner version using a switch statement.
func (p *Player) HandleUserInput() {
	// Check for movement keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Vx = -RunSpeed
		p.facingLeft = true
		p.state = Walking
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Vx = RunSpeed
		p.facingLeft = false
		p.state = Walking
	} else {
		p.Vx = 0.0
		p.state = Idle
	}

	p.flipHoriz = p.facingLeft
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (p *Player) HandleGravity(gravity float64) {
	p.Vy += gravity
	p.Vy = clamp(p.Vy, -MaxFallSpeed, MaxFallSpeed)
	p.flipVert = gravity < 0
}

func (p *Player) IsOnGround() bool {
	return p.onGround
}

// resolveCollision handles collision resolution on a specified axis.
func (p *Player) resolveCollision(otherRect Rect, axis CollisionAxis) {
	playerRect := p.FlippedHitbox()

	if !playerRect.Intersects(otherRect) {
		return
	}

	if axis == AxisX {
		overlap := 0.0
		if p.Vx > 0 { // moving right
			overlap = playerRect.right - otherRect.left
		} else if p.Vx < 0 { // moving left
			overlap = playerRect.left - otherRect.right
		}

		if math.Abs(overlap) > 0 {
			p.X -= overlap
			p.Vx = 0.0
		}
	} else if axis == AxisY {
		overlap := 0.0
		if p.Vy > 0 { // moving down
			overlap = playerRect.bottom - otherRect.top
		} else if p.Vy < 0 { // moving up
			overlap = playerRect.top - otherRect.bottom
		}

		if math.Abs(overlap) > 0 {
			p.Y -= overlap
			p.Vy = 0.0
			p.onGround = true
		}
	}
}

// HandleCollisions checks for and resolves collisions for the player.
func (p *Player) HandleTileCollisions(level *Level, axis CollisionAxis) {
	// Only check the tiles near the player to improve performance
	playerHitBox := p.FlippedHitbox()
	minX := int(math.Floor(playerHitBox.left/TileSize)) - 1
	maxX := int(math.Floor(playerHitBox.right/TileSize)) + 1
	minY := int(math.Floor(playerHitBox.top/TileSize)) - 1
	maxY := int(math.Floor(playerHitBox.bottom/TileSize)) + 1

	for _, tile := range level.tiles {
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

// HandleObjectCollisions checks for collisions with a slice of GameObjects and handles them.
func (p *Player) HandleObjectCollisions(objects []GameObject, axis CollisionAxis) {
	playerRect := p.FlippedHitbox()

	for _, obj := range objects {
		if platform, ok := obj.(*Platform); ok {
			if playerRect.Intersects(obj.HitBox()) {
				p.resolveCollision(platform.HitBox(), axis)
				// Move the player along with the platform if on ground
				if axis == AxisY && p.onGround && platform.horiz {
					p.X += platform.delta
				}
				// TODO: Handle vertical platforms
			}
		} else if breakingFloor, ok := obj.(*BreakingFloor); ok {
			if playerRect.Intersects(obj.HitBox()) {
				if breakingFloor.IsSolid() {
					p.resolveCollision(breakingFloor.HitBox(), axis)
					breakingFloor.StartBreak()
				} else {
					breakingFloor.KeepBroken()
				}
			}
		} else if crystal, ok := obj.(*Crystal); ok {
			if playerRect.Intersects(obj.HitBox()) && !crystal.Collected {
				crystal.Collected = true
				p.numCrystals++
			}
		}
	}
}

// checkAllEvents checks for collisions with event objects (non-solid) and returns an event action.
func (p *Player) checkAllEvents(level *Level) PlayerActionEvent {
	playerRect := p.FlippedHitbox()

	for _, obj := range level.objects {
		if playerRect.Intersects(obj.HitBox()) {
			switch o := obj.(type) {
			case *Spike:
				return PlayerActionEvent{Action: RespawnAction}
			case LevelExit:
				return PlayerActionEvent{Action: SwitchLevelAction, Payload: o}
			case *Checkpoint:
				return PlayerActionEvent{Action: CheckpointReachedAction, Payload: o}
			case *HelicopterMonster:
				return PlayerActionEvent{Action: RespawnAction}
			}
		}
	}
	if p.numCrystals >= NumCrystals {
		return PlayerActionEvent{Action: WinGameAction}
	}
	return PlayerActionEvent{Action: NoAction}
}

// Update moves the player and handles collisions, returning any requested game actions.
func (p *Player) Update(level *Level, gravity float64) PlayerActionEvent {
	p.animations[p.state].Update()

	p.HandleUserInput()
	p.HandleGravity(gravity)

	p.onGround = false

	p.X += p.Vx
	p.HandleTileCollisions(level, AxisX)
	p.HandleObjectCollisions(level.objects, AxisX)

	p.Y += p.Vy
	p.HandleTileCollisions(level, AxisY)
	p.HandleObjectCollisions(level.objects, AxisY)

	event := p.checkAllEvents(level)
	if event.Action != NoAction {
		return event
	}

	return PlayerActionEvent{Action: NoAction}
}
