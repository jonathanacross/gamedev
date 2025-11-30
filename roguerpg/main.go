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

// TODO: make these fields capital
type GameContext struct {
	level         *Level
	player        *Player
	camera        *Camera
	damageSources []*DamageSource // current damage sources; saved here for debug drawing
}

type Game struct {
	GameContext
	StateStack []GameState
}

func NewGame() *Game {
	level := BuildLevel(70, 50)
	level.AddObjects()
	level.AddEnemies()

	player := NewPlayer()
	player.SetLocation(level.FindRandomFloorLocation())

	return &Game{
		GameContext: GameContext{
			level:         level,
			player:        player,
			camera:        NewCamera(ScreenWidth, ScreenHeight),
			damageSources: []*DamageSource{},
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
		for _, ds := range g.damageSources {
			// Note: The handleDamageSource function itself would need access to g.Context.Level/g.Context.Player,
			// or you would pass the context to it as well.
			mainState.handleDamageSource(ctx, ds)
		}

		g.level.Enemies = g.cleanupDeadEnemies(g.level.Enemies)
		g.level.Objects = g.cleanupObjects(g.level.Objects)
		g.camera.CenterOn(g.player.Location())
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Iterate from the bottom up to draw all states/windows.
	for _, state := range g.StateStack {
		state.Draw(screen, &g.GameContext)
	}
}

func (g *Game) executeActions(actions []Action) {
	g.damageSources = []*DamageSource{}
	for _, action := range actions {
		switch action.Type {
		case ActionPushState:
			g.StateStack = append(g.StateStack, action.GameState)
		case ActionPopState:
			g.StateStack = g.StateStack[:len(g.StateStack)-1]
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
