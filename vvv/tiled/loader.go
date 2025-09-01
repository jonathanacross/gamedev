package tiled

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// The Loader interface defines the contract for loading Tiled data.
// It abstracts away the details of file system access and caching.
type Loader interface {
	LoadMap(path string) (*Map, error)
}

// ImageConverter is a function type that converts a standard image.Image to a custom type.
type ImageConverter func(img image.Image) (ImageProvider, error)

// FsLoader is an implementation of the Loader interface that
// uses an embedded file system (fs.FS) to read files.
type FsLoader struct {
	fs fs.FS

	// Cache for loaded files to avoid redundant reads.
	// We use a map where the key is the file path and the value is the raw JSON data.
	// The mutex protects concurrent access to the cache.
	cache map[string][]byte
	mu    sync.Mutex

	// Cache for loaded images.
	imageCache map[string]ImageProvider
	imageMu    sync.Mutex

	converter ImageConverter
}

// NewFsLoader creates a new FsLoader instance.
// This uses the default *image.Image type for images.
func NewFsLoader(fsys fs.FS) *FsLoader {
	return &FsLoader{
		fs:         fsys,
		cache:      make(map[string][]byte),
		imageCache: make(map[string]ImageProvider),
		converter:  func(img image.Image) (ImageProvider, error) { return img, nil },
	}
}

// NewFsLoaderWithImageConverter allows users to provide a custom image converter.
func NewFsLoaderWithImageConverter(fsys fs.FS, converter ImageConverter) *FsLoader {
	return &FsLoader{
		fs:         fsys,
		cache:      make(map[string][]byte),
		imageCache: make(map[string]ImageProvider),
		converter:  converter,
	}
}

// LoadMap loads a Tiled map from the specified path.
func (l *FsLoader) LoadMap(filePath string) (*Map, error) {
	// Step 1: Load the raw map data.
	tiledMapData, err := l.loadMapData(filePath)
	if err != nil {
		return nil, err
	}

	// Step 2: Load and convert all tilesets.
	allTiles, err := l.loadTilesets(tiledMapData.Tilesets, path.Dir(filePath))
	if err != nil {
		return nil, err
	}

	// Step 3: Convert the raw layers to game layers.
	gameLayers, err := convertLayers(tiledMapData.Layers, &allTiles)
	if err != nil {
		return nil, err
	}

	// Step 4: Construct and return the final game map.
	gameMap := &Map{
		Name:          filepath.Base(filePath),
		WidthInTiles:  tiledMapData.Width,
		HeightInTiles: tiledMapData.Height,
		TileWidth:     tiledMapData.TileWidth,
		TileHeight:    tiledMapData.TileHeight,
		Layers:        gameLayers,
		Tiles:         allTiles,
	}

	// // Populate the Tiles slice from the map for easier access later.
	// for _, tile := range allTiles {
	// 	gameMap.Tiles = append(gameMap.Tiles, tile)
	// }

	return gameMap, nil
}

// loadMapData handles reading and unmarshalling the main map JSON file.
func (l *FsLoader) loadMapData(filePath string) (*tiledMap, error) {
	data, err := l.loadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load map file %s: %w", filePath, err)
	}

	var tiledMap tiledMap
	if err := json.Unmarshal(data, &tiledMap); err != nil {
		return nil, fmt.Errorf("failed to parse map JSON %s: %w", filePath, err)
	}

	return &tiledMap, nil
}

// loadTilesets iterates through the tileset references in the map and loads them.
func (l *FsLoader) loadTilesets(tsRefs []tiledTileset, mapDir string) (map[int]Tile, error) {
	allTiles := make(map[int]Tile)

	for _, tsRef := range tsRefs {
		var tsData tiledTileset
		tsPath := mapDir

		if tsRef.Source != "" {
			// Normalize the path by replacing backslashes with forward slashes.
			normalizedSource := strings.ReplaceAll(tsRef.Source, "\\", "/")
			tsPath = path.Join(mapDir, normalizedSource)
			data, err := l.loadFile(tsPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load tileset file %s: %w", tsPath, err)
			}
			if err := json.Unmarshal(data, &tsData); err != nil {
				return nil, fmt.Errorf("failed to parse tileset JSON %s: %w", tsPath, err)
			}
		} else {
			tsData = tsRef
		}

		imageMap, err := l.loadTilesetImages(&tsData, tsPath)
		if err != nil {
			return nil, err
		}

		convertedTiles, err := ConvertTileset(&tsData, imageMap, tsRef.FirstGID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tileset %s: %w", tsPath, err)
		}

		for _, tile := range convertedTiles {
			allTiles[tile.ID] = tile
		}
	}

	return allTiles, nil
}

