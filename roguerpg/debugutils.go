package main

import (
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

// GetDebugRectImage checks the cache and creates/stores the image if missing.
func GetDebugRectImage(r Rect) *ebiten.Image {
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
