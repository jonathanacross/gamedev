package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// BaseSprite provides common fields and methods for any visible game entity.
// It handles drawing a single sprite or the current frame of an animation.
// This implements the GameObject interface.
type BaseSprite struct {
	Location
	image      *ebiten.Image
	srcRect    image.Rectangle
	drawOffset Location
}

// GetBounds returns the drawing rectangle for the BaseSprite.
func (bs *BaseSprite) GetBounds() Rect {
	x := bs.X - bs.drawOffset.X
	y := bs.Y - bs.drawOffset.Y

	width := float64(bs.srcRect.Dx())
	height := float64(bs.srcRect.Dy())

	return Rect{
		Left:   x,
		Top:    y,
		Right:  x + width,
		Bottom: y + height,
	}
}

func (bs *BaseSprite) GetX() float64 { return bs.X }

func (bs *BaseSprite) GetY() float64 { return bs.Y }

func (bs *BaseSprite) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}

	// Draw the pushbox rectangle
	hb := bs.GetBounds()
	debugImage := GetDebugRectImage(hb)

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(hb.Left, hb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(debugImage, opRect)

	// Draw the Location Dot
	opDot := &ebiten.DrawImageOptions{}
	opDot.GeoM.Translate(bs.X-dotSize/2, bs.Y-dotSize/2)
	opDot.GeoM.Concat(cameraMatrix)
	screen.DrawImage(dotImage, opDot)
}

func (bs *BaseSprite) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X-bs.drawOffset.X, bs.Y-bs.drawOffset.Y)
	op.GeoM.Concat(cameraMatrix)
	currImage := bs.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

// Base entity for physical objects.  Tiles, monsters, the player, items,
// will be based on this struct.
type BasePhysical struct {
	BaseSprite

	pushBoxOffset Rect // The offset of the physical box relative to Location.
}

func (bp *BasePhysical) Location() Location {
	return Location{X: bp.X, Y: bp.Y}
}

func (bp *BasePhysical) SetLocation(l Location) {
	bp.X = l.X
	bp.Y = l.Y
}

// GetPushBox implements the PhysicalObject interface.
func (bp *BasePhysical) GetPushBox() Rect {
	return bp.pushBoxOffset.Offset(bp.X, bp.Y)
}

// DrawDebugInfo overrides the BaseSprite version to draw the PushBox.
func (bp *BasePhysical) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	// Draw base debug info (Location Dot)
	bp.BaseSprite.DrawDebugInfo(screen, cameraMatrix)

	if !ShowDebugInfo {
		return
	}

	// Draw the PushBox rectangle
	pb := bp.GetPushBox()
	debugImage := GetDebugRectImage(pb)

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(pb.Left, pb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(debugImage, opRect)
}

type Tile struct {
	BasePhysical
	solid bool
}

func NewTile(location Location, image *ebiten.Image, srcRect image.Rectangle, solid bool) *Tile {
	pushBox := Rect{
		Left:   0,
		Top:    0,
		Right:  float64(srcRect.Dx()),
		Bottom: float64(srcRect.Dy()),
	}
	return &Tile{
		BasePhysical: BasePhysical{
			BaseSprite: BaseSprite{
				Location: location,
				image:    image,
				srcRect:  srcRect,
				drawOffset: Location{
					X: 0,
					Y: 0,
				},
			},
			pushBoxOffset: pushBox,
		},
		solid: solid,
	}
}

func (t *Tile) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	t.BasePhysical.DrawDebugInfo(screen, cameraMatrix)
}
