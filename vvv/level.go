package main

import (
	"image"
	"image/color"
	"log"

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

type LevelExit struct {
	Rect
	ToLevel int
}

type BaseSprite struct {
	Location
	spriteSheet *SpriteSheet
	srcRect     image.Rectangle
}

func (bs *BaseSprite) HitRect() Rect {
	return Rect{
		left:   bs.X,
		top:    bs.Y,
		right:  bs.X + float64(bs.spriteSheet.tileWidth),
		bottom: bs.Y + float64(bs.spriteSheet.tileHeight),
	}
}

func (bs *BaseSprite) GetX() float64 { return bs.X }

func (bs *BaseSprite) GetY() float64 { return bs.Y }

func (bs *BaseSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.X, bs.Y)
	currImage := bs.spriteSheet.image.SubImage(bs.srcRect).(*ebiten.Image)
	screen.DrawImage(currImage, op)
}

type Tile struct {
	BaseSprite
	solid bool
}

type Spike struct {
	BaseSprite
	hitbox Rect
}

type Level struct {
	tilemapJson LevelJSON
	spriteSheet *SpriteSheet
	tiles       *[]Tile
	spikes      []Spike
	exits       []LevelExit
	width       float64
	height      float64
}

func DrawRectFrame(screen *ebiten.Image, rect Rect, clr color.RGBA) {
	lineWidth := float32(1)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.right), float32(rect.top), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.bottom), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.left), float32(rect.top), float32(rect.left), float32(rect.bottom), lineWidth, clr, false)
	vector.StrokeLine(screen, float32(rect.right), float32(rect.top), float32(rect.right), float32(rect.bottom), lineWidth, clr, false)
}

func findTilesetTileData(tilesetData TilesetDataJSON, gid int) *TilesetTileJSON {
	for _, tile := range tilesetData.Tiles {
		if tile.ID == gid-1 {
			return &tile
		}
	}
	return nil
}

func isSolid(tilesetData TilesetDataJSON, id int) bool {
	tileData := findTilesetTileData(tilesetData, id)
	if tileData == nil {
		return false
	}
	for _, prop := range tileData.Properties {
		if prop.Name == "solid" {
			if isSolid, ok := prop.BoolValue(); ok {
				return isSolid
			}
		}
	}
	return false
}

func getHitboxFromTileData(obj ObjectJSON, tilesetTileData *TilesetTileJSON) Rect {
	if tilesetTileData != nil && len(tilesetTileData.ObjectGroup.Objects) > 0 {
		rectData := tilesetTileData.ObjectGroup.Objects[0]
		return Rect{
			left:   obj.X + rectData.X,
			top:    obj.Y + rectData.Y,
			right:  obj.X + rectData.X + rectData.Width,
			bottom: obj.Y + rectData.Y + rectData.Height,
		}
	} else {
		return Rect{
			left:   obj.X,
			top:    obj.Y,
			right:  obj.X + obj.Width,
			bottom: obj.Y + obj.Height,
		}
	}
}

func getLevelObjectsAndExits(leveljson LevelJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet) ([]Spike, []LevelExit) {
	spikes := []Spike{}
	exits := []LevelExit{}
	for _, layer := range leveljson.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				objType := obj.Type
				if objType == "" && obj.Gid > 0 {
					tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
					if tilesetTileData != nil && len(tilesetTileData.ObjectGroup.Objects) > 0 && tilesetTileData.ObjectGroup.Objects[0].Name == "Spikes" {
						objType = "Spike"
					}
				}

				switch objType {
				case "Spike":
					tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
					if tilesetTileData == nil {
						log.Println("Tileset tile data not found for Spike, Gid:", obj.Gid)
						continue
					}
					// Tiled JSON uses the bottom-left corner for an object's y-position.
					// We need to adjust it to be the top-left for drawing.
					adjustedY := obj.Y - obj.Height

					// Recalculate the hitbox based on the corrected y-position
					var hitbox Rect
					if len(tilesetTileData.ObjectGroup.Objects) > 0 {
						rectData := tilesetTileData.ObjectGroup.Objects[0]
						hitbox = Rect{
							left:   obj.X + rectData.X,
							top:    adjustedY + rectData.Y,
							right:  obj.X + rectData.X + rectData.Width,
							bottom: adjustedY + rectData.Y + rectData.Height,
						}
					} else {
						// Fallback to the object's dimensions if no specific hitbox is defined
						hitbox = Rect{
							left:   obj.X,
							top:    adjustedY,
							right:  obj.X + obj.Width,
							bottom: adjustedY + obj.Height,
						}
					}

					spike := Spike{
						BaseSprite: BaseSprite{
							Location: Location{
								X: obj.X,
								Y: adjustedY,
							},
							spriteSheet: spriteSheet,
							srcRect:     spriteSheet.Rect(obj.Gid - 1),
						},
						hitbox: hitbox,
					}
					spikes = append(spikes, spike)
				case "LevelExit":
					toLevel := 0
					for _, prop := range obj.Properties {
						if prop.Name == "ToLevel" {
							if toLevelVal, ok := prop.IntValue(); ok {
								toLevel = toLevelVal
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
	return spikes, exits
}

func getTiles(leveljson LevelJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet) *[]Tile {
	tiles := []Tile{}
	for layerIdx, layer := range leveljson.Layers {
		if layerIdx != 0 {
			continue
		}
		for idx, id := range layer.Data {
			x := (idx % layer.Width) * TileSize
			y := (idx / layer.Width) * TileSize
			tile := Tile{
				BaseSprite: BaseSprite{
					Location: Location{
						X: float64(x),
						Y: float64(y),
					},
					spriteSheet: spriteSheet,
					srcRect:     spriteSheet.Rect(id - 1),
				},
				solid: isSolid(tilesetData, id),
			}
			tiles = append(tiles, tile)
		}
	}
	return &tiles
}

func NewLevel(leveljson LevelJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet) *Level {
	spikes, exits := getLevelObjectsAndExits(leveljson, tilesetData, spriteSheet)
	return &Level{
		tilemapJson: leveljson,
		spriteSheet: spriteSheet,
		tiles:       getTiles(leveljson, tilesetData, spriteSheet),
		spikes:      spikes,
		exits:       exits,
		width:       float64(leveljson.Width * TileSize),
		height:      float64(leveljson.Height * TileSize),
	}
}

func (level *Level) Draw(screen *ebiten.Image) {
	for _, tile := range *level.tiles {
		tile.Draw(screen)
	}
	for _, spike := range level.spikes {
		spike.BaseSprite.Draw(screen)
		DrawRectFrame(screen, spike.hitbox, color.RGBA{255, 165, 0, 255})
	}
	for _, exit := range level.exits {
		DrawRectFrame(screen, exit.Rect, color.RGBA{0, 255, 0, 255})
	}
}

func (level *Level) Update() {
}
