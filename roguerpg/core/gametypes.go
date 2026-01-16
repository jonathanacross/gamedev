package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// EntityTag is used to categorize game objects for collision filtering (e.g., friendly fire)
type EntityTag int

const (
	TagPlayer EntityTag = iota
	TagEnemy
	TagTile
)

const (
	TileSize = 16

	KnockbackForce    = 3.0
	KnockbackDuration = 6
)

type DamageType int

const (
	DamageTypePhysical DamageType = iota
	DamageTypeStun
	DamageTypeImpact
	DamageTypeMagic
)

// DamageSourceConfig holds the Rect offset and damage value for a specific attack frame.
// Rect is relative to the player's center (Location).
type DamageSourceConfig struct {
	HitBox Rect
	Damage int
}

// DamageSource is a generic struct for anything that deals damage or applies effects (like stun).
type DamageSource struct {
	SourceTag     EntityTag // e.g., TagPlayer, TagEnemy
	HitBox        Rect      // The current world-space hitbox of the attack
	Damage        int       // Used for DamageTypePhysical
	Type          DamageType
	Duration      int                           // Used for DamageTypeStun (in frames)
	OnHit         func()                        // Callback when damage source hits a target
	IsReflectable bool                          // Can be reflected by a reflector
	IsReflector   bool                          // Can reflect reflectable sources
	Direction     Direction                     // Direction of the source (needed for reflection axis)
	OnReflect     func(reflector *DamageSource) // Callback when reflected
}

func NewDamageSource(sourceTag EntityTag, hitBox Rect, damageType DamageType, damage int) *DamageSource {
	return &DamageSource{
		SourceTag: sourceTag,
		HitBox:    hitBox,
		Damage:    damage,
		Type:      damageType,
	}
}

// Specialized version of NewDamageSource for stun damage.
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
	ActionCreateFire
	ActionThrowBoulder
	ActionSwitchWeapon
	ActionGoUpLevel
	ActionGoDownLevel
	ActionGainXP
	ActionDropHeart
	ActionShowChestItem
	ActionDoUpgrade
)

type WeaponType int

const (
	WeaponSword WeaponType = iota
	WeaponBomb
	WeaponBoomerang
	WeaponShield
	WeaponBow
	WeaponWand
	WeaponNone
)

type UpgradeType int

const (
	UpgradeTypeSword UpgradeType = iota
	UpgradeTypeBomb
	UpgradeTypeBoomerang
	UpgradeTypeShield
	UpgradeTypeBow
	UpgradeTypeWand
	UpgradeTypeHeart
	UpgradeTypeNone
)

type PlayerUpgrades map[UpgradeType]int

// Action is the generic struct returned by any GameObject
// to signal an intent to change the game state.
// Various other fields are populated depending on the ActionType.
type Action struct {
	Type         ActionType
	Location     Location
	Direction    Vector
	GameState    GameState     // used for ActionPushState
	DamageSource *DamageSource // used for ActionCreateDamageSource
	WeaponType   WeaponType    // used for ActionSwitchWeapon
	Experience   int           // used for ActionGainXP
	Target       Character     // used for ActionCreateStar
	UpgradeType  UpgradeType   // used for ActionShowChestItem
	Level        int           // used for weapon attacks
}

type UpdateResult struct {
	Actions []Action
}

func (ds *DamageSource) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	hitBoxColor := color.RGBA{R: 255, G: 255, B: 0, A: 255}
	DrawDebugRect(screen, ds.HitBox, hitBoxColor, cameraMatrix)
}
