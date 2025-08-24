package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

func (r Rect) Width() float64 {
	return r.right - r.left
}

func (r Rect) Height() float64 {
	return r.bottom - r.top
}

type Object struct {
	Rect
}

type LevelExit struct {
	Rect
	ToLevel int
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
}

type Level struct {
	tilemapJson TilemapJSON
	spriteSheet *SpriteSheet
	tiles       *[]Tile
	exits       []LevelExit
	width       float64
	height      float64
}

// Hack: these functions depend on the particular number/location
// of tiles in the tileset.
// TODO: read from tileset properties
func isSolid(id int) bool    { return id <= 20 }
func isDamaging(id int) bool { return id == 22 || id == 23 }

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

func getLevelExits(tilemapJson TilemapJSON) []LevelExit {
	exits := []LevelExit{}
	for _, layer := range tilemapJson.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				if obj.Type == "LevelExit" {
					toLevel := 0
					for _, prop := range obj.Properties {
						if prop.Name == "ToLevel" {
							var ok bool
							toLevel, ok = prop.IntValue()
							if !ok {
								// If the value couldn't be decoded, skip this object or handle the error.
								continue
							}
							break
						}
					}
					exit := LevelExit{
						Rect: Rect{
							left:   obj.X,
							top:    obj.Y,
							right:  obj.X + obj.Width,
							bottom: obj.Y + obj.Height,
						},
						ToLevel: toLevel,
					}
					exits = append(exits, exit)
				}
			}
		}
	}
	return exits
}

func NewLevel(tilemapJson TilemapJSON, spriteSheet *SpriteSheet) *Level {
	return &Level{
		tilemapJson: tilemapJson,
		spriteSheet: spriteSheet,
		tiles:       getTiles(tilemapJson, spriteSheet),
		exits:       getLevelExits(tilemapJson),
		width:       float64(tilemapJson.Width * TileSize),
		height:      float64(tilemapJson.Height * TileSize),
	}
}

// DrawRectFrame draws a 1-pixel wide frame around the given Rect with the specified color.
func DrawRectFrame(screen *ebiten.Image, rect Rect, clr color.RGBA) {
	lineWidth := float32(1)

	// Draw top line
	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.right), float32(rect.top), lineWidth, clr, false)
	// Draw bottom line
	vector.StrokeLine(screen, float32(rect.left), float32(rect.bottom), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
	// Draw left line
	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.left), float32(rect.bottom), lineWidth, clr, false)
	// Draw right line
	vector.StrokeLine(screen, float32(rect.right), float32(rect.top), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
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

	// Draw the exits for debugging
	for _, exit := range level.exits {
		DrawRectFrame(screen, exit.Rect, color.RGBA{255, 0, 0, 255})
	}
}

func (level *Level) Update() {
}
