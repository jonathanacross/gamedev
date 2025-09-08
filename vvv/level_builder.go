package main

import (
	"image"
	"log"
	"time"
	"vvv/tiled"

	"github.com/hajimehoshi/ebiten/v2"
)

func toRect(r tiled.Rect) Rect {
	return Rect{
		left:   r.X,
		top:    r.Y,
		right:  r.X + r.Width,
		bottom: r.Y + r.Height,
	}
}

func toImageRectangle(r tiled.Rect) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(r.X), Y: int(r.Y)},
		Max: image.Point{X: int(r.X + r.Width), Y: int(r.Y + r.Height)},
	}
}

func getLocation(o tiled.Object) Location {
	return Location{
		X: o.Location.X,
		Y: o.Location.Y,
	}
}

func GetLevelObjects(tm *tiled.Map, levelNum int) ([]GameObject, Location) {
	gameObjects := []GameObject{}
	var startPoint Location

	for _, layer := range tm.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				switch obj.Type {
				case "Spikes":
					spike := processSpikeObject(obj, tm.Tiles[obj.GID])
					gameObjects = append(gameObjects, spike)
				case "LevelExit":
					exit := processLevelExit(obj)
					gameObjects = append(gameObjects, exit)
				case "Checkpoint":
					checkpoint := processCheckpointObject(obj, tm.Tiles[obj.GID], levelNum)
					gameObjects = append(gameObjects, checkpoint)
					if checkpoint.Active {
						startPoint = checkpoint.Location
					}
				case "Platform":
					platform := processPlatformObject(obj, tm.Tiles[obj.GID])
					gameObjects = append(gameObjects, platform)
				case "HelicopterMonster":
					helicopterMonster := processHelicopterMonsterObject(obj, tm.Tiles[obj.GID])
					gameObjects = append(gameObjects, helicopterMonster)
				case "Crystal":
					crystal := processCrystalObject(obj, tm.Tiles[obj.GID])
					gameObjects = append(gameObjects, crystal)
				case "BreakingFloor":
					breakingFloor := processBreakingFloorObject(obj, tm.Tiles[obj.GID])
					gameObjects = append(gameObjects, breakingFloor)
				default:
					log.Printf("Unknown object type: %s.  Object = %v\n", obj.Type, obj)
				}
			}
		}
	}
	return gameObjects, startPoint
}

// New helper function to process a single Spike object.
func processSpikeObject(obj tiled.Object, tile tiled.Tile) *Spike {
	return &Spike{
		BaseSprite: BaseSprite{
			Location: getLocation(obj),
			image:    tile.SrcImage.(*ebiten.Image),
			srcRect:  toImageRectangle(tile.SrcRect),
			hitbox:   toRect(tile.HitRect).Offset(obj.Location.X, obj.Location.Y),
		},
	}
}

// New helper function to process a single LevelExit object.
func processLevelExit(obj tiled.Object) LevelExit {
	toLevel, err := obj.Properties.GetPropertyInt("ToLevel")
	if err != nil {
		log.Println("Error reading ToLevel property for LevelExit:", err)
	}
	return LevelExit{
		Rect:    toRect(obj.Location),
		ToLevel: toLevel,
	}
}

// New helper function to process a single Checkpoint object.
func processCheckpointObject(obj tiled.Object, tile tiled.Tile, levelNum int) *Checkpoint {
	isActive, err := obj.Properties.GetPropertyBool("Active")
	if err != nil {
		log.Println("Error reading Active property for Checkpoint:", err)
	}
	spriteSheet := NewGridTileSet(16, 16, 2, 1)
	srcRect2 := spriteSheet.Rect(0)
	checkpoint := Checkpoint{
		BaseSprite: BaseSprite{
			image:    CheckpointSprite,
			Location: getLocation(obj),
			srcRect:  srcRect2,
			hitbox:   toRect(tile.HitRect).Offset(obj.Location.X, obj.Location.Y),
		},
		spriteSheet: spriteSheet,
		Active:      isActive,
		Id:          levelNum*1000 + obj.GID,
		LevelNum:    levelNum,
	}
	checkpoint.SetActive(isActive)
	return &checkpoint
}

