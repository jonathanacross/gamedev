package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// struct to handle the main game logic
type MainGameState struct {
}

// checkDamageAgainstTargets checks a damage source against a list of characters (enemies or player)
// and applies damage and knockback, returning the characters that were hit.
func (g *MainGameState) checkDamageAgainstTargets(damageSource *DamageSource, targets []Character) (hitCharacters []Character) {
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

func (s *MainGameState) handleDamageSource(ctx *GameContext, damageSource *DamageSource) {
	if damageSource.SourceTag == TagPlayer {
		// Player attack hits enemies
		s.checkDamageAgainstTargets(damageSource, ctx.Level.Enemies)
	} else if damageSource.SourceTag == TagEnemy {
		// Enemy attack hits the player
		s.checkDamageAgainstTargets(damageSource, []Character{ctx.Player})
	}

	// TODO: Handle other tags like TagBomb, which could hit both the player and enemies.
}

func (mg *MainGameState) handleInput(ctx *GameContext) []Action {
	var actions []Action
	player := ctx.Player
	isMoving := false

	// --- Handle State Transition Input (Global to MainGameState) ---

	// Open Weapon Selector Menu
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		// This remains an action to modify the global game state stack
		action := Action{
			Type:      ActionPushState,
			GameState: NewWeaponSelector(),
		}
		actions = append(actions, action)
	}

	if !player.IsActive() {
		return actions
	}

	// --- Handle Movement Input ---
	var moveVector Vector
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		moveVector.Y = -1
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		moveVector.Y = 1
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveVector.X = -1
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveVector.X = 1
		isMoving = true
	}

	if !isMoving {
		player.StopMoving()
	} else {
		player.Move(moveVector)
	}

	// --- Handle Attack/Item Input (Take precedence over movement commands) ---

	// Sword Attack
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		player.AttackSword()
	}

	// Shield Attack
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		player.AttackShield()
	}

	if ebiten.IsKeyPressed(ebiten.KeyB) {
		// Call the new method and collect the resulting action if it's not nil
		if bombAction := player.UseBomb(); bombAction != nil {
			actions = append(actions, *bombAction)
		}
	}

	// Boomerang
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		// Call the new method and collect the resulting action if it's not nil
		if boomerangAction := player.UseBoomerang(); boomerangAction != nil {
			actions = append(actions, *boomerangAction)
		}
	}

	return actions
}

func (mg *MainGameState) Update(ctx *GameContext) []Action {

	var actions []Action

	// Handle Input (New Step)
	inputActions := mg.handleInput(ctx)
	actions = append(actions, inputActions...)

	// Update Game Objects
	for _, object := range ctx.Level.Objects {
		result := object.Update(ctx.Level, ctx.Player)
		actions = append(actions, result.Actions...)
	}

	// Update Enemies
	for _, enemy := range ctx.Level.Enemies {
		// Only update enemies if they are not currently being knocked back,
		// otherwise, the knockback movement should be prioritized.
		if !enemy.IsKnockedBack() {
			result := enemy.Update(ctx.Level, ctx.Player)
			actions = append(actions, result.Actions...)
		}
		result := enemy.Update(ctx.Level, ctx.Player)
		actions = append(actions, result.Actions...)
	}

	// Update Player (handles state transitions, knockback, and physics)
	playerResult := ctx.Player.Update(ctx.Level, ctx.Player)
	actions = append(actions, playerResult.Actions...)

	return actions
}

func (g *MainGameState) Draw(screen *ebiten.Image, ctx *GameContext) {
	cameraMatrix := ctx.Camera.WorldToScreen()
	viewRect := ctx.Camera.GetViewRect()

	for _, row := range ctx.Level.Tiles {
		for _, tile := range row {
			if tile.GetBounds().Intersects(viewRect) {
				tile.Draw(screen, cameraMatrix)
				if tile.solid {
					tile.DrawDebugInfo(screen, cameraMatrix)
				}
			}
		}
	}

	// TODO: update drawing of objects/enemies/player so they are drawn based on sorted y coordinate
	for _, object := range ctx.Level.Objects {
		if object.GetBounds().Intersects(viewRect) {
			object.Draw(screen, cameraMatrix)
			object.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	for _, enemy := range ctx.Level.Enemies {
		if enemy.GetBounds().Intersects(viewRect) {
			enemy.Draw(screen, cameraMatrix)
			enemy.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	ctx.Player.Draw(screen, cameraMatrix)
	ctx.Player.DrawDebugInfo(screen, cameraMatrix)

	if ShowDebugInfo {
		for _, ds := range ctx.DamageSources {
			ds.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	DrawHeadsUpDisplay(screen, ctx.Player)
}
