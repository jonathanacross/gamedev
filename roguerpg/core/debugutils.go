package core

import (
	"image/color"
	"sync"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// Size of the debug dot in pixels
const dotSize = 2.0

// Global ebiten.Image for the location marker dot (initialized in init()).
var dotImage *ebiten.Image

func init() {
	// Create the dot image using the gg package.
	ctx := gg.NewContext(int(dotSize), int(dotSize))

	// Draw a filled green circle in the center
	// TODO: draw crosshairs instead of a circle
	ctx.DrawCircle(dotSize/2, dotSize/2, dotSize/2)
	ctx.SetColor(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green
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

	// Draw the rectangle with a 1-pixel white stroke
	ctx.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	ctx.SetLineWidth(1.0)

	// The rectangle starts at (0.5, 0.5) to keep the 1-pixel line entirely inside
	// the image bounds (standard gg practice for lines on integer coordinates).
	ctx.DrawRectangle(0.5, 0.5, float64(w)-1, float64(h)-1)
	ctx.Stroke()

	// Convert the result to an *ebiten.Image
	return ebiten.NewImageFromImage(ctx.Image())
}

func DrawDebugRect(screen *ebiten.Image, rect Rect, color color.Color, cameraMatrix ebiten.GeoM) {
	debugImage := getDebugRectImage(rect)
	if debugImage == nil {
		return
	}

	var cm colorm.ColorM
	r, g, b, a := color.RGBA()
	cm.Scale(float64(r)/255.0, float64(g)/255.0, float64(b)/255.0, float64(a)/255.0)

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(rect.Left, rect.Top)
	op.GeoM.Concat(cameraMatrix)
	colorm.DrawImage(screen, debugImage, cm, op)
}

func DrawLocationDot(screen *ebiten.Image, loc Location, cameraMatrix ebiten.GeoM) {
	if dotImage == nil {
		return
	}

	opDot := &ebiten.DrawImageOptions{}
	opDot.GeoM.Translate(loc.X-dotSize/2, loc.Y-dotSize/2)
	opDot.GeoM.Concat(cameraMatrix)
	screen.DrawImage(dotImage, opDot)
}
