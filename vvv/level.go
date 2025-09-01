package main

import (
	"image/color"
	"vvv/tiled"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Level struct {
	tiles      []Tile
	objects    []GameObject
	width      float64
	height     float64
	startPoint Location
	levelImage *ebiten.Image
}

func NewLevel(tm *tiled.Map, levelNum int) *Level {
	objects, startPoint := GetLevelObjects(tm, levelNum)

	levelImage := ebiten.NewImage(tm.WidthInTiles*tm.TileWidth, tm.HeightInTiles*tm.TileHeight)

	// Get all tiles and draw them to the offscreen image.
	tiles := GetTiles(tm)
	for _, tile := range tiles {
		tile.Draw(levelImage)
	}

	return &Level{
		tiles:      tiles,
		objects:    objects,
		width:      float64(tm.WidthInTiles * tm.TileWidth),
		height:     float64(tm.HeightInTiles * TileSize),
		startPoint: startPoint,
		levelImage: levelImage,
	}
}

func (level *Level) FindCheckpoint(id int) *Checkpoint {
	// Find the checkpoint by its ID.
	for _, obj := range level.objects {
		if cp, ok := obj.(*Checkpoint); ok {
			if cp.Id == id {
				return cp
			}
		}
	}
	return nil
}

// DrawRectFrame draws a 1-pixel wide frame around the given Rect with the specified color.
func DrawRectFrame(screen *ebiten.Image, rect Rect, clr color.RGBA) {
	lineWidth := float32(1)

	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.right), float32(rect.top), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.bottom), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.left), float32(rect.bottom), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.right), float32(rect.top), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
}

func (level *Level) Draw(screen *ebiten.Image, debug bool) {
	// Draw the pre-rendered level image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(level.levelImage, op)

	// Draw dynamic objects (spikes, exits, checkpoints)
	for _, object := range level.objects {
		if d, ok := object.(Drawable); ok {
			d.Draw(screen)
		}
		if debug {
			DrawRectFrame(screen, object.HitBox(), color.RGBA{255, 255, 255, 255})
		}
	}
}

func (level *Level) Update() {
}
