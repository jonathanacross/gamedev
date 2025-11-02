package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strconv"
)

// This tool constructs a "blob" tilemap image from a sub-blob tilemap image.
// For background, see
// https://www.boristhebrave.com/2013/07/14/tileset-roundup/.
//
// The sub-blob tileset must be a 12x2 tile sheet in the following layout:
// +--+--+--+--+--+--+--+--+--+--+--+--+
// |##|##|.#|#.|..|..|.#|#.|..|..|..|..|
// |##|##|##|##|##|##|.#|#.|.#|#.|..|..|
// +--+--+--+--+--+--+--+--+--+--+--+--+
// |##|##|##|##|##|##|.#|#.|.#|#.|..|..|
// |##|##|.#|#.|..|..|.#|#.|..|..|..|..|
// +--+--+--+--+--+--+--+--+--+--+--+--+
//
// The resulting blob tilemap will be 7x7 tiles, where each tile
// is made up of 4 sub-blob tiles arranged in a 2x2 grid.
// The layout of the blob tileset is from
// https://www.boristhebrave.com/permanent/24/06/cr31/stagecast/art/atlas/blob/wangbl.png
// Numbers correspond to bits being enabled in the following directions:
//
//	 7  0  1
//	  \ | /
//	6 - * - 2
//	  / | \
//	 5  4  3

func readImage(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error opening image file: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("Error decoding image: %v", err)
	}

	return img, nil
}

type SubBlobTileMap struct {
	sourceImage image.Image
	tileSize    int
	tiles       []image.Rectangle
}

type Direction int

const (
	NW Direction = iota // 0
	NE                  // 1
	SW                  // 2
	SE                  // 3
)

// getSubBlobKey computes the unique key for the subBlobLookup map.
// The key is composed of (Pattern << 2) | Direction.
func getSubBlobKey(pattern int, dir Direction) int {
	return (pattern << 2) | int(dir)
}

var subBlobLookup = map[int]int{
	getSubBlobKey(0, NW):  10,
	getSubBlobKey(0, NE):  11,
	getSubBlobKey(0, SW):  22,
	getSubBlobKey(0, SE):  23,
	getSubBlobKey(1, NW):  8,
	getSubBlobKey(2, NE):  9,
	getSubBlobKey(3, NW):  4,
	getSubBlobKey(3, NE):  5,
	getSubBlobKey(4, SW):  20,
	getSubBlobKey(5, NW):  6,
	getSubBlobKey(5, SW):  18,
	getSubBlobKey(7, NW):  2,
	getSubBlobKey(8, SE):  21,
	getSubBlobKey(10, NE): 7,
	getSubBlobKey(10, SE): 19,
	getSubBlobKey(11, NE): 3,
	getSubBlobKey(12, SW): 16,
	getSubBlobKey(12, SE): 17,
	getSubBlobKey(13, SW): 14,
	getSubBlobKey(14, SE): 15,
	getSubBlobKey(15, NW): 0,
	getSubBlobKey(15, NE): 1,
	getSubBlobKey(15, SW): 12,
	getSubBlobKey(15, SE): 13,
}

// blobIDToGrid converts a blob ID into a 3x3 grid of 'on' (1) or 'off' (0) values.
func blobIDToGrid(blobID int) [9]int {
	var grid [9]int

	// Bit 0:N, 1:NE, 2:E, 3:SE, 4:S, 5:SW, 6:W, 7:NW
	// Grid Index: 0:NW, 1:N, 2:NE, 3:W, 4:C, 5:E, 6:SW, 7:S, 8:SE
	bitToGridIndex := [8]int{1, 2, 5, 8, 7, 6, 3, 0}

	for bit := 0; bit < 8; bit++ {
		if (blobID>>bit)&1 == 1 {
			grid[bitToGridIndex[bit]] = 1
		}
	}

	if blobID != 0 {
		grid[4] = 1 // Center is 'on'
	}

	return grid
}

