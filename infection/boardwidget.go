package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DragInfo struct {
	isDragging bool
	startLoc   int
	offsetX    int
	offsetY    int
	image      *ebiten.Image
}

func EmptyDragInfo() DragInfo {
	return DragInfo{
		isDragging: false,
		startLoc:   -1,
		offsetX:    0,
		offsetY:    0,
		image:      nil,
	}
}

func NewDragInfo(startLoc int, x, y int, image *ebiten.Image) DragInfo {
	return DragInfo{
		isDragging: true,
		startLoc:   startLoc,
		offsetX:    x,
		offsetY:    y,
		image:      image,
	}
}

type ComputerDragInfo struct {
	isDragging      bool
	totalFrameCount int
	currFrameCount  int
	image           *ebiten.Image
	move            Move
	// fromX           int
	// fromY           int
	// toX             int
	// toY             int
}

func NewComputerDragInfo() *ComputerDragInfo {
	return &ComputerDragInfo{
		isDragging:      false,
		totalFrameCount: 30,
		currFrameCount:  0,
		image:           nil,
		move:            Move{},
		// fromX:           0,
		// fromY:           0,
		// toX:             0,
		// toY:             0,
	}
}

func (g *BoardWidget) makeComputerDragInfo(m Move) *ComputerDragInfo {
	// TODO: handle the case where the move is a pass better
	image := WhiteSquare
	if g.gameBoard.playerToMove == Black {
		image = BlackSquare
	}

	return &ComputerDragInfo{
		isDragging:      true,
		totalFrameCount: 30,
		currFrameCount:  0,
		image:           image,
		move:            m,
		// fromX:           fromX*TileSize + g.bounds.Min.X,
		// fromY:           fromY*TileSize + g.bounds.Min.Y,
		// toX:             toX*TileSize + g.bounds.Min.X,
		// toY:             toY*TileSize + g.bounds.Min.Y,
	}
}

func (g *BoardWidget) UpdateComputerDragInfo() {
	d := g.computerDragInfo
	if !d.isDragging {
		return
	}
	d.currFrameCount++
	fmt.Printf("draginfo = %v\n", d)
	if d.currFrameCount >= d.totalFrameCount {
		d.isDragging = false
		// make the computer move
		fmt.Printf("BoardWidget.UpdateComputerDragInfo: making move = %v\n", d.move)
		g.gameBoard.Move(d.move)
		fmt.Println(g.gameBoard.String())
	}
}

func (g *BoardWidget) DrawComputerDragInfo(screen *ebiten.Image) {
	d := g.computerDragInfo
	if !d.isDragging {
		return
	}
	gameFromY, gameFromX := IndexToRowCol(d.move.from)
	gameToY, gameToX := IndexToRowCol(d.move.to)

	pixelFromX := float64(gameFromX*TileSize + g.bounds.Min.X)
	pixelFromY := float64(gameFromY*TileSize + g.bounds.Min.Y)
	pixelToX := float64(gameToX*TileSize + g.bounds.Min.X)
	pixelToY := float64(gameToY*TileSize + g.bounds.Min.Y)

	percentDone := float64(d.currFrameCount) / float64(d.totalFrameCount)
	x := (1-percentDone)*pixelFromX + percentDone*pixelToX
	y := (1-percentDone)*pixelFromY + percentDone*pixelToY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(d.image, op)
}

type BoardWidget struct {
	gameBoard *Board

	bounds           image.Rectangle
	dragInfo         DragInfo
	computerDragInfo *ComputerDragInfo
	allowUserInput   bool
}

func NewBoardWidget() *BoardWidget {
	margin := 40
	widget := BoardWidget{
		bounds: image.Rectangle{
			Min: image.Point{X: margin, Y: margin},
			Max: image.Point{X: margin + BoardSize*TileSize, Y: margin + BoardSize*TileSize},
		},
		gameBoard:        NewBoard(),
		dragInfo:         EmptyDragInfo(),
		computerDragInfo: NewComputerDragInfo(),
		allowUserInput:   false,
	}
	return &widget
}

func (g *BoardWidget) DoComputerMove(m Move) {
	g.computerDragInfo = g.makeComputerDragInfo(m)
	fmt.Printf("BoardWidget.DoComputerMove: draginfo = %v\n", g.computerDragInfo)
	// Update will update the drag position, which will
	// be drawn in the next Draw call.
	// At the end of the drag, the move will be applied.
}

func (g *BoardWidget) Draw(screen *ebiten.Image) {
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize; c++ {
			x := float64(c*TileSize + g.bounds.Min.X)
			y := float64(r*TileSize + g.bounds.Min.Y)

			backgroundImge := Empty1Square
			if (r+c)%2 == 0 {
				backgroundImge = Empty2Square
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			screen.DrawImage(backgroundImge, op)

			idx := GetIndex(r, c)
			if g.gameBoard.white.Get(idx) {
				screen.DrawImage(WhiteSquare, op)
			} else if g.gameBoard.black.Get(idx) {
				screen.DrawImage(BlackSquare, op)
			}
		}
	}

	if g.dragInfo.isDragging {
		op := &ebiten.DrawImageOptions{}
		x, y := ebiten.CursorPosition()
		op.GeoM.Translate(float64(x-g.dragInfo.offsetX), float64(y-g.dragInfo.offsetY))
		screen.DrawImage(g.dragInfo.image, op)
	}

	if g.computerDragInfo.isDragging {
		g.DrawComputerDragInfo(screen)
	}
}

func (g *BoardWidget) pointToIndex(x, y int) int {
	sqX := (x - g.bounds.Min.X) / TileSize
	sqY := (y - g.bounds.Min.Y) / TileSize
	if sqX < 0 || sqX >= BoardSize || sqY < 0 || sqY >= BoardSize {
		return -1
	}
	return GetIndex(sqY, sqX)
}

func (g *BoardWidget) Update() {

	g.UpdateComputerDragInfo()

	if !g.allowUserInput {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		index := g.pointToIndex(x, y)
		if index >= 0 {
			sqX := ((x-g.bounds.Min.X)/TileSize)*TileSize + g.bounds.Min.X
			sqY := ((y-g.bounds.Min.Y)/TileSize)*TileSize + g.bounds.Min.Y
			if g.gameBoard.playerToMove == White && g.gameBoard.white.Get(index) {
				g.dragInfo = NewDragInfo(index, x-sqX, y-sqY, WhiteSquare)
			} else if g.gameBoard.playerToMove == Black && g.gameBoard.black.Get(index) {
				g.dragInfo = NewDragInfo(index, x-sqX, y-sqY, BlackSquare)
			} else {
				g.dragInfo = EmptyDragInfo()
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.dragInfo.isDragging {

			startPos := g.dragInfo.startLoc
			endPos := g.pointToIndex(x, y)

			m, err := CreateMove(startPos, endPos)
			if err == nil {
				valid, _ := IsValidMove(g.gameBoard, m)
				if valid {
					g.gameBoard.Move(m)
				}
			}
		}
		g.dragInfo = EmptyDragInfo()
	}
}
