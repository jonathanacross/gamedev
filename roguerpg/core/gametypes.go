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

var ShowDebugInfo bool = false

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
