package main

import (
	"fmt"
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

	// Player state
	Vx         float64
	Vy         float64
	onGround   bool
	facingLeft bool
}

func NewPlayer() *Player {
	spriteSheet := NewGridTileSet(PlayerSprite, TileSize, TileSize, 1, 1)
	return &Player{
		FlippableSprite: FlippableSprite{
			BaseSprite: BaseSprite{
				Location: Location{
					X: 100.0,
					Y: 100.0,
				},
				spriteSheet: spriteSheet,
				srcRect:     spriteSheet.Rect(0),
				hitbox: Rect{ // Initialize the hitbox
					left:   100.0,
					top:    100.0,
					right:  100.0 + TileSize,
					bottom: 100.0 + TileSize,
				},
			},
		},
		Vx:         0.0,
		Vy:         0.0,
		onGround:   false,
		facingLeft: false,
	}
}

// HandleUserInput is a cleaner version using a switch statement.
func (p *Player) HandleUserInput() {
	// Check for movement keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Vx = -RunSpeed
		p.facingLeft = true
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Vx = RunSpeed
		p.facingLeft = false
	} else {
		p.Vx = 0.0
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
func (p *Player) HandleObjectCollisions(objects []GameObject, axis CollisionAxis) PlayerActionEvent {
	playerRect := p.FlippedHitbox()

	for _, obj := range objects {
		if playerRect.Intersects(obj.HitBox()) {
			switch o := obj.(type) {
			case *Spike:
				return PlayerActionEvent{Action: RespawnAction}
			case LevelExit:
				return PlayerActionEvent{Action: SwitchLevelAction, Payload: o}
			case *Checkpoint:
				return PlayerActionEvent{Action: CheckpointReachedAction, Payload: o}
			case *Platform:
				p.resolveCollision(o.HitBox(), axis)
			default:
				fmt.Printf("hit something unknown %T\n", o)
			}
		}
	}

	return PlayerActionEvent{Action: NoAction}
}

// Update moves the player and handles collisions, returning any requested game actions.
func (p *Player) Update(level *Level, gravity float64) PlayerActionEvent {
	p.HandleUserInput()
	p.HandleGravity(gravity)

	p.X += p.Vx
	p.HandleTileCollisions(level, AxisX)
	event := p.HandleObjectCollisions(level.objects, AxisX)
	if event.Action != NoAction {
		return event
	}

	p.onGround = false
	p.Y += p.Vy
	p.HandleTileCollisions(level, AxisY)
	event = p.HandleObjectCollisions(level.objects, AxisY)
	if event.Action != NoAction {
		return event
	}

	return PlayerActionEvent{Action: NoAction}
}
