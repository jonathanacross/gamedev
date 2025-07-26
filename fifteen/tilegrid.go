package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"image/color"
)

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (r *Rect) Contains(x, y float64) bool {
	return r.X <= x && x <= r.X+r.Width && r.Y <= y && y <= r.Y+r.Height
}

var blackImage = getBlackImage()

type Tile struct {
	screenX float64
	screenY float64
	image   *ebiten.Image
}

type TileGrid struct {
	fifteen Fifteen

	background *ebiten.Image
	picture    *ebiten.Image
	tiles      [16]*Tile
}

func NewTileGrid(picture *ebiten.Image) *TileGrid {
	return &TileGrid{
		fifteen:    *NewFifteen(),
		background: Background,
		picture:    picture,
		tiles:      generateTiles(picture),
	}
}

func (g *TileGrid) updateTileLocations() {
	for i, v := range g.fifteen.Grid {
		if v == -1 {
			continue
		}
		g.tiles[v].screenX = float64(i%4) * TileSize
		g.tiles[v].screenY = float64(i/4) * TileSize
	}
}

func (g *TileGrid) Update() error {
	if !g.fifteen.IsSolved() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		imgX, imgY := ebiten.CursorPosition()

		x := (imgX - GridMargin) / TileSize
		y := (imgY - GridMargin) / TileSize

		err := g.fifteen.Move(x, y)
		if err == nil {
			// errors can/will occur if user clicks in an illegal place
			g.updateTileLocations()

			if g.fifteen.IsSolved() {
				PlayWinSound()
			}
		}
	}
	return nil
}

func (g *TileGrid) Draw(screen *ebiten.Image) {
	// Fill Background
	scale := float64(ScreenHeight) / float64(g.background.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	screen.DrawImage(g.background, op)

	if g.fifteen.IsSolved() {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(GridMargin, GridMargin)
		screen.DrawImage(g.picture, op)
	} else {
		// Deepen background for missing tile
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(GridMargin, GridMargin)
		screen.DrawImage(blackImage, op)

		// Fill Tiles
		for _, i := range g.fifteen.Grid {
			if i == -1 {
				continue
			}
			tile := g.tiles[i]
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tile.screenX+GridMargin, tile.screenY+GridMargin)
			screen.DrawImage(tile.image, op)
		}
	}
}

func (g *TileGrid) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func getBlackImage() *ebiten.Image {
	color := color.RGBA{R: 0, G: 0, B: 0, A: 128}
	blackImage := ebiten.NewImage(4*int(TileSize), 4*int(TileSize))
	blackImage.Fill(color)
	return blackImage
}

func (g *TileGrid) Randomize() {
	g.fifteen.Randomize()
	g.updateTileLocations()
}

// Generates 16 tiles, with numbers, from the source image.
// Assumes the source is square.
func generateTiles(picture *ebiten.Image) [16]*Tile {
	tiles := [16]*Tile{}

	strokewidth := 3

	srcTileSize := picture.Bounds().Dx() / 4
	scale := TileSize / float64(srcTileSize)

	for x := range 4 {
		for y := range 4 {
			srcRect := image.Rect(x*srcTileSize, y*srcTileSize, (x+1)*srcTileSize, (y+1)*srcTileSize)
			subImage := picture.SubImage(srcRect).(*ebiten.Image)
			tile := ebiten.NewImage(int(TileSize), int(TileSize))

			// scale subimage to fit tile
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			tile.DrawImage(subImage, op)

			// add border
			vector.StrokeRect(tile, 0, 0, float32(TileSize), float32(TileSize), float32(strokewidth), color.Black, false)

			// draw number
			tileNum := 4*y + x + 1
			dx := 52
			dy := 80
			if tileNum >= 10 {
				dx = 35
			}
			text.Draw(tile, fmt.Sprintf("%d", tileNum), NumberFont, dx-1, dy, color.White)
			text.Draw(tile, fmt.Sprintf("%d", tileNum), NumberFont, dx+1, dy, color.White)
			text.Draw(tile, fmt.Sprintf("%d", tileNum), NumberFont, dx, dy-1, color.White)
			text.Draw(tile, fmt.Sprintf("%d", tileNum), NumberFont, dx, dy+1, color.White)
			text.Draw(tile, fmt.Sprintf("%d", tileNum), NumberFont, dx, dy, color.Black)

			tiles[4*y+x] = &Tile{
				screenX: float64(x) * TileSize,
				screenY: float64(y) * TileSize,
				image:   tile,
			}
		}
	}

	return tiles
}
