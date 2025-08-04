package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DragInfo struct {
	isDragging bool
	startLoc   SquareIndex
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

func NewDragInfo(startLoc SquareIndex, x, y int, image *ebiten.Image) DragInfo {
	return DragInfo{
		isDragging: true,
		startLoc:   startLoc,
		offsetX:    x,
		offsetY:    y,
		image:      image,
	}
}

type ComputerDragInfo struct {
	isAnimating     bool
	totalFrameCount int
	currFrameCount  int
	image           *ebiten.Image
	move            Move
}

func NewComputerDragInfo() *ComputerDragInfo {
	return &ComputerDragInfo{
		isAnimating:     false,
		totalFrameCount: 30,
		currFrameCount:  0,
		image:           nil,
		move:            Move{},
	}
}

func (g *BoardWidget) makeComputerDragInfo(m Move, playerToMove Player) {
	image := WhiteSquare
	if playerToMove == Black {
		image = BlackSquare
	}
	g.computerDragInfo.isAnimating = true
	g.computerDragInfo.currFrameCount = 0
	g.computerDragInfo.image = image
	g.computerDragInfo.move = m
}

func (g *BoardWidget) UpdateComputerDragInfo() {
	d := g.computerDragInfo
	if !d.isAnimating {
		return
	}
	d.currFrameCount++
}

func (g *BoardWidget) DrawComputerDragInfo(screen *ebiten.Image, gameBoard *Board) {
	d := g.computerDragInfo
	if !d.isAnimating {
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
	bounds           image.Rectangle
	dragInfo         DragInfo
	computerDragInfo *ComputerDragInfo
	humanMove        *Move
}

func NewBoardWidget() *BoardWidget {
	margin := 40
	widget := BoardWidget{
		bounds: image.Rectangle{
			Min: image.Point{X: margin, Y: margin},
			Max: image.Point{X: margin + BoardSize*TileSize, Y: margin + BoardSize*TileSize},
		},
		dragInfo:         EmptyDragInfo(),
		computerDragInfo: NewComputerDragInfo(),
		humanMove:        nil,
	}
	return &widget
}

func (g *BoardWidget) DoComputerMove(m Move, playerToMove Player) {
	g.makeComputerDragInfo(m, playerToMove)
}

func (g *BoardWidget) Draw(screen *ebiten.Image, gameBoard *Board) {
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
			if gameBoard.white.Get(idx) {
				screen.DrawImage(WhiteSquare, op)
			} else if gameBoard.black.Get(idx) {
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

	if g.computerDragInfo.isAnimating {
		g.DrawComputerDragInfo(screen, gameBoard)
	}
}

func (g *BoardWidget) pointToIndex(x, y int) SquareIndex {
	sqX := (x - g.bounds.Min.X) / TileSize
	sqY := (y - g.bounds.Min.Y) / TileSize
	if sqX < 0 || sqX >= BoardSize || sqY < 0 || sqY >= BoardSize {
		return -1
	}
	return GetIndex(sqY, sqX)
}

func (g *BoardWidget) Update(gameBoard *Board) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		index := g.pointToIndex(x, y)
		if index >= 0 {
			sqX := ((x-g.bounds.Min.X)/TileSize)*TileSize + g.bounds.Min.X
			sqY := ((y-g.bounds.Min.Y)/TileSize)*TileSize + g.bounds.Min.Y
			if gameBoard.playerToMove == White && gameBoard.white.Get(index) {
				g.dragInfo = NewDragInfo(index, x-sqX, y-sqY, WhiteSquare)
			} else if gameBoard.playerToMove == Black && gameBoard.black.Get(index) {
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
			if endPos != -1 {
				m, err := CreateMove(startPos, endPos)
				if err == nil {
					g.humanMove = &m
				}
			}
		}
		g.dragInfo = EmptyDragInfo()
	}
}

func (g *BoardWidget) GetAndClearHumanMove() (Move, bool) {
	if g.humanMove != nil {
		move := *g.humanMove
		g.humanMove = nil
		return move, true
	}
	return Move{}, false
}

func (g *BoardWidget) GetAndClearComputerMove() (Move, bool) {
	if g.computerDragInfo.isAnimating && g.computerDragInfo.currFrameCount >= g.computerDragInfo.totalFrameCount {
		move := g.computerDragInfo.move
		g.computerDragInfo.isAnimating = false
		return move, true
	}
	return Move{}, false
}
