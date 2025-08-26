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
// It now takes a generic Rect, allowing it to work with any object.
func (p *Player) resolveCollision(otherRect Rect, axis CollisionAxis) {
	playerRect := p.FlippedHitbox()

	if !playerRect.Intersects(&otherRect) {
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
func (p *Player) HandleCollisions(level *Level, axis CollisionAxis) {
	// Only check the tiles near the player to improve performance
	playerHitRect := p.FlippedHitbox()
	minX := int(math.Floor(playerHitRect.left/TileSize)) - 1
	maxX := int(math.Floor(playerHitRect.right/TileSize)) + 1
	minY := int(math.Floor(playerHitRect.top/TileSize)) - 1
	maxY := int(math.Floor(playerHitRect.bottom/TileSize)) + 1

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

		p.resolveCollision(tile.HitRect(), axis)
	}
}

func (p *Player) HandlePlatformCollisions(platforms []*Platform, axis CollisionAxis) {
	playerHitRect := p.FlippedHitbox()

	// Check for collision with each platform
	for _, platform := range platforms {
		platformRect := platform.hitbox
		if playerHitRect.Intersects(&platformRect) {
			p.resolveCollision(platform.HitRect(), axis)
		}
	}
}

func (p *Player) CheckCheckpointCollisions(level *Level) *Checkpoint {
	playerHitRect := p.FlippedHitbox()

	for _, cp := range level.checkpoints {
		if playerHitRect.Intersects(&cp.hitbox) {
			return cp
		}
	}
	return nil
}

// HandleSpikeCollisions now returns a boolean indicating if a respawn is needed.
func (p *Player) HandleSpikeCollisions(level *Level) bool {
	playerRect := p.FlippedHitbox()
	for _, spike := range level.spikes {
		if playerRect.Intersects(&spike.hitbox) {
			return true // Player needs to respawn
		}
	}
	return false // No respawn needed
}

func (p *Player) checkLevelExits(level *Level) PlayerActionEvent {
	playerHitRect := p.FlippedHitbox()

	for _, exit := range level.exits {
		if playerHitRect.Intersects(&exit.Rect) {
			return PlayerActionEvent{Action: SwitchLevelAction, Payload: exit}
		}
	}
	return PlayerActionEvent{Action: NoAction}
}

// Update moves the player and handles collisions, returning any requested game actions.
func (p *Player) Update(level *Level, gravity float64) PlayerActionEvent {
	p.HandleUserInput()
	p.HandleGravity(gravity)

	p.X += p.Vx
	p.HandleCollisions(level, AxisX)
	p.HandlePlatformCollisions(level.platforms, AxisX)

	p.onGround = false
	p.Y += p.Vy
	p.HandleCollisions(level, AxisY)
	p.HandlePlatformCollisions(level.platforms, AxisY)

	if newCheckpoint := p.CheckCheckpointCollisions(level); newCheckpoint != nil {
		return PlayerActionEvent{Action: CheckpointReachedAction, Payload: newCheckpoint}
	}

	if p.HandleSpikeCollisions(level) {
		return PlayerActionEvent{Action: RespawnAction}
	}

	return p.checkLevelExits(level)
}
