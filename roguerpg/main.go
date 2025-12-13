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

type GameContext struct {
	Level         *Level
	Player        *Player
	Camera        *Camera
	DamageSources []*DamageSource // current damage sources; saved here for debug drawing
}

type Game struct {
	GameContext
	StateStack []GameState
}

func buildLevels() []*Level {
	levels := []*Level{}
	numLevels := 5
	for i := range numLevels {
		level := BuildLevel(70, 50)
		isFirstLevel := (i == 0)
		isFinalLevel := (i == numLevels-1)
		level.AddObjects(isFirstLevel, isFinalLevel)
		level.AddEnemies(i)

		levels = append(levels, level)
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
	player := NewPlayer()
	player.SetLocation(levels[0].FindRandomFloorLocation())

	return &Game{
		GameContext: GameContext{
			Level:         levels[0],
			Player:        player,
			Camera:        NewCamera(ScreenWidth, ScreenHeight),
			DamageSources: []*DamageSource{},
		},
		StateStack: []GameState{&MainGameState{}},
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
			// Note: The handleDamageSource function itself would need access to g.Context.Level/g.Context.Player,
			// or you would pass the context to it as well.
			mainState.handleDamageSource(ctx, ds)
		}

		g.Level.Enemies = g.cleanupDeadEnemies(g.Level.Enemies)
		g.Level.Objects = g.cleanupObjects(g.Level.Objects)
		g.Camera.CenterOn(g.Player.Location())
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Iterate from the bottom up to draw all states/windows.
	for _, state := range g.StateStack {
		state.Draw(screen, &g.GameContext)
	}
}

func (g *Game) GoDownLevel() {
	if g.Level.DownLevel != nil {
		g.Level = g.Level.DownLevel
		newLoc, err := g.Level.GetUpstairsLocation()
		if err != nil {
			log.Printf("Error finding upstairs location: %v", err)
			return
		}
		g.Player.SetLocation(newLoc)
	}
}

func (g *Game) GoUpLevel() {
	if g.Level.UpLevel != nil {
		g.Level = g.Level.UpLevel
		newLoc, err := g.Level.GetDownstairsLocation()
		if err != nil {
			log.Printf("Error finding downstairs location: %v", err)
			return
		}
		g.Player.SetLocation(newLoc)
	}
}

func (g *Game) executeActions(actions []Action) {
	g.DamageSources = []*DamageSource{}
	for _, action := range actions {
		switch action.Type {
		case ActionPushState:
			g.StateStack = append(g.StateStack, action.GameState)
		case ActionPopState:
			g.StateStack = g.StateStack[:len(g.StateStack)-1]
		case ActionCreateDamageSource:
			g.DamageSources = append(g.DamageSources, action.DamageSource)
		case ActionDropBomb:
			newBomb := NewBomb(action.Location)
			g.Level.Objects = append(g.Level.Objects, newBomb)
		case ActionShootArrow:
			newArrow := NewArrow(action.Location, action.Direction)
			g.Level.Objects = append(g.Level.Objects, newArrow)
		case ActionThrowBoomerang:
			newBoomerang := NewBoomerang(action.Location, action.Direction, rand.IntN(3)+1)
			g.Level.Objects = append(g.Level.Objects, newBoomerang)
		case ActionCreateStar:
			newStar := NewStar(action.Location, action.Direction)
			g.Level.Objects = append(g.Level.Objects, newStar)
		case ActionReturnBoomerang:
			g.Player.ReturnBoomerang()
		case ActionExplosion:
			NewBombExplosion := NewBombExplosion(action.Location)
			g.Level.Objects = append(g.Level.Objects, NewBombExplosion)
		case ActionSwitchWeapon:
			g.Player.SwitchWeapon(action.WeaponType)
		case ActionGoUpLevel:
			g.GoUpLevel()
		case ActionGoDownLevel:
			g.GoDownLevel()
		case ActionGainXP:
			g.Player.Experience += action.Experience
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