// subgridToSubBlobID converts a 2x2 sub-square pattern and its direction to a sub-blob ID.
func subgridToSubBlobID(tl, tr, bl, br int, dir Direction) int {
	pattern := (tl << 3) | (tr << 2) | (bl << 1) | br
	key := getSubBlobKey(pattern, dir)

	if id, ok := subBlobLookup[key]; ok {
		return id
	}

	return -1 // Should not be reached
}

// convertBlobIDToSubBlobs converts a blob ID into 4 sub-blob IDs: [NW, NE, SW, SE].
func convertBlobIDToSubBlobs(blobID int) [4]int {
	grid := blobIDToGrid(blobID)

	nwID := subgridToSubBlobID(grid[0], grid[1], grid[3], grid[4], NW)
	neID := subgridToSubBlobID(grid[1], grid[2], grid[4], grid[5], NE)
	swID := subgridToSubBlobID(grid[3], grid[4], grid[6], grid[7], SW)
	seID := subgridToSubBlobID(grid[4], grid[5], grid[7], grid[8], SE)

	return [4]int{nwID, neID, swID, seID}
}

func getRect(idx int, tileSize int, widthInTiles int) image.Rectangle {
	x := (idx % widthInTiles) * tileSize
	y := (idx / widthInTiles) * tileSize

	return image.Rectangle{Min: image.Point{X: x, Y: y}, Max: image.Point{X: x + tileSize, Y: y + tileSize}}
}

func NewSubBlobTileMap(img image.Image, tileSize int) *SubBlobTileMap {
	widthInTiles := 12
	heightInTiles := 2
	numTiles := widthInTiles * heightInTiles
	tiles := make([]image.Rectangle, numTiles)
	for i := range 24 {
		tiles[i] = getRect(i, tileSize, widthInTiles)
	}

	return &SubBlobTileMap{
		sourceImage: img,
		tileSize:    tileSize,
		tiles:       tiles,
	}
}

func (sb *SubBlobTileMap) BuildBlobImage() image.Image {
	layout := [][]int{
		{0, 4, 92, 124, 116, 80, 0},
		{16, 20, 87, 223, 241, 21, 64},
		{29, 117, 85, 71, 221, 125, 112},
		{31, 253, 113, 28, 127, 247, 209},
		{23, 199, 213, 95, 255, 245, 81},
		{5, 84, 93, 119, 215, 193, 17},
		{0, 1, 7, 197, 69, 68, 65},
	}

	outputImage := image.NewRGBA(image.Rect(0, 0, 2*len(layout[0])*sb.tileSize, 2*len(layout)*sb.tileSize))

	for rowIdx, row := range layout {
		for colIdx, blobID := range row {
			subTileIndices := convertBlobIDToSubBlobs(blobID)

			destX := colIdx * 2 * sb.tileSize
			destY := rowIdx * 2 * sb.tileSize

			offsets := []image.Point{
				{X: 0, Y: 0},
				{X: sb.tileSize, Y: 0},
				{X: 0, Y: sb.tileSize},
				{X: sb.tileSize, Y: sb.tileSize},
			}

			for i := 0; i < 4; i++ {
				srcRect := sb.tiles[subTileIndices[i]]
				destBounds := image.Rect(destX+offsets[i].X, destY+offsets[i].Y, destX+offsets[i].X+sb.tileSize, destY+offsets[i].Y+sb.tileSize)
				fmt.Printf("Placing Blob ID %d Sub-tile %d at dest %v from src %v\n", blobID, subTileIndices[i], destBounds, srcRect)

				draw.Draw(outputImage, destBounds, sb.sourceImage, srcRect.Min, draw.Src)
			}
		}
	}
	return outputImage
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: blob_builder [input_image.png] [input_tile_size] [output_image.png]")
		return
	}
	inputFileName := os.Args[1]
	tileSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Error parsing tile size: %v", err)
	}
	outputFileName := os.Args[3]

	img, err := readImage(inputFileName)
	if err != nil {
		log.Fatalf("Error reading image: %v", err)
	}

	outputImage := NewSubBlobTileMap(img, tileSize).BuildBlobImage()
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, outputImage)
	if err != nil {
		log.Fatalf("Error encoding output image: %v", err)
	}

}
