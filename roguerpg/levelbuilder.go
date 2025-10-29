package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

const (
	floor = 0
	wall  = 1
)

type LevelBlueprint struct {
	Width   int
	Height  int
	Squares [][]int
}

func (lb *LevelBlueprint) String() string {
	result := ""
	for y := range lb.Height {
		for x := range lb.Width {
			if lb.Squares[y][x] == floor {
				result = result + ". "
			} else {
				result = result + "# "
			}
		}
		result = result + "\n"
	}
	return result
}

func FillSolidWalls(width, height int) [][]int {
	squares := [][]int{}
	for range height {
		row := make([]int, width)
		for x := range width {
			row[x] = wall
		}
		squares = append(squares, row)
	}
	return squares
}

type RoomInfo struct {
	X        int
	Y        int
	XRadius  int
	YRadius  int
	Vertices Polygon
}

func (room *RoomInfo) Area() float64 {
	return room.Vertices.area()
}

func (room *RoomInfo) Overlaps(rooms []RoomInfo) bool {
	for _, prevRoom := range rooms {
		if room.Vertices.overlaps(prevRoom.Vertices) {
			return true
		}
	}
	return false
}

func genRoom(mapWidth, mapHeight int) RoomInfo {
	minRadius := 5
	maxRadius := 15
	xRadius := rand.IntN(maxRadius-minRadius) + minRadius
	yRadius := rand.IntN(maxRadius-minRadius) + minRadius
	x := rand.IntN(mapWidth-2*(xRadius+1)) + xRadius + 1
	y := rand.IntN(mapHeight-2*(yRadius+1)) + yRadius + 1

	vertices := MakeShape(x, y, xRadius, yRadius)

	return RoomInfo{
		X:        x,
		Y:        y,
		XRadius:  xRadius,
		YRadius:  yRadius,
		Vertices: vertices,
	}
}

func MakeShape(centerX, centerY int, xRadius, yRadius int) Polygon {
	r := (xRadius + yRadius) / 2
	numSpokes := 3 * r / 2
	// Set initial random radius for each spoke
	radii := make([]float64, numSpokes)
	for i := range numSpokes {
		radii[i] = float64(r/2) + float64(r)*rand.Float64()
	}

	// Smooth the radii
	smoothedRadii := make([]float64, numSpokes)
	for i := range numSpokes {
		prev := radii[(i-1+numSpokes)%numSpokes]
		curr := radii[i]
		next := radii[(i+1)%numSpokes]
		smoothedRadii[i] = 0.15*prev + 0.7*curr + 0.15*next
	}
	radii = smoothedRadii

	// Use FPoint for vertex positions calculation
	dTheta := 2 * math.Pi / float64(numSpokes)
	fVertices := make([]FPoint, numSpokes) // Changed to FPoint
	for i := range numSpokes {
		theta := dTheta * float64(i)
		r := radii[i]
		// Calculations remain in float64
		x := float64(centerX) + r*math.Cos(theta)
		y := float64(centerY) + r*math.Sin(theta)
		fVertices[i] = FPoint{X: x, Y: y}
	}

	// Calculate bounds using FPoint values
	xMin, xMax := fVertices[0].X, fVertices[0].X
	yMin, yMax := fVertices[0].Y, fVertices[0].Y
	for _, v := range fVertices {
		xMin = math.Min(xMin, v.X)
		xMax = math.Max(xMax, v.X)
		yMin = math.Min(yMin, v.Y)
		yMax = math.Max(yMax, v.Y)
	}

	// Check for zero-division risk (in case shape collapsed to a line/point)
	var scaleX, scaleY float64 = 1.0, 1.0
	if xMax-xMin > 0 {
		scaleX = float64(2*xRadius+1) / (xMax - xMin)
	}
	if yMax-yMin > 0 {
		scaleY = float64(2*yRadius+1) / (yMax - yMin)
	}

	// Final conversion and scaling to integer Point
	vertices := make(Polygon, numSpokes)
	for i, v := range fVertices {
		// Scaling and offset calculation using floats
		scaledX := (v.X-xMin)*scaleX + float64(centerX-xRadius)
		scaledY := (v.Y-yMin)*scaleY + float64(centerY-yRadius)

		// Final conversion to integer Point
		vertices[i].X = int(scaledX)
		vertices[i].Y = int(scaledY)
	}

	return vertices
}

func BuildLevel(width, height int) *LevelBlueprint {
	// initialize to solid wall
	data := FillSolidWalls(width, height)

	// pick rooms until we get some percent of the area
	rooms := []RoomInfo{}
	area := 0.0
	attempts := 0
	for area < 0.50*float64(width*height) && attempts < 1000 {
		attempts++
		room := genRoom(width, height)
		if !room.Overlaps(rooms) {
			rooms = append(rooms, room)
			area += room.Area()
		}
	}

	for _, room := range rooms {
		fmt.Printf("room = %v\n", room)
	}

	// clear out the area in the rooms
	for y := range height {
		for x := range width {
			point := Point{X: x, Y: y}
			for _, room := range rooms {
				if room.Vertices.contains(point) {
					data[y][x] = floor
					break
				}
			}
		}
	}

	return &LevelBlueprint{
		Width:   width,
		Height:  height,
		Squares: data,
	}
}

func main() {
	lb := BuildLevel(70, 50)
	fmt.Printf("%v", lb)
}
