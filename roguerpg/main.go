package main

import (
	"log"
	"math/rand/v2"

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
	level.AddObjects()
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

// checkDamageAgainstTargets checks a damage source against a list of characters (enemies or player)
// and applies damage and knockback, returning the characters that were hit.
func (g *Game) checkDamageAgainstTargets(damageSource *DamageSource, targets []Character) (hitCharacters []Character) {
	for _, target := range targets {
		if target.IsDead() {
			continue
		}

		if damageSource.HitBox.Intersects(target.GetHurtBox()) {
			// Apply damage and knockback
			target.TakeDamage(damageSource.Damage)

			// The attacker's location for knockback calculation is the center of the HitBox.
			// This is better than the Character location for area-of-effect attacks (like bombs).
			attackerLoc := Location{
				X: (damageSource.HitBox.Left + damageSource.HitBox.Right) / 2,
				Y: (damageSource.HitBox.Top + damageSource.HitBox.Bottom) / 2,
			}

			force := CalculateKnockbackForce(attackerLoc, target.Location(), KnockbackForce)
			target.ApplyKnockback(force, KnockbackDuration)
			hitCharacters = append(hitCharacters, target)
		}
	}
	return hitCharacters
}

// handleDamageSource checks if a damage source hits the player or any enemy,
// depending on its SourceTag.
func (g *Game) handleDamageSource(damageSource *DamageSource) {
	if damageSource.SourceTag == TagPlayer {
		// Player attack hits enemies
		g.checkDamageAgainstTargets(damageSource, g.level.Enemies)
	} else if damageSource.SourceTag == TagEnemy {
		// Enemy attack hits the player
		// We only check the player here.
		g.checkDamageAgainstTargets(damageSource, []Character{g.player})
	}

	// TODO: Handle other tags like TagBomb, which could hit both the player and enemies.
}

func (g *Game) Update() error {
	var actions []Action
	for _, object := range g.level.Objects {
		result := object.Update(g.level, g.player)
		actions = append(actions, result.Actions...)
	}
	for _, enemy := range g.level.Enemies {
		result := enemy.Update(g.level, g.player)
		actions = append(actions, result.Actions...)
	}
	playerResult := g.player.Update(g.level, g.player)
	actions = append(actions, playerResult.Actions...)
	g.executeActions(actions)

	for _, ds := range g.damageSources {
		g.handleDamageSource(ds)
	}

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
		case ActionThrowBoomerang:
			newBoomerang := NewBoomerang(action.Location, action.Direction, rand.IntN(3)+1)
			g.level.Objects = append(g.level.Objects, newBoomerang)
		case ActionReturnBoomerang:
			g.player.ReturnBoomerang()
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
