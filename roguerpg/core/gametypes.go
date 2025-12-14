package core

import "github.com/hajimehoshi/ebiten/v2"

// EntityTag is used to categorize game objects for collision filtering (e.g., friendly fire)
type EntityTag int

const (
	TagPlayer EntityTag = iota
	TagEnemy
	TagTile
)

const TileSize = 16

type DamageType int

const (
	DamageTypePhysical DamageType = iota
	DamageTypeStun
)

// DamageSourceConfig holds the Rect offset and damage value for a specific attack frame.
// Rect is relative to the player's center (Location).
type DamageSourceConfig struct {
	HitBox Rect
	Damage int
}

// DamageSource is a generic struct for anything that deals damage or applies effects (like stun).
type DamageSource struct {
	SourceTag EntityTag  // e.g., TagPlayer, TagEnemy
	HitBox    Rect       // The current world-space hitbox of the attack
	Damage    int        // Used for DamageTypePhysical
	Type      DamageType // Physical or Stun
	Duration  int        // Used for DamageTypeStun (in frames)
	OnHit     func()     // Callback when damage source hits a target
}

var ShowDebugInfo bool = false // Global flag, or move to a config? For now assuming global.
// Actually ShowDebugInfo was in main.go constants usually or types.go?
// It was in main.go as const ShowDebugInfo = false.
// But DamageSource.DrawDebugInfo uses it.
// If it was a const in main, it's not accessible here unless passed or duplicated.
// Let's assume we can remove the check or add a global variable in core.
// Better: Pass it in DrawDebugInfo or just always draw if called?
// The DrawDebugInfo method usually checks the flag.
// I'll add a SetShowDebugInfo in core or just export a variable.
// For now, I'll add it as a variable.

func NewDamageSource(sourceTag EntityTag, hitBox Rect, damage int) *DamageSource {
	return &DamageSource{
		SourceTag: sourceTag,
		HitBox:    hitBox,
		Damage:    damage,
		Type:      DamageTypePhysical,
	}
}

func NewStunSource(sourceTag EntityTag, hitBox Rect, duration int) *DamageSource {
	return &DamageSource{
		SourceTag: sourceTag,
		HitBox:    hitBox,
		Duration:  duration,
		Type:      DamageTypeStun,
	}
}

// ActionType defines the different kinds of actions that can be requested.
type ActionType int

const (
	ActionPushState ActionType = iota
	ActionPopState
	ActionCreateDamageSource
	ActionDropBomb
	ActionShootArrow
	ActionCreateStar
	ActionThrowBoomerang
	ActionReturnBoomerang
	ActionExplosion
	ActionSwitchWeapon
	ActionGoUpLevel
	ActionGoDownLevel
	ActionGainXP
)

// Helper for WeaponType to avoid circular dep if WeaponType is in projectile...
// Wait, WeaponType was in types.go? No, I don't see it in the `types.go` dump I read.
// I must have missed it or it was in `player.go`?
// Let me check `player.go` for `WeaponType`.
// If it's used in `Action`, it needs to be in `core`.

type WeaponType int

const (
	WeaponSword WeaponType = iota
	WeaponBomb
	WeaponBoomerang
	WeaponShield
	WeaponBow
	WeaponWand
)

// Action is the generic struct returned by any GameObject
// to signal an intent to change the game state.
// Various other fields are populated depending on the ActionType.
type Action struct {
	Type         ActionType
	Location     Location
	Direction    Vector
	GameState    GameState     // only used for ActionPushState
	DamageSource *DamageSource // only populated for ActionCreateDamageSource
	WeaponType   WeaponType    // only used for ActionSwitchWeapon
	Experience   int           // only used for ActionGainXP
	Target       Character     // only used for ActionCreateStar
}

type UpdateResult struct {
	Actions []Action
}

func (ds *DamageSource) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}

	// Draw the Hitbox rectangle
	hb := ds.HitBox
	debugImage := GetDebugRectImage(hb)
	if debugImage == nil {
		return
	}

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(hb.Left, hb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(debugImage, opRect)
}
