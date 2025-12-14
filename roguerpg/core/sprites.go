package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// BaseSprite provides common fields and methods for any visible game entity.
type BaseSprite struct {
	Loc        Location // Renamed from Location to avoid method conflict
	Image      *ebiten.Image
	SrcRect    image.Rectangle
	DrawOffset Location
}

// GetBounds returns the drawing rectangle for the BaseSprite.
func (bs *BaseSprite) GetBounds() Rect {
	x := bs.Loc.X - bs.DrawOffset.X
	y := bs.Loc.Y - bs.DrawOffset.Y

	width := float64(bs.SrcRect.Dx())
	height := float64(bs.SrcRect.Dy())

	return Rect{
		Left:   x,
		Top:    y,
		Right:  x + width,
		Bottom: y + height,
	}
}

func (bs *BaseSprite) GetX() float64 { return bs.Loc.X }

func (bs *BaseSprite) GetY() float64 { return bs.Loc.Y }

func (bs *BaseSprite) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}
	// Using global variables from same package (debugutils.go)
	if DotImage == nil {
		return
	}

	// Draw the bounds rectangle
	hb := bs.GetBounds()
	debugImage := GetDebugRectImage(hb)
	if debugImage == nil {
		return
	}

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(hb.Left, hb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(debugImage, opRect)

	// Draw the Location Dot
	opDot := &ebiten.DrawImageOptions{}
	opDot.GeoM.Translate(bs.Loc.X-1.5, bs.Loc.Y-1.5) // Assuming dotSize approx 3-4, centering.
	opDot.GeoM.Concat(cameraMatrix)
	screen.DrawImage(DotImage, opDot)
}

func (bs *BaseSprite) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if bs.Image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.Loc.X-bs.DrawOffset.X, bs.Loc.Y-bs.DrawOffset.Y)
	op.GeoM.Concat(cameraMatrix)
	currImage := bs.Image.SubImage(bs.SrcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

// Base entity for physical objects.  Tiles, monsters, the player, items,
// will be based on this struct.
type BasePhysical struct {
	BaseSprite

	PushBoxOffset Rect // The offset of the physical box relative to Location.
}

func (bp *BasePhysical) Location() Location {
	return bp.Loc
}

func (bp *BasePhysical) SetLocation(l Location) {
	bp.Loc = l
}

// GetPushBox implements the PhysicalObject interface.
func (bp *BasePhysical) GetPushBox() Rect {
	return bp.PushBoxOffset.Offset(bp.Loc.X, bp.Loc.Y)
}

// DrawDebugInfo overrides the BaseSprite version to draw the PushBox.
func (bp *BasePhysical) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}
	if DotImage == nil {
		return
	}

	// Draw the PushBox rectangle
	pb := bp.GetPushBox()
	debugImagePB := GetDebugRectImage(pb)
	if debugImagePB == nil {
		return
	}

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(pb.Left, pb.Top)

	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(debugImagePB, opRect)

	// Draw the Location Dot
	opDot := &ebiten.DrawImageOptions{}
	opDot.GeoM.Translate(bp.Loc.X-1.5, bp.Loc.Y-1.5)
	opDot.GeoM.Concat(cameraMatrix)
	screen.DrawImage(DotImage, opDot)
}

// Tile
type Tile struct {
	BasePhysical
	Solid bool
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
				Loc:     location,
				Image:   image,
				SrcRect: srcRect,
				DrawOffset: Location{
					X: 0,
					Y: 0,
				},
			},
			PushBoxOffset: pushBox,
		},
		Solid: solid,
	}
}

func (t *Tile) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	t.BasePhysical.DrawDebugInfo(screen, cameraMatrix)
}
