package main

import (
	"fmt"
	"log"
	"math/rand/v2"

	"roguerpg/core"
	"roguerpg/enemy"
	"roguerpg/level"
	"roguerpg/objects"
	"roguerpg/player"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 384
	ScreenHeight = 240

	ShowDebugInfo = true
)

type Game struct {
	core.GameContext
	StateStack []core.GameState
}

// AddObjectsToLevel populates the level with objects (stairs, chests).
func AddObjectsToLevel(lvl *level.Level, isFirstLevel, isFinalLevel bool) {
	if !isFirstLevel {
		pos := lvl.FindRandomFloorLocation()
		lvl.Objects = append(lvl.Objects, objects.NewStairs(pos, true))
	}

	if !isFinalLevel {
		pos := lvl.FindRandomFloorLocation()
		lvl.Objects = append(lvl.Objects, objects.NewStairs(pos, false))
	}

	// Add Chests
	numChests := 10
	for range numChests {
		upgradeType := core.UpgradeType(rand.IntN(7))
		pos := lvl.FindRandomFloorLocation()
		lvl.Objects = append(lvl.Objects, objects.NewChest(pos, upgradeType))
	}
}

// AddEnemiesToLevel populates the level with enemies.
func AddEnemiesToLevel(lvl *level.Level, difficulty int) {
	numEnemies := 15 + difficulty
	for range numEnemies {
		pos := lvl.FindRandomFloorLocation()
		// enemyType := rand.IntN(5)
		var newEnemy core.Character
		newEnemy = enemy.NewLichEnemy(pos)

		// switch enemyType {
		// case 0:
		// 	newEnemy = enemy.NewBatEnemy(pos)
		// case 1:
		// 	newEnemy = enemy.NewBlobEnemy(pos)
		// case 2:
		// 	newEnemy = enemy.NewGoblinEnemy(pos)
		// case 3:
		// 	newEnemy = enemy.NewGhostEnemy(pos)
		// case 4:
		// 	newEnemy = enemy.NewSpikeTurtleEnemy(pos)
		// default:
		// 	newEnemy = enemy.NewBatEnemy(pos)
		// }
		lvl.Enemies = append(lvl.Enemies, newEnemy)
	}
}

func GetUpstairsLocation(lvl *level.Level) (core.Location, error) {
	for _, obj := range lvl.Objects {
		if stairs, ok := obj.(*objects.Stairs); ok {
			if stairs.IsUpstairs {
				return stairs.Location(), nil
			}
		}
	}
	return core.Location{}, fmt.Errorf("upstairs not found")
}

func GetDownstairsLocation(lvl *level.Level) (core.Location, error) {
	for _, obj := range lvl.Objects {
		if stairs, ok := obj.(*objects.Stairs); ok {
			if !stairs.IsUpstairs {
				return stairs.Location(), nil
			}
		}
	}
	return core.Location{}, fmt.Errorf("downstairs not found")
}

func buildLevels() []*level.Level {
	levels := []*level.Level{}
	numLevels := 5
	for i := range numLevels {
		lvl := level.BuildLevel(70, 50)
		isFirstLevel := (i == 0)
		isFinalLevel := (i == numLevels-1)
		AddObjectsToLevel(lvl, isFirstLevel, isFinalLevel)
		AddEnemiesToLevel(lvl, i)

		levels = append(levels, lvl)
	}

	// Link levels
	for i := range numLevels {
		if i > 0 {
			levels[i].UpLevel = levels[i-1]
		}
		if i < numLevels-1 {
			levels[i].DownLevel = levels[i+1]
		}
	}
	return levels
}

func NewGame() *Game {
	levels := buildLevels()
	p := player.NewPlayer()
	p.SetLocation(levels[0].FindRandomFloorLocation())

	return &Game{
		GameContext: core.GameContext{
			Level:         levels[0],
			Player:        p,
			Camera:        NewCamera(ScreenWidth, ScreenHeight), // Camera is in main (camera.go is main package)
			DamageSources: []*core.DamageSource{},
		},
		StateStack: []core.GameState{&MainGameState{}},
	}
}

func (g *Game) Update() error {
	ctx := &g.GameContext

	// Update only the top-most state on the stack.
	activeState := g.StateStack[len(g.StateStack)-1]
	actions := activeState.Update(ctx)

	// Execute Actions (handles push/pop state, damage source creation)
	g.executeActions(actions)

	// Global Game-World Updates
	if mainState, ok := activeState.(*MainGameState); ok {
		// Access level, player, camera via the embedded Context field
		for _, ds := range g.DamageSources {
			mainState.handleDamageSource(ctx, ds)
		}

		// Use type assertions/helpers to update Level enemies/objects if context is holding core.Level
		if lvl, ok := g.Level.(*level.Level); ok {
			lvl.Enemies = g.cleanupDeadEnemies(lvl.Enemies)
			lvl.Objects = g.cleanupObjects(lvl.Objects)
		}

		g.Camera.CenterOn(g.Player.Location())
	}

	PlayMusic()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Iterate from the bottom up to draw all states/windows.
	for _, state := range g.StateStack {
		state.Draw(screen, &g.GameContext)
	}
}

