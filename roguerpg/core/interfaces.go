package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GameObject is an interface for any entity in the game world.
type GameObject interface {
	GetBounds() Rect // General bounding box for drawing
	Update(level Level, player Player) UpdateResult
	Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM)
	DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM)
	CanRemove() bool // indicate object can be removed from the game
}

// PhysicalObject is anything that participates in collisions and pushing.
type PhysicalObject interface {
	GameObject
	GetPushBox() Rect
	Location() Location
	SetLocation(loc Location)
}

// An object that the player can interact with.
type Interactable interface {
	PhysicalObject
	// Interact is called by MainGameState when the player presses the
	// interaction key while overlapping the object's PushBox.
	Interact(level Level, player Player) []Action
}

// Character is a specialized entity that can take damage and be knocked back.
type Character interface {
	PhysicalObject
	GetHurtBox() Rect
	TakeDamage(damage int)
	ApplyKnockback(force Vector, duration int)
	IsKnockedBack() bool

	ApplyStun(duration int)
	IsStunned() bool

	GetHealth() int
	GetMaxHealth() int
	IsDead() bool
}

type GameState interface {
	// Update handles input and logic, using the global context.
	Update(context *GameContext) []Action

	// Draw renders the state.
	Draw(screen *ebiten.Image, context *GameContext)
}

type GameContext struct {
	Level         Level
	Player        Player
	Camera        Camera
	DamageSources []*DamageSource
}

type Player interface {
	Character
	IsActive() bool
	PrimaryAttack() *Action
	SecondaryAttack() *Action
	StopMoving()
	Move(velocity Vector)
	SwitchWeapon(weapon WeaponType)
	ReturnBoomerang() // Used by ActionReturnBoomerang
	AddExperience(amount int)
	GetExperience() int
	GetPrimaryWeapon() WeaponType
	GetSecondaryWeapon() WeaponType
	GetWeaponProgress(weapon WeaponType) float64
}

type Level interface {
	GetTile(x, y int) *Tile

	IsTileSolid(tileX, tileY int) bool
	WorldToTile(loc Location) (int, int)
	TileToWorld(tileX, tileY int) Location
	FindNearestEnemy(loc Location, radius float64) Character

	GetEnemies() []Character
	GetObjects() []GameObject
}

type Camera interface {
	WorldToScreen() ebiten.GeoM
	GetViewRect() Rect
	CenterOn(loc Location)
}
