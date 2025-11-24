package main

import (
	"image"
	"image/color"
	"sync"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Size of the debug dot in pixels
const dotSize = 4.0

// Global ebiten.Image for the location marker dot (initialized in init()).
var dotImage *ebiten.Image

func init() {
	// Create the dot image using the gg package.
	ctx := gg.NewContext(int(dotSize), int(dotSize))

	// Draw a filled blue circle in the center
	ctx.DrawCircle(dotSize/2, dotSize/2, dotSize/2)
	ctx.SetColor(color.RGBA{B: 255, A: 255}) // Blue
	ctx.Fill()

	// Convert the resulting image to an *ebiten.Image
	dotImage = ebiten.NewImageFromImage(ctx.Image())
}

// DebugRectCacheKey is used to store unique dimensions in the debug image map.
type DebugRectCacheKey struct {
	Width  int
	Height int
}

// Global cache for debug images (rectangles).
var debugImageCache = make(map[DebugRectCacheKey]*ebiten.Image)

// Mutex to protect concurrent access to the cache (good practice).
var debugImageCacheMutex sync.Mutex

// getDebugRectImage checks the cache and creates/stores the image if missing.
func getDebugRectImage(r Rect) *ebiten.Image {
	// We only need integers for the cache key.
	w := int(r.Width())
	h := int(r.Height())

	// Safety: don't create images for zero/negative dimensions
	if w <= 0 || h <= 0 {
		return nil
	}

	key := DebugRectCacheKey{Width: w, Height: h}

	// Use mutex for thread-safe access
	debugImageCacheMutex.Lock()
	defer debugImageCacheMutex.Unlock()

	if img, ok := debugImageCache[key]; ok {
		return img
	}

	// If not found, create the image and store it
	img := createDebugRectImage(r)
	debugImageCache[key] = img
	return img
}

func createDebugRectImage(r Rect) *ebiten.Image {
	w := int(r.Width())
	h := int(r.Height())

	// Create a new drawing context for the size of the hitbox
	ctx := gg.NewContext(w, h)

	// Draw the rectangle with a 1-pixel red stroke
	ctx.SetColor(color.RGBA{R: 255, A: 255}) // Red
	ctx.SetLineWidth(1.0)

	// The rectangle starts at (0.5, 0.5) to keep the 1-pixel line entirely inside
	// the image bounds (standard gg practice for lines on integer coordinates).
	ctx.DrawRectangle(0.5, 0.5, float64(w)-1, float64(h)-1)
	ctx.Stroke()

	// Convert the result to an *ebiten.Image
	return ebiten.NewImageFromImage(ctx.Image())
}

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
	debugImage := getDebugRectImage(hb)

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
	debugImage := getDebugRectImage(pb)

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