func (g *Game) GoDownLevel() {
	lvl, ok := g.Level.(*level.Level)
	if !ok {
		return
	}

	if lvl.DownLevel != nil {
		g.Level = lvl.DownLevel
		newLvl := lvl.DownLevel
		newLoc, err := GetUpstairsLocation(newLvl)
		if err != nil {
			log.Printf("Error finding upstairs location: %v", err)
			// Fallback
			newLoc = newLvl.FindRandomFloorLocation()
		}
		g.Player.SetLocation(newLoc)
	}
}

func (g *Game) GoUpLevel() {
	lvl, ok := g.Level.(*level.Level)
	if !ok {
		return
	}

	if lvl.UpLevel != nil {
		g.Level = lvl.UpLevel
		newLvl := lvl.UpLevel
		newLoc, err := GetDownstairsLocation(newLvl)
		if err != nil {
			log.Printf("Error finding downstairs location: %v", err)
			// Fallback
			newLoc = newLvl.FindRandomFloorLocation()
		}
		g.Player.SetLocation(newLoc)
	}
}

func (g *Game) executeActions(actions []core.Action) {
	g.DamageSources = []*core.DamageSource{}

	// Check if level is concrete type to append objects
	lvl, ok := g.Level.(*level.Level)
	if !ok {
		return
	}

	for _, action := range actions {
		switch action.Type {
		case core.ActionPushState:
			g.StateStack = append(g.StateStack, action.GameState)
		case core.ActionPopState:
			g.StateStack = g.StateStack[:len(g.StateStack)-1]
		case core.ActionCreateDamageSource:
			g.DamageSources = append(g.DamageSources, action.DamageSource)
		case core.ActionDropBomb:
			newBomb := objects.NewBomb(action.Location, action.Level)
			lvl.Objects = append(lvl.Objects, newBomb)
		case core.ActionShootArrow:
			newArrow := objects.NewArrow(action.Location, action.Direction, action.Level)
			lvl.Objects = append(lvl.Objects, newArrow)
		case core.ActionThrowBoomerang:
			newBoomerang := objects.NewBoomerang(action.Location, action.Direction, action.Level)
			lvl.Objects = append(lvl.Objects, newBoomerang)
		case core.ActionCreateStar:
			newStars := objects.NewStarSet(action.Location, action.Direction, action.Target, action.Level)
			for _, s := range newStars {
				lvl.Objects = append(lvl.Objects, s)
			}
		case core.ActionReturnBoomerang:
			g.Player.ReturnBoomerang()
		case core.ActionExplosion:
			newBombExplosion := objects.NewBombExplosion(action.Location, action.Level)
			lvl.Objects = append(lvl.Objects, newBombExplosion)
		case core.ActionCreateFire:
			newFire := objects.NewFire(action.Location, action.Direction)
			lvl.Objects = append(lvl.Objects, newFire)
		case core.ActionSwitchWeapon:
			g.Player.SwitchWeapon(action.WeaponType)
		case core.ActionGoUpLevel:
			g.GoUpLevel()
		case core.ActionGoDownLevel:
			g.GoDownLevel()
		case core.ActionGainXP:
			g.Player.AddExperience(action.Experience)
		case core.ActionDropHeart:
			newHeart := objects.NewHeart(action.Location)
			lvl.Objects = append(lvl.Objects, newHeart)
		case core.ActionShowChestItem:
			newChestItem := objects.NewChestItem(action.UpgradeType, action.Location)
			g.Player.AddUpgrade(action.UpgradeType)
			lvl.Objects = append(lvl.Objects, newChestItem)
		default:
		}
	}
}

func (g *Game) cleanupDeadEnemies(enemies []core.Character) []core.Character {
	liveEnemies := enemies[:0]
	for _, enemy := range enemies {
		if !enemy.IsDead() {
			liveEnemies = append(liveEnemies, enemy)
		}
	}
	return liveEnemies
}

func (g *Game) cleanupObjects(objs []core.GameObject) []core.GameObject {
	liveObjects := objs[:0]
	for _, object := range objs {
		if !object.CanRemove() {
			liveObjects = append(liveObjects, object)
		}
	}
	return liveObjects
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
