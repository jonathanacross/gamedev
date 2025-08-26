package main

import (
	"log"
	"vvv/tiled"
)

func GetLevelObjects(leveljson tiled.LevelJSON, tilesetData tiled.TilesetDataJSON, spriteSheet *GridTileSet, levelNum int) ([]GameObject, Location) {
	gameObjects := []GameObject{}
	var startPoint Location

	for _, layer := range leveljson.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				objType := obj.Type
				if objType == "" && obj.Gid > 0 {
					tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
					objType = tilesetTileData.Type
				}

				switch objType {
				case "Spikes":
					spike := processSpikeObject(obj, tilesetData, spriteSheet)
					if spike != nil {
						gameObjects = append(gameObjects, spike)
					}

				case "LevelExit":
					exit := processLevelExit(obj)
					gameObjects = append(gameObjects, exit)
				case "Checkpoint":
					checkpoint := processCheckpointObject(obj, tilesetData, spriteSheet, levelNum)
					gameObjects = append(gameObjects, checkpoint)
					if checkpoint.Active {
						startPoint = checkpoint.Location
					}
				case "Platform":
					platform := processPlatformObject(obj, tilesetData, spriteSheet)
					if platform != nil {
						gameObjects = append(gameObjects, platform)
					}
				}
			}
		}
	}
	return gameObjects, startPoint
}

// New helper function to process a single Spike object.
func processSpikeObject(obj tiled.ObjectJSON, tilesetData tiled.TilesetDataJSON, spriteSheet *GridTileSet) *Spike {
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
			hitbox:      hitbox,
		},
	}
}

// New helper function to process a single LevelExit object.
func processLevelExit(obj tiled.ObjectJSON) LevelExit {
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
func processCheckpointObject(obj tiled.ObjectJSON, tilesetData tiled.TilesetDataJSON, spriteSheet *GridTileSet, levelNum int) *Checkpoint {
	adjustedY := obj.Y - obj.Height

	isActive := false
	tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
	if tilesetTileData != nil {
		if activeVal, ok := tiled.GetBoolProperty(tilesetTileData.Properties, "Active"); ok {
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
			hitbox: Rect{
				left:   obj.X,
				top:    adjustedY,
				right:  obj.X + obj.Width,
				bottom: adjustedY + obj.Height,
			},
		},
		Active:   isActive,
		Id:       obj.ID,
		LevelNum: levelNum,
	}
}

// New helper function to process a single Platform object.
func processPlatformObject(obj tiled.ObjectJSON, tilesetData tiled.TilesetDataJSON, spriteSheet *GridTileSet) *Platform {
	adjustedY := obj.Y - obj.Height

	// Read custom properties for movement from the JSON
	var endX, endY float64

	tilesetTileData := findTilesetTileData(tilesetData, obj.Gid)
	if tilesetTileData != nil {
		// TODO: read in platform data from json files
		// if val, ok := getStringProperty(tilesetTileData.Properties, "endPoint"); ok {
		// 	// Assuming "endPoint" is in the format "x,y"
		// 	// You'll need a helper function to parse this string
		// 	// For now, let's assume you'll manually set the destination
		// 	endX = obj.X + obj.Width
		// 	endY = adjustedY
		// }
	}

	return &Platform{
		BaseSprite: BaseSprite{
			Location: Location{
				X: obj.X,
				Y: adjustedY,
			},
			spriteSheet: spriteSheet,
			srcRect:     spriteSheet.Rect(obj.Gid - 1),
			hitbox: Rect{
				left:   obj.X,
				top:    adjustedY,
				right:  obj.X + obj.Width,
				bottom: adjustedY + obj.Height,
			},
		},
		startX: obj.X,
		startY: adjustedY,
		endX:   endX,
		endY:   endY,
	}
}

// findTilesetTileData returns the TilesetTileJSON for a given GID.
func findTilesetTileData(tilesetData tiled.TilesetDataJSON, gid int) *tiled.TilesetTileJSON {
	for _, tile := range tilesetData.Tiles {
		// TODO: gid needs to be computed from JSON filed instead of
		// assuming the offset is 1.
		if tile.ID == gid-1 {
			return &tile
		}
	}
	return nil
}

func isSolid(tilesetData tiled.TilesetDataJSON, id int) bool {
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

func GetTiles(leveljson tiled.LevelJSON, tilesetData tiled.TilesetDataJSON, spriteSheet *GridTileSet) []Tile {
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
						hitbox: Rect{
							left:   float64(x),
							top:    float64(y),
							right:  float64(x + TileSize),
							bottom: float64(y + TileSize),
						},
					},
					solid: isSolid(tilesetData, id),
				}
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles
}
