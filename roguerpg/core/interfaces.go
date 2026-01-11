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
	// Touch is called by MainGameState when the player overlaps the object's PushBox.
	Touch(level Level, player Player) []Action
	// Interact is called by MainGameState when the player presses the
	// interaction key while overlapping the object's PushBox.
	Interact(level Level, player Player) []Action
}

// Character is a specialized entity that can take damage and be knocked back.
type Character interface {
	PhysicalObject
	GetHurtBox() Rect
	HandleHit(ds *DamageSource)
	// TODO: see if I can remove TakeDamage, ApplyKnockback, and IsKnockedBack, ApplyStun, and IsStunned.
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
	AddHealth(amount int)
	GetExperience() int
	GetPrimaryWeapon() WeaponType
	GetSecondaryWeapon() WeaponType
	GetWeaponProgress(weapon WeaponType) float64
	AddUpgrade(upgradeType UpgradeType)
	DoUpgrade(upgradeType UpgradeType)
	GetCurrentStats() PlayerUpgrades
	GetFutureUpgrades() PlayerUpgrades
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

type RenderOrder int

const (
	RenderOrderBehind RenderOrder = iota
	RenderOrderFront
)

type PlayerContext interface {
	Location() Location
	GetDirection() Direction
	AddAction(action Action)
	CreateDamageSource(ds *DamageSource)
	SetBlocking(active bool)
}

type Weapon interface {
	// Update handles weapon logic (cooldowns, etc.) and returns actions if needed
	Update(ctx PlayerContext, level Level)

	// OnAttack is called when the player triggers the weapon
	OnAttack(ctx PlayerContext)

	// GetRenderOrder returns whether the weapon should be drawn behind or in front of the player
	GetRenderOrder(playerDirection Direction) RenderOrder

	// Draw renders the weapon relative to the player
	Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM, playerPos Location, playerDir Direction, playerFrame int)

	// GetCooldownProgress returns 0.0 to 1.0 indicating how ready the weapon is (1.0 = ready)
	GetCooldownProgress() float64

	// SetLevel sets the weapon level
	SetLevel(level int)

	// IsAttacking returns true if the weapon is currently in an attack animation
	IsAttacking() bool
}
