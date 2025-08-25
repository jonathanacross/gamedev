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

type Player struct {
	BaseSprite

	// Player state
	Vx               float64
	Vy               float64
	onGround         bool
	checkpointId     int
	activeCheckpoint *Checkpoint
}

func NewPlayer() *Player {
	spriteSheet := NewSpriteSheet(PlayerSprite, TileSize, TileSize, 1, 1)
	return &Player{
		BaseSprite: BaseSprite{
			Location: Location{
				X: 100.0,
				Y: 100.0,
			},
			spriteSheet: spriteSheet,
			srcRect:     spriteSheet.Rect(0),
		},
		Vx:               0.0,
		Vy:               0.0,
		onGround:         false,
		activeCheckpoint: nil,
	}
}

// HandleUserInput is a cleaner version using a switch statement.
func (p *Player) HandleUserInput() {
	runSpeed := 100.0 / float64(ebiten.TPS())

	// Check for movement keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Vx = -runSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Vx = runSpeed
	} else {
		p.Vx = 0.0
	}
}

func (p *Player) HandleGravity(gravity float64) {
	p.Vy += gravity
}

func (p *Player) IsOnGround() bool {
	return p.onGround
}

// ResolveCollision handles collision resolution on a specified axis.
func (p *Player) ResolveCollision(tile Tile, axis CollisionAxis) {
	playerRect := p.HitRect()
	tileRect := tile.HitRect()

	if !playerRect.Intersects(&tileRect) {
		return
	}

	if axis == AxisX {
		overlap := 0.0
		if p.Vx > 0 { // moving right
			overlap = playerRect.right - tileRect.left
		} else if p.Vx < 0 { // moving left
			overlap = playerRect.left - tileRect.right
		}

		if math.Abs(overlap) > 0 {
			p.X -= overlap
			p.Vx = 0.0
		}
	} else if axis == AxisY {
		overlap := 0.0
		if p.Vy > 0 { // moving down
			overlap = playerRect.bottom - tileRect.top
		} else if p.Vy < 0 { // moving up
			overlap = playerRect.top - tileRect.bottom
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
	playerHitRect := p.HitRect()
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

		p.ResolveCollision(tile, axis)
	}
}

func (p *Player) HandleCheckpoints(level *Level) {
	playerHitRect := p.HitRect()

	for _, cp := range level.checkpoints {
		if playerHitRect.Intersects(&cp.hitbox) {
			if p.activeCheckpoint != nil && p.activeCheckpoint.Id != cp.Id {
				if p.activeCheckpoint.LevelNum == cp.LevelNum {
					p.activeCheckpoint.SetActive(false)
				}
			}
			cp.SetActive(true)
			p.activeCheckpoint = cp
			p.checkpointId = cp.Id
		}
	}
}

// HandleSpikeCollisions now returns a boolean indicating if a respawn is needed.
func (p *Player) HandleSpikeCollisions(level *Level) bool {
	playerRect := p.HitRect()
	for _, spike := range level.spikes {
		if playerRect.Intersects(&spike.hitbox) {
			return true // Player needs to respawn
		}
	}
	return false // No respawn needed
}

func (p *Player) checkLevelExits(level *Level) PlayerActionEvent {
	playerHitRect := p.HitRect()

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

	p.onGround = false
	p.Y += p.Vy
	p.HandleCollisions(level, AxisY)

	p.HandleCheckpoints(level)

	if p.HandleSpikeCollisions(level) {
		return PlayerActionEvent{Action: RespawnAction}
	}

	return p.checkLevelExits(level)
}
