package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const TileSize = 16

type Location struct {
	X float64
	Y float64
}

type Rect struct {
	left   float64
	top    float64
	right  float64
	bottom float64
}

// BaseSprite provides common fields and methods for any visible game entity.
// It handles drawing a single sprite or the current frame of an animation.
type BaseSprite struct {
	Location
	spriteSheet *SpriteSheet
	srcRect     image.Rectangle // The specific rectangle on the sprite sheet to draw
}

// HitRect returns the collision rectangle for the BaseSprite.
func (bs *BaseSprite) HitRect() Rect {
	return Rect{
		left:   bs.X,
		top:    bs.Y,
		right:  bs.X + float64(bs.spriteSheet.tileWidth),
		bottom: bs.Y + float64(bs.spriteSheet.tileHeight),
	}
}

// GetX returns the X coordinate of the BaseSprite.
func (bs *BaseSprite) GetX() float64 { return bs.X }

// GetY returns the Y coordinate of the BaseSprite.
func (bs *BaseSprite) GetY() float64 { return bs.Y }

type Tile struct {
	BaseSprite
	solid    bool
	damaging bool
	isLeft   bool
	isRight  bool
	isUp     bool
	isDown   bool
}

type Level struct {
	tilemapJson TilemapJSON
	spriteSheet *SpriteSheet
	tiles       *[]Tile
	levelLeft   int
	levelRight  int
	levelUp     int
	levelDown   int
}

// Hack: these functions depend on the particular number/location
// of tiles in the tileset.
// TODO: read from tileset properties
// TODO: consider directions to objects with level ids
func isSolid(id int) bool    { return id <= 20 }
func isDamaging(id int) bool { return id == 22 || id == 23 }
func isLeft(id int) bool     { return id == 31 }
func isRight(id int) bool    { return id == 32 }
func isUp(id int) bool       { return id == 33 }
func isDown(id int) bool     { return id == 34 }

func getTiles(tilemapJson TilemapJSON, spriteSheet *SpriteSheet) *[]Tile {
	tiles := []Tile{}

	for layerIdx, layer := range tilemapJson.Layers {
		if layerIdx != 0 {
			// TODO: consider filtering based on layer name
			continue
		}
		for idx, id := range layer.Data {
			x := (idx % layer.Width) * TileSize
			y := (idx / layer.Width) * TileSize

			// Json stores ids as one-based
			tile := Tile{
				BaseSprite: BaseSprite{
					Location: Location{
						X: float64(x),
						Y: float64(y),
					},
					spriteSheet: spriteSheet,
					srcRect:     spriteSheet.Rect(id - 1),
				},
				solid:    isSolid(id),
				damaging: isDamaging(id),
			}
			tiles = append(tiles, tile)
		}
	}

	return &tiles
}

func NewLevel(tilemapJson TilemapJSON, spriteSheet *SpriteSheet) *Level {
	return &Level{
		tilemapJson: tilemapJson,
		spriteSheet: spriteSheet,
		tiles:       getTiles(tilemapJson, spriteSheet),
	}
}

func (bs *BaseSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X, bs.Y)
	currImage := bs.spriteSheet.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

func (level *Level) Draw(screen *ebiten.Image) {
	for _, tile := range *level.tiles {
		tile.Draw(screen)
	}
}

func (level *Level) Update() {
}
