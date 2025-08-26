package main

import (
	"image/color"
	"vvv/tiled"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Level struct {
	spriteSheet *GridTileSet
	tiles       []Tile
	spikes      []Spike
	exits       []LevelExit
	checkpoints []*Checkpoint
	platforms   []*Platform
	width       float64
	height      float64
	startPoint  Location
	levelImage  *ebiten.Image
}

// DrawRectFrame draws a 1-pixel wide frame around the given Rect with the specified color.
func DrawRectFrame(screen *ebiten.Image, rect Rect, clr color.RGBA) {
	lineWidth := float32(1)

	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.right), float32(rect.top), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.bottom), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.left), float32(rect.bottom), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.right), float32(rect.top), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
}

func NewLevel(leveljson tiled.LevelJSON, tilesetData tiled.TilesetDataJSON, spriteSheet *GridTileSet, levelNum int) *Level {
	spikes, exits, checkpoints, platforms, startPoint := GetLevelObjects(leveljson, tilesetData, spriteSheet, levelNum)

	levelImage := ebiten.NewImage(leveljson.Width*TileSize, leveljson.Height*TileSize)

	// Get all tiles and draw them to the offscreen image.
	tiles := GetTiles(leveljson, tilesetData, spriteSheet)
	for _, tile := range tiles {
		tile.Draw(levelImage)
	}

	return &Level{
		spriteSheet: spriteSheet,
		tiles:       tiles,
		spikes:      spikes,
		exits:       exits,
		checkpoints: checkpoints,
		platforms:   platforms,
		width:       float64(leveljson.Width * TileSize),
		height:      float64(leveljson.Height * TileSize),
		startPoint:  startPoint,
		levelImage:  levelImage,
	}
}

func (level *Level) FindCheckpoint(id int) *Checkpoint {
	for _, cp := range level.checkpoints {
		if cp.Id == id {
			return cp
		}
	}
	return nil
}

func (level *Level) Draw(screen *ebiten.Image, debug bool) {
	// Draw the pre-rendered level image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(level.levelImage, op)

	// Draw dynamic objects (spikes, exits, checkpoints)
	for _, spike := range level.spikes {
		spike.BaseSprite.Draw(screen)
		if debug {
			DrawRectFrame(screen, spike.hitbox, color.RGBA{255, 165, 0, 255})
		}
	}
	for _, exit := range level.exits {
		if debug {
			DrawRectFrame(screen, exit.Rect, color.RGBA{0, 255, 0, 255})
		}
	}
	for _, cp := range level.checkpoints {
		cp.Draw(screen)
		if debug {
			DrawRectFrame(screen, cp.hitbox, color.RGBA{0, 0, 255, 255})
		}
	}
	for _, platform := range level.platforms {
		platform.Draw(screen)
		if debug {
			DrawRectFrame(screen, platform.hitbox, color.RGBA{0, 128, 255, 255})
		}
	}
}

func (level *Level) Update() {
}

func (c *Checkpoint) SetActive(active bool) {
	c.Active = active
	if active {
		c.srcRect = c.spriteSheet.Rect(23)
	} else {
		c.srcRect = c.spriteSheet.Rect(24)
	}
}