func processPlatformObject(obj tiled.Object, tile tiled.Tile) *Platform {
	low, err := obj.Properties.GetPropertyFloat64("low")
	if err != nil {
		log.Println("Error reading 'low' property for platform:", err)
	}
	high, err := obj.Properties.GetPropertyFloat64("high")
	if err != nil {
		log.Println("Error reading 'high' property for platform:", err)
	}
	delta, err := obj.Properties.GetPropertyFloat64("delta")
	if err != nil {
		log.Println("Error reading 'delta' property for platform:", err)
	}
	horiz, err := obj.Properties.GetPropertyBool("horiz")
	if err != nil {
		log.Println("Error reading 'horiz' property for platform:", err)
	}

	return &Platform{
		BaseSprite: BaseSprite{
			Location: getLocation(obj),
			image:    tile.SrcImage.(*ebiten.Image),
			srcRect:  toImageRectangle(tile.SrcRect),
			hitbox:   toRect(tile.HitRect).Offset(obj.Location.X, obj.Location.Y),
		},
		low:   low,
		high:  high,
		delta: delta,
		horiz: horiz,
	}
}

func processHelicopterMonsterObject(obj tiled.Object, tile tiled.Tile) *HelicopterMonster {
	low, err := obj.Properties.GetPropertyFloat64("low")
	if err != nil {
		log.Println("Error reading 'low' property for helicoptermonster:", err)
	}
	high, err := obj.Properties.GetPropertyFloat64("high")
	if err != nil {
		log.Println("Error reading 'high' property for helicoptermonster:", err)
	}
	delta, err := obj.Properties.GetPropertyFloat64("delta")
	if err != nil {
		log.Println("Error reading 'delta' property for helicoptermonster:", err)
	}
	horiz, err := obj.Properties.GetPropertyBool("horiz")
	if err != nil {
		log.Println("Error reading 'horiz' property for helicoptermonster:", err)
	}

	return &HelicopterMonster{
		BaseSprite: BaseSprite{
			Location: getLocation(obj),
			image:    MonsterSprite,
			srcRect:  toImageRectangle(tile.SrcRect),
			hitbox:   toRect(tile.HitRect).Offset(obj.Location.X, obj.Location.Y),
		},
		spriteSheet: NewGridTileSet(16, 16, 2, 1),
		animation:   NewAnimation(0, 1, 20),
		low:         low,
		high:        high,
		delta:       delta,
		horiz:       horiz,
	}
}

func processCrystalObject(obj tiled.Object, tile tiled.Tile) *Crystal {
	return &Crystal{
		BaseSprite: BaseSprite{
			Location: getLocation(obj),
			image:    tile.SrcImage.(*ebiten.Image),
			srcRect:  toImageRectangle(tile.SrcRect),
			hitbox:   toRect(tile.HitRect).Offset(obj.Location.X, obj.Location.Y),
		},
		Collected: false,
	}
}

func processBreakingFloorObject(obj tiled.Object, tile tiled.Tile) *BreakingFloor {
	return &BreakingFloor{
		BaseSprite: BaseSprite{
			Location: getLocation(obj),
			image:    BreakingFloorSprite,
			srcRect:  toImageRectangle(tile.SrcRect),
			hitbox:   toRect(tile.HitRect).Offset(obj.Location.X, obj.Location.Y),
		},
		spriteSheet: NewGridTileSet(16, 16, 5, 1),
		state:       Intact,
		breakTimer:  NewTimer(500 * time.Millisecond),
	}
}

func isSolid(tile tiled.Tile) bool {
	solid, _ := tile.Properties.GetPropertyBool("solid")
	return solid
}

func GetTiles(tm *tiled.Map) []Tile {
	tiles := []Tile{}
	for _, layer := range tm.Layers {
		if layer.Type == "tilelayer" {
			for idx, id := range layer.TileIds {
				if id <= 0 {
					continue
				}
				t := tm.Tiles[id]
				x := (idx % layer.Width) * TileSize
				y := (idx / layer.Width) * TileSize
				tile := Tile{
					BaseSprite: BaseSprite{
						Location: Location{
							X: float64(x),
							Y: float64(y),
						},
						image:   t.SrcImage.(*ebiten.Image),
						srcRect: toImageRectangle(t.SrcRect),
						hitbox: Rect{
							left:   float64(x),
							top:    float64(y),
							right:  float64(x + TileSize),
							bottom: float64(y + TileSize),
						},
					},
					solid: isSolid(t),
				}
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles
}
