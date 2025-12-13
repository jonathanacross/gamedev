package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Location Vector

type Vector struct {
	X float64
	Y float64
}

func ZeroVector() Vector {
	return Vector{X: 0, Y: 0}
}

func (v Vector) Minus(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector) Plus(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v Vector) Normalize() Vector {
	length := v.Length()
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
	}
}

func (v Vector) Scale(scalar float64) Vector {
	return Vector{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v Vector) Rotate(angleRadians float64) Vector {
	s := math.Sin(angleRadians)
	c := math.Cos(angleRadians)
	return Vector{
		X: v.X*c + v.Y*s,
		Y: -v.X*s + v.Y*c,
	}
}

type Rect struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (r Rect) Width() float64 {
	return r.Right - r.Left
}

func (r Rect) Height() float64 {
	return r.Bottom - r.Top
}

func (r Rect) Offset(x, y float64) Rect {
	return Rect{
		Left:   r.Left + x,
		Top:    r.Top + y,
		Right:  r.Right + x,
		Bottom: r.Bottom + y,
	}
}

func (r1 Rect) IntersectsX(r2 Rect) bool {
	return r1.Left < r2.Right && r1.Right > r2.Left
}

func (r1 Rect) IntersectsY(r2 Rect) bool {
	return r1.Top < r2.Bottom && r1.Bottom > r2.Top
}

func (r1 Rect) Intersects(r2 Rect) bool {
	return r1.IntersectsX(r2) && r1.IntersectsY(r2)
}

type CollisionAxis int

const (
	AxisX CollisionAxis = iota
	AxisY
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

func VectorToDirection(dirVector Vector) Direction {
	angle := math.Atan2(dirVector.Y, dirVector.X)
	if angle >= -math.Pi/4 && angle < math.Pi/4 {
		return Right
	} else if angle >= math.Pi/4 && angle < 3*math.Pi/4 {
		return Down
	} else if angle >= -3*math.Pi/4 && angle < -math.Pi/4 {
		return Up
	} else {
		return Left
	}
}

func DirectionToVector(dir Direction) Vector {
	switch dir {
	case Left:
		return Vector{X: -1, Y: 0}
	case Right:
		return Vector{X: 1, Y: 0}
	case Up:
		return Vector{X: 0, Y: -1}
	case Down:
		return Vector{X: 0, Y: 1}
	}
	return Vector{X: 0, Y: 0}
}

// EntityTag is used to categorize game objects for collision filtering (e.g., friendly fire)
type EntityTag int

const (
	TagPlayer EntityTag = iota
	TagEnemy
	TagTile
)

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
}

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
}

type UpdateResult struct {
	Actions []Action
}

func (ds *DamageSource) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}

	debugImage := GetDebugRectImage(ds.HitBox)

	// Draw the Hitbox rectangle
	hb := ds.HitBox

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(hb.Left, hb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(debugImage, opRect)
}

// GameObject is an interface for any entity in the game world.
type GameObject interface {
	GetBounds() Rect // General bounding box for drawing
	Update(level *Level, player *Player) UpdateResult
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
	Interact(level *Level, player *Player) []Action
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

	IsDead() bool
}

// CalculateKnockbackForce computes a normalized, scaled vector pointing from the attacker to the defender.
func CalculateKnockbackForce(attackerLoc Location, defenderLoc Location, speed float64) Vector {
	direction := Vector{
		X: defenderLoc.X - attackerLoc.X,
		Y: defenderLoc.Y - attackerLoc.Y,
	}

	if direction.Length() == 0 {
		return Vector{X: 0, Y: 0}
	}
	return direction.Normalize().Scale(speed)
}

type GameState interface {
	// Update handles input and logic, using the global context.
	// It returns Actions that modify the main Game struct (e.g., PushState, PopState).
	Update(context *GameContext) []Action

	// Draw renders the state. We pass the camera matrix since some states (like MainGame)
	// need to translate world coordinates, while others ignore it.
	Draw(screen *ebiten.Image, context *GameContext)
}