// loadTilesetImages loads all images associated with a tileset.
func (l *FsLoader) loadTilesetImages(tsData *tiledTileset, tsPath string) (map[string]ImageProvider, error) {
	imageMap := make(map[string]ImageProvider)

	// Get the directory of the tileset file to correctly resolve relative image paths.
	tsDir := path.Dir(tsPath)

	if tsData.Image != "" { // Sprite sheet tileset
		normalizedImage := strings.ReplaceAll(tsData.Image, "\\", "/")
		imgPath := path.Join(tsDir, normalizedImage)
		img, err := l.loadImage(imgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tileset image %s: %w", imgPath, err)
		}
		imageMap[tsData.Image] = img
	} else { // Collection tileset
		for _, tile := range tsData.Tiles {
			if tile.Image != "" {
				normalizedImage := strings.ReplaceAll(tile.Image, "\\", "/")
				imgPath := path.Join(tsDir, normalizedImage)
				img, err := l.loadImage(imgPath)
				if err != nil {
					return nil, fmt.Errorf("failed to load image for tile %s: %w", imgPath, err)
				}
				imageMap[tile.Image] = img
			}
		}
	}
	return imageMap, nil
}

// convertLayers converts a slice of tiledLayer structs to a slice of MapLayer structs.
func convertLayers(tiledLayers []tiledLayer, tiles *map[int]Tile) ([]MapLayer, error) {
	gameLayers := make([]MapLayer, len(tiledLayers))
	for i, layerJSON := range tiledLayers {
		newLayer := MapLayer{
			Name:   layerJSON.Name,
			Type:   layerJSON.Type,
			Width:  layerJSON.Width,
			Height: layerJSON.Height,
		}

		switch newLayer.Type {
		case "tilelayer":
			newLayer.TileIds = layerJSON.Data
		case "objectgroup":
			objects, err := convertObjectGroup(layerJSON.Objects, tiles)
			if err != nil {
				return nil, err
			}
			newLayer.Objects = objects
		}
		gameLayers[i] = newLayer
	}
	return gameLayers, nil
}

// convertObjectGroup converts a slice of tiledObjects into a slice of Objects.
func convertObjectGroup(tiledObjects []tiledObject, tiles *map[int]Tile) ([]Object, error) {
	objects := make([]Object, len(tiledObjects))
	for i, objJSON := range tiledObjects {
		objType := ""
		objProperties := PropertySet{}

		// Look up the tile data if this object has a GID.
		yOffset := 0.0
		tileData, ok := (*tiles)[objJSON.GID]
		if ok {
			objType = tileData.Type
			objProperties = *tileData.Properties
			// For tile objects, adjust the Y position by the tile height.
			yOffset = objJSON.Height
		}

		// Merge object-specific properties, which override tile properties.
		if objJSON.Type != "" {
			objType = objJSON.Type
		}
		properties, err := GetProperties(objJSON.Properties)
		if err == nil {
			for k, v := range properties {
				objProperties[k] = v
			}
		}

		objects[i] = Object{
			Name:       objJSON.Name,
			Type:       objType,
			Properties: &objProperties,
			Location: Rect{
				X:      objJSON.X,
				Y:      objJSON.Y - yOffset,
				Width:  objJSON.Width,
				Height: objJSON.Height,
			},
			GID: objJSON.GID,
		}
	}
	return objects, nil
}

// loadFile is a helper method that reads a file from the embedded file system.
// It uses a cache to avoid reading the same file multiple times.
func (l *FsLoader) loadFile(path string) ([]byte, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check the cache first.
	if data, ok := l.cache[path]; ok {
		return data, nil
	}

	// Open the file from the embedded file system.
	f, err := l.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read all the data from the file.
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// Store the data in the cache for future use.
	l.cache[path] = data
	return data, nil
}

// loadImage is a helper method that reads an image from the embedded file system and uses the converter.
func (l *FsLoader) loadImage(path string) (ImageProvider, error) {
	l.imageMu.Lock()
	defer l.imageMu.Unlock()

	// Check the image cache first.
	if img, ok := l.imageCache[path]; ok {
		return img, nil
	}

	// Open the file from the embedded file system.
	f, err := l.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Decode the image. This will automatically detect the format.
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s: %w", path, err)
	}

	// Use the custom converter to process the image
	convertedImage, err := l.converter(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image '%s': %w", path, err)
	}

	// Store the decoded image in the cache.
	l.imageCache[path] = convertedImage
	return convertedImage, nil
}
