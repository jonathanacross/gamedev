package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const TileSize = 16

type LevelExit struct {
	Rect
	ToLevel int
}

// BaseSprite provides common fields and methods for any visible game entity.
// It handles drawing a single sprite or the current frame of an animation.
type BaseSprite struct {
	Location
	spriteSheet *SpriteSheet
	srcRect     image.Rectangle
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

type Checkpoint struct {
	BaseSprite
	hitbox   Rect
	Active   bool
	Id       int
	LevelNum int
}

type Level struct {
	tilemapJson LevelJSON
	spriteSheet *SpriteSheet
	tiles       []Tile
	spikes      []Spike
	exits       []LevelExit
	checkpoints []*Checkpoint
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

// findTilesetTileData returns the TilesetTileJSON for a given GID.
func findTilesetTileData(tilesetData TilesetDataJSON, gid int) *TilesetTileJSON {
	for _, tile := range tilesetData.Tiles {
		// The `id` in tileset JSON is 0-based, while `gid` in level JSON is 1-based.
		// We need to account for this and any potential offset from `firstgid`.
		// For simplicity, let's assume `firstgid` is always 1 for now.
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

func getLevelObjects(leveljson LevelJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet, levelNum int) ([]Spike, []LevelExit, []*Checkpoint, Location) {
	spikes := []Spike{}
	exits := []LevelExit{}
	checkpoints := []*Checkpoint{}
	var startPoint Location

	for _, layer := range leveljson.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				objType := obj.Type
				if objType == "" && obj.Gid > 0 {
					tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
					if tilesetTileData != nil && len(tilesetTileData.ObjectGroup.Objects) > 0 && tilesetTileData.ObjectGroup.Objects[0].Name == "Spikes" {
						objType = "Spike"
					}
					if tilesetTileData != nil {
						if tileType, ok := getStringProperty(tilesetTileData.Properties, "Type"); ok {
							objType = tileType
						}
					}
				}

				switch objType {
				case "Spike":
					spike := processSpikeObject(obj, tilesetData, spriteSheet)
					if spike != nil {
						spikes = append(spikes, *spike)
					}
				case "LevelExit":
					exit := processLevelExit(obj)
					exits = append(exits, exit)
				case "Checkpoint":
					checkpoint := processCheckpointObject(obj, tilesetData, spriteSheet, levelNum)
					checkpoints = append(checkpoints, checkpoint)
					if checkpoint.Active {
						startPoint = checkpoint.Location
					}
				}
			}
		}
	}
	return spikes, exits, checkpoints, startPoint
}

// New helper function to process a single Spike object.
func processSpikeObject(obj ObjectJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet) *Spike {
	tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
	if tilesetTileData == nil {
		log.Println("Tileset tile data not found for Spike, Gid:", obj.Gid)
		return nil
	}
	adjustedY := obj.Y - obj.Height

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
		hitbox = Rect{
			left:   obj.X,
			top:    adjustedY,
			right:  obj.X + obj.Width,
			bottom: adjustedY + obj.Height,
		}
	}

	return &Spike{
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
}

// New helper function to process a single LevelExit object.
func processLevelExit(obj ObjectJSON) LevelExit {
	toLevel := 0
	for _, prop := range obj.Properties {
		if prop.Name == "ToLevel" {
			if toLevelVal, ok := prop.IntValue(); ok {
				toLevel = toLevelVal
			}
			break
		}
	}
	return LevelExit{
		Rect: Rect{
			left:   obj.X,
			top:    obj.Y,
			right:  obj.X + obj.Width,
			bottom: obj.Y + obj.Height,
		},
		ToLevel: toLevel,
	}
}

// New helper function to process a single Checkpoint object.
func processCheckpointObject(obj ObjectJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet, levelNum int) *Checkpoint {
	adjustedY := obj.Y - obj.Height

	isActive := false
	tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
	if tilesetTileData != nil {
		if activeVal, ok := getBoolProperty(tilesetTileData.Properties, "Active"); ok {
			isActive = activeVal
		}
	}

	return &Checkpoint{
		BaseSprite: BaseSprite{
			Location: Location{
				X: obj.X,
				Y: adjustedY,
			},
			spriteSheet: spriteSheet,
			srcRect:     spriteSheet.Rect(obj.Gid - 1),
		},
		hitbox: Rect{
			left:   obj.X,
			top:    adjustedY,
			right:  obj.X + obj.Width,
			bottom: adjustedY + obj.Height,
		},
		Active:   isActive,
		Id:       obj.ID,
		LevelNum: levelNum,
	}
}

func getTiles(leveljson LevelJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet) []Tile {
	tiles := []Tile{}
	for _, layer := range leveljson.Layers {
		if layer.Type == "tilelayer" {
			for idx, id := range layer.Data {
				if id == 0 {
					continue
				}
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
	}
	return tiles
}

func NewLevel(leveljson LevelJSON, tilesetData TilesetDataJSON, spriteSheet *SpriteSheet, levelNum int) *Level {
	spikes, exits, checkpoints, startPoint := getLevelObjects(leveljson, tilesetData, spriteSheet, levelNum)

	levelImage := ebiten.NewImage(leveljson.Width*TileSize, leveljson.Height*TileSize)

	// Get all tiles and draw them to the offscreen image.
	tiles := getTiles(leveljson, tilesetData, spriteSheet)
	for _, tile := range tiles {
		tile.Draw(levelImage)
	}

	return &Level{
		tilemapJson: leveljson,
		spriteSheet: spriteSheet,
		tiles:       tiles,
		spikes:      spikes,
		exits:       exits,
		checkpoints: checkpoints,
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
