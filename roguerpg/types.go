package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Location struct {
	X float64
	Y float64
}

type Vector struct {
	X float64
	Y float64
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

func (r1 Rect) Intersects(r2 Rect) bool {
	return r1.Left < r2.Right && r1.Right > r2.Left &&
		r1.Top < r2.Bottom && r1.Bottom > r2.Top
}

type CollisionAxis int

const (
	AxisX CollisionAxis = iota
	AxisY
)

// EntityTag is used to categorize game objects for collision filtering (e.g., friendly fire)
type EntityTag int

const (
	TagPlayer EntityTag = iota
	TagEnemy
	TagTile
)

// DamageSourceConfig holds the Rect offset and damage value for a specific attack frame.
// Rect is relative to the player's center (Location).
type DamageSourceConfig struct {
	HitBox Rect
	Damage int
}

// DamageSource represents an active attack hitbox in the world.
type DamageSource struct {
	SourceTag  EntityTag // e.g., TagPlayer, TagEnemy
	HitBox     Rect      // The current world-space hitbox of the attack
	Damage     int
	debugImage *ebiten.Image
}

func NewDamageSource(sourceTag EntityTag, hitBox Rect, damage int) *DamageSource {
	return &DamageSource{
		SourceTag:  sourceTag,
		HitBox:     hitBox,
		Damage:     damage,
		debugImage: createDebugRectImage(hitBox),
	}
}

func (ds *DamageSource) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}

	if ds.debugImage == nil || dotImage == nil {
		return
	}

	// Draw the Hitbox rectangle
	hb := ds.HitBox

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(hb.Left, hb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(ds.debugImage, opRect)
}

// GameObject is an interface for any entity in the game world.
type GameObject interface {
	GetBounds() Rect // General bounding box for drawing
	Update()
	DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM)
}

// PhysicalObject is anything that participates in collisions and pushing.
type PhysicalObject interface {
	GameObject
	GetPushBox() Rect
}

// Character is a specialized entity that can take damage and be knocked back.
type Character interface {
	GameObject
	GetHurtBox() Rect
	TakeDamage(damage int)
	ApplyKnockback(force Vector, duration int)
	IsKnockedBack() bool
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
