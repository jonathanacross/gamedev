package main

import (
	"roguerpg/core"
	"roguerpg/level"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// struct to handle the main game logic
type MainGameState struct {
}

// applyDamageToTargets applies damage to a list of characters (enemies or player)
func (g *MainGameState) applyDamageToTargets(damageSource *core.DamageSource, targets []core.Character) {
	for _, target := range targets {
		if target.IsDead() {
			continue
		}

		if damageSource.HitBox.Intersects(target.GetHurtBox()) {
			// Let the character decide how to react to this specific damage source
			target.HandleHit(damageSource)

			if damageSource.OnHit != nil {
				damageSource.OnHit()
			}
		}
	}
}

// handleDamageSource checks if a damage source hits the player or any enemy,
// depending on its SourceTag.
func (s *MainGameState) handleDamageSource(ctx *core.GameContext, damageSource *core.DamageSource) {
	if damageSource.SourceTag == core.TagPlayer {
		// Player attack hits enemies
		s.applyDamageToTargets(damageSource, ctx.Level.GetEnemies())
	} else if damageSource.SourceTag == core.TagEnemy {
		reflected := false
		// Check for reflection against player reflectors (e.g. Shield)
		if damageSource.IsReflectable {
			for _, otherDS := range ctx.DamageSources {
				if otherDS.SourceTag == core.TagPlayer && otherDS.IsReflector {
					if damageSource.HitBox.Intersects(otherDS.HitBox) {
						if damageSource.OnReflect != nil {
							damageSource.OnReflect(otherDS)
							reflected = true
							break
						}
					}
				}
			}
		}

		if !reflected {
			// Enemy attack hits the player
			s.applyDamageToTargets(damageSource, []core.Character{ctx.Player})
		}
	}
}

func (mg *MainGameState) handleInput(ctx *core.GameContext) []core.Action {
	var actions []core.Action
	player := ctx.Player
	level := ctx.Level
	isMoving := false

	// --- Handle State Transition Input (Global to MainGameState) ---

	// Open Weapon Selector Menu
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		action := core.Action{
			Type:      core.ActionPushState,
			GameState: WeaponSelectorInstance,
		}
		actions = append(actions, action)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		action := core.Action{
			Type:      core.ActionPushState,
			GameState: NewUpgradeSelector(player),
		}
		actions = append(actions, action)
	}

	if !player.IsActive() {
		return actions
	}

	// --- Handle Movement Input ---
	var moveVector core.Vector
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
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		if action := player.PrimaryAttack(); action != nil {
			actions = append(actions, *action)
		}
	}

	// Secondary Attack
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if action := player.SecondaryAttack(); action != nil {
			actions = append(actions, *action)
		}
	}

	// Interact with objects
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		playerPushBox := player.GetPushBox()

		// Check all objects for interaction
		for _, object := range ctx.Level.GetObjects() {
			if interactable, ok := object.(core.Interactable); ok {
				if playerPushBox.Intersects(interactable.GetPushBox()) {
					interactionActions := interactable.Interact(level, player)
					actions = append(actions, interactionActions...)
					// Only interact with one object at a time
					break
				}
			}
		}
	}

	return actions
}

func (mg *MainGameState) Update(ctx *core.GameContext) []core.Action {

	var actions []core.Action

	// Handle Input (New Step)
	inputActions := mg.handleInput(ctx)
	actions = append(actions, inputActions...)

	// Update Game Objects
	for _, object := range ctx.Level.GetObjects() {
		result := object.Update(ctx.Level, ctx.Player)
		actions = append(actions, result.Actions...)
	}

	// Update Enemies
	for _, enemy := range ctx.Level.GetEnemies() {
		result := enemy.Update(ctx.Level, ctx.Player)
		actions = append(actions, result.Actions...)
	}

	// Update Player (handles state transitions, knockback, and physics)
	playerResult := ctx.Player.Update(ctx.Level, ctx.Player)
	actions = append(actions, playerResult.Actions...)

	return actions
}

func (g *MainGameState) Draw(screen *ebiten.Image, ctx *core.GameContext) {
	cameraMatrix := ctx.Camera.WorldToScreen()
	viewRect := ctx.Camera.GetViewRect()

	// Cast core.Level to *level.Level to access tiles for drawing
	if lvl, ok := ctx.Level.(*level.Level); ok {
		for _, row := range lvl.Tiles {
			for _, tile := range row {
				if tile.GetBounds().Intersects(viewRect) {
					tile.Draw(screen, cameraMatrix)
					if tile.Solid && ShowDebugInfo {
						tile.DrawDebugInfo(screen, cameraMatrix)
					}
				}
			}
		}
	}

	// TODO: update drawing of objects/enemies/player so they are drawn based on sorted y coordinate
	for _, object := range ctx.Level.GetObjects() {
		if object.GetBounds().Intersects(viewRect) {
			object.Draw(screen, cameraMatrix)
			if ShowDebugInfo {
				object.DrawDebugInfo(screen, cameraMatrix)
			}
		}
	}

	for _, enemy := range ctx.Level.GetEnemies() {
		if enemy.GetBounds().Intersects(viewRect) {
			enemy.Draw(screen, cameraMatrix)
			if ShowDebugInfo {
				enemy.DrawDebugInfo(screen, cameraMatrix)
			}
		}
	}

	ctx.Player.Draw(screen, cameraMatrix)
	if ShowDebugInfo {
		ctx.Player.DrawDebugInfo(screen, cameraMatrix)
	}

	if ShowDebugInfo {
		for _, ds := range ctx.DamageSources {
			ds.DrawDebugInfo(screen, cameraMatrix)
		}
	}

	DrawHeadsUpDisplay(screen, ctx.Player)
}
