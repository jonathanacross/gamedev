package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	BaseSprite

	// Player state
	Vx               float64
	Vy               float64
	onGround         bool
	checkpointId     int
	game             *Game
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

func (p *Player) FixPlayerX(tile *Tile) {
	playerRect := p.HitRect()
	tileRect := tile.HitRect()
	if !playerRect.Intersects(&tileRect) {
		return
	}

	if p.Vx > 0 { // player moving right
		if playerRect.right < tileRect.right { // and to the left of the tile
			overlap := playerRect.right - tileRect.left
			p.X -= float64(overlap)
			p.Vx = 0.0
		}
	} else if p.Vx < 0 { // player moving left
		if playerRect.left > tileRect.left { // and to the right of the tilec
			overlap := tileRect.right - playerRect.left
			p.X += float64(overlap)
			p.Vx = 0.0
		}
	}
}

func (p *Player) FixPlayerY(tile *Tile) {
	playerRect := p.HitRect()
	tileRect := tile.HitRect()
	if !playerRect.Intersects(&tileRect) {
		return
	}

	if p.Vy > 0 { // player moving down
		if playerRect.bottom < tileRect.bottom { // and above the tile
			overlap := playerRect.bottom - tileRect.top
			p.Y -= overlap
			p.Vy = 0.0
			p.onGround = true
		}
	} else if p.Vy < 0 { // player moving up
		if playerRect.top > tileRect.top { // and below the tile
			overlap := tileRect.bottom - playerRect.top
			p.Y += overlap
			p.Vy = 0.0
			p.onGround = true
		}
	}
}

func (p *Player) HandleCollisions(level *Level, horiz bool) {
	for _, tile := range *level.tiles {
		if !tile.solid {
			continue
		}
		if horiz {
			p.FixPlayerX(&tile)
		} else {
			p.FixPlayerY(&tile)
		}
	}
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
	// enforce a max speed so the player can't fall
	// so fast they could fall through an entire tile
	p.Vy = clamp(p.Vy, -5.0, 5.0)
}

func (p *Player) HandleUserInput() {
	runSpeed := 100.0 / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Vx = -runSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Vx = runSpeed
	} else {
		p.Vx = 0.0
	}
}
func (p *Player) HandleCheckpoints() {
	for _, cp := range p.game.currentLevel.checkpoints {
		if p.HitRect().Intersects(&cp.hitbox) {
			if !cp.Active {
				if p.activeCheckpoint != nil {
					p.activeCheckpoint.SetActive(false)
				}
				cp.SetActive(true)
				p.activeCheckpoint = cp
				p.checkpointId = cp.Id
			}
		}
	}
}

func (p *Player) HandleSpikeCollisions() {
	playerRect := p.HitRect()
	for _, spike := range p.game.currentLevel.spikes {
		if playerRect.Intersects(&spike.hitbox) {
			p.game.Respawn()
			return
		}
	}
}

func (p *Player) Update(game *Game) {
	p.HandleUserInput()
	p.HandleGravity(game.gravity)
	p.X += p.Vx
	p.HandleCollisions(game.currentLevel, true)
	p.onGround = false
	p.Y += p.Vy
	p.HandleCollisions(game.currentLevel, false)
	p.HandleCheckpoints()
	p.HandleSpikeCollisions()
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.BaseSprite.Draw(screen)
}

func (p *Player) IsOnGround() bool {
	return p.onGround
}
