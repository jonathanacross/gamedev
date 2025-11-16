package main

import (
	"image"
	"image/color"

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
type BaseSprite struct {
	Location
	image      *ebiten.Image
	srcRect    image.Rectangle
	hitbox     Rect
	drawOffset Location
	debugImage *ebiten.Image
}

// HitBox returns the collision rectangle for the BaseSprite.
func (bs *BaseSprite) HitBox() Rect {
	return bs.hitbox.Offset(bs.X, bs.Y)
}

func (bs *BaseSprite) GetX() float64 { return bs.X }

func (bs *BaseSprite) GetY() float64 { return bs.Y }

func (bs *BaseSprite) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	if !ShowDebugInfo {
		return
	}

	if bs.debugImage == nil || dotImage == nil {
		return
	}

	// Draw the Hitbox rectangle
	hb := bs.HitBox()

	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(hb.Left, hb.Top)
	opRect.GeoM.Concat(cameraMatrix)
	screen.DrawImage(bs.debugImage, opRect)

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

type Tile struct {
	BaseSprite
	solid bool
}

func NewTile(location Location, image *ebiten.Image, srcRect image.Rectangle, solid bool) *Tile {
	hitBox := Rect{
		Left:   0,
		Top:    0,
		Right:  float64(srcRect.Dx()),
		Bottom: float64(srcRect.Dy()),
	}
	return &Tile{
		BaseSprite: BaseSprite{
			Location: location,
			image:    image,
			srcRect:  srcRect,
			hitbox:   hitBox,
			drawOffset: Location{
				X: 0,
				Y: 0,
			},
			debugImage: createDebugRectImage(hitBox),
		},
		solid: solid,
	}
}
