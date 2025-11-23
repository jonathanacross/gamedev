package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 240

	TileSize = 16

	ShowDebugInfo = false
)

type Game struct {
	level  *Level
	player *Player
	camera *Camera
}

func NewGame() *Game {
	level := BuildLevel(70, 50)
	level.AddEnemies()

	player := NewPlayer()
	player.Location = level.FindRandomFloorLocation()
	return &Game{
		level:  level,
		player: player,
		camera: NewCamera(ScreenWidth, ScreenHeight),
	}
}

// handlePlayerAttackCollisions checks if the player's active attack hits any enemies.
func (g *Game) handlePlayerAttackCollisions() {
	damageSource := g.player.GetActiveDamageSource()
	if damageSource == nil {
		return // Player is not currently attacking or hitbox is inactive
	}

	// Use an array to track enemies that have been hit in this frame
	// to prevent a single attack frame from hitting one enemy multiple times.
	var hitEnemies []*BlobEnemy

	for _, enemy := range g.level.Enemies {
		if enemy.IsDead {
			continue
		}

		// Check for intersection between the DamageSource's HitBox and the enemy's HurtBox (which is its HitBox())
		// TODO: why don't we need to use damageSource.HitBox()?
		if damageSource.HitBox.Intersects(enemy.HitBox()) {
			// Check if this enemy has already been hit by this damage source
			alreadyHit := false
			for _, hit := range hitEnemies {
				if hit == enemy {
					alreadyHit = true
					break
				}
			}

			if !alreadyHit {
				// Apply damage
				enemy.TakeDamage(damageSource.Damage)
				hitEnemies = append(hitEnemies, enemy)
			}
		}
	}
}

func (g *Game) HandleEnemyAttackCollisions() {
	for _, enemy := range g.level.Enemies {
		if enemy.IsDead {
			continue
		}

		if enemy.HitBox().Intersects(g.player.HitBox()) {
			g.player.TakeDamage(1)
			break
		}
	}
}

func (g *Game) Update() error {
	for _, enemy := range g.level.Enemies {
		enemy.Update(g.level)
	}
	g.HandleEnemyAttackCollisions()
	g.player.HandleUserInput()
	g.player.Update(g.level)
	g.handlePlayerAttackCollisions()
	g.level.Enemies = g.cleanupDeadEnemies(g.level.Enemies)

	g.camera.CenterOn(g.player.Location)

	return nil
}

// cleanupDeadEnemies iterates through the slice and removes enemies marked for deletion.
func (g *Game) cleanupDeadEnemies(enemies []*BlobEnemy) []*BlobEnemy {
	liveEnemies := enemies[:0] // Creates a zero-length slice backed by the original array
	for _, enemy := range enemies {
		// For now, we only remove if IsDead is true. You might add a death animation check here later.
		if !enemy.IsDead {
			liveEnemies = append(liveEnemies, enemy)
		}
	}
	// Return the newly filtered slice.
	return liveEnemies
}

func (g *Game) Draw(screen *ebiten.Image) {
	cameraMatrix := g.camera.WorldToScreen()
	viewRect := g.camera.GetViewRect()

	for _, row := range g.level.Tiles {
		for _, tile := range row {
			if tile.HitBox().Intersects(viewRect) {
				tile.Draw(screen, cameraMatrix)
				if tile.solid {
					tile.DrawDebugInfo(screen, cameraMatrix)
				}
			}
		}
	}
	for _, enemy := range g.level.Enemies {
		if enemy.HitBox().Intersects(viewRect) {
			enemy.Draw(screen, cameraMatrix)
			enemy.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	g.player.Draw(screen, cameraMatrix)
	g.player.DrawDebugInfo(screen, cameraMatrix)

	// Draw active player attack hitbox for debugging
	if ShowDebugInfo {
		if ds := g.player.GetActiveDamageSource(); ds != nil {
			ds.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	DrawHeadsUpDisplay(screen, g.player)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(ScreenWidth*3, ScreenHeight*3)
	ebiten.SetWindowTitle("Rogue RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
