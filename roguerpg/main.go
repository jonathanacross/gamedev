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

	KnockbackForce    = 3.0
	KnockbackDuration = 6
)

type Game struct {
	level         *Level
	player        *Player
	camera        *Camera
	damageSources []*DamageSource // current damage sources; saved here for debug drawing
}

func NewGame() *Game {
	level := BuildLevel(70, 50)
	level.AddEnemies()

	player := NewPlayer()
	player.SetLocation(level.FindRandomFloorLocation())
	return &Game{
		level:         level,
		player:        player,
		camera:        NewCamera(ScreenWidth, ScreenHeight),
		damageSources: []*DamageSource{},
	}
}

// handlePlayerAttackCollisions checks any damage source hits enemies
// TODO: also handle player getting hit.
func (g *Game) handleDamageSource(damageSource *DamageSource) {
	// Use an array to track enemies that have been hit in this frame
	// to prevent a single attack frame from hitting one enemy multiple times.
	var hitEnemies []Character

	for _, enemy := range g.level.Enemies {
		if enemy.IsDead() {
			continue
		}

		if damageSource.HitBox.Intersects(enemy.GetHurtBox()) {
			// Check if this enemy has already been hit by this damage source
			alreadyHit := false
			for _, hit := range hitEnemies {
				if hit == enemy {
					alreadyHit = true
					break
				}
			}

			if !alreadyHit {
				enemy.TakeDamage(damageSource.Damage)
				force := CalculateKnockbackForce(g.player.Location(), enemy.Location(), KnockbackForce)
				enemy.ApplyKnockback(force, KnockbackDuration)
				hitEnemies = append(hitEnemies, enemy)
			}
		}
	}
}

func (g *Game) HandleEnemyAttackCollisions() {
	for _, enemy := range g.level.Enemies {
		if enemy.IsDead() {
			continue
		}

		if enemy.GetPushBox().Intersects(g.player.GetHurtBox()) {
			g.player.TakeDamage(1)
			force := CalculateKnockbackForce(enemy.Location(), g.player.Location(), KnockbackForce)
			g.player.ApplyKnockback(force, KnockbackDuration)
			break
		}
	}
}

func (g *Game) Update() error {
	var actions []Action
	for _, object := range g.level.Objects {
		result := object.Update(g.level)
		actions = append(actions, result.Actions...)
	}
	for _, enemy := range g.level.Enemies {
		result := enemy.Update(g.level)
		actions = append(actions, result.Actions...)
	}
	playerResult := g.player.Update(g.level)
	actions = append(actions, playerResult.Actions...)
	g.executeActions(actions)

	for _, ds := range g.damageSources {
		g.handleDamageSource(ds)
	}

	g.HandleEnemyAttackCollisions()

	g.level.Enemies = g.cleanupDeadEnemies(g.level.Enemies)
	g.level.Objects = g.cleanupObjects(g.level.Objects)

	g.camera.CenterOn(g.player.Location())

	return nil
}

func (g *Game) executeActions(actions []Action) {
	g.damageSources = []*DamageSource{}
	for _, action := range actions {
		switch action.Type {
		case ActionCreateDamageSource:
			g.damageSources = append(g.damageSources, action.DamageSource)
		case ActionDropBomb:
			newBomb := NewBomb(action.Location)
			g.level.Objects = append(g.level.Objects, newBomb)
		case ActionExplosion:
			NewBombExplosion := NewBombExplosion(action.Location)
			g.level.Objects = append(g.level.Objects, NewBombExplosion)
		default:
		}
	}
}

func (g *Game) cleanupDeadEnemies(enemies []Character) []Character {
	liveEnemies := enemies[:0] // Create a zero-length slice backed by the original array
	for _, enemy := range enemies {
		if !enemy.IsDead() {
			liveEnemies = append(liveEnemies, enemy)
		}
	}
	// Return the newly filtered slice.
	return liveEnemies
}

func (g *Game) cleanupObjects(objects []GameObject) []GameObject {
	liveObjects := objects[:0] // Create a zero-length slice backed by the original array
	for _, object := range objects {
		if !object.CanRemove() {
			liveObjects = append(liveObjects, object)
		}
	}
	// Return the newly filtered slice.
	return liveObjects
}

func (g *Game) Draw(screen *ebiten.Image) {
	cameraMatrix := g.camera.WorldToScreen()
	viewRect := g.camera.GetViewRect()

	for _, row := range g.level.Tiles {
		for _, tile := range row {
			if tile.GetBounds().Intersects(viewRect) {
				tile.Draw(screen, cameraMatrix)
				if tile.solid {
					tile.DrawDebugInfo(screen, cameraMatrix)
				}
			}
		}
	}

	for _, object := range g.level.Objects {
		if object.GetBounds().Intersects(viewRect) {
			object.Draw(screen, cameraMatrix)
			object.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	for _, enemy := range g.level.Enemies {
		if enemy.GetBounds().Intersects(viewRect) {
			enemy.Draw(screen, cameraMatrix)
			enemy.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	g.player.Draw(screen, cameraMatrix)
	g.player.DrawDebugInfo(screen, cameraMatrix)

	if ShowDebugInfo {
		for _, ds := range g.damageSources {
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
