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

type Point struct {
	X int
	Y int
}

// PointInPolygon implements the Ray Casting Algorithm (Even-Odd Rule).
// It determines if a given point is inside the polygon defined by the vertices.
func PointInPolygon(point Point, vertices []Point) bool {
	n := len(vertices)
	if n < 3 {
		return false
	}

	x := float64(point.X)
	//y := float64(point.Y)
	inside := false

	// Iterate over each edge of the polygon
	// (vj, vi) represents the current edge
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		vi := vertices[i]
		vj := vertices[j]

		// Check if the horizontal ray from (x, y) intersects the edge (vi, vj).
		// The condition (vi.Y > y) != (vj.Y > y) ensures the edge crosses the line y=Y.
		if (vi.Y > point.Y) != (vj.Y > point.Y) {
			// Calculate the X-intercept (intersectX) of the edge with the horizontal ray y=Y.
			intersectX := (float64(vj.X-vi.X)*float64(point.Y-vi.Y))/(float64(vj.Y-vi.Y)) + float64(vi.X)

			// If the X-intercept is to the right of the test point's X (x < intersectX),
			// the ray crosses the edge, so we flip the 'inside' state.
			if x < intersectX {
				inside = !inside
			}
		}
	}

	return inside
}

type RoomInfo struct {
	X        int
	Y        int
	XRadius  int
	YRadius  int
	Vertices []Point
}

func (room *RoomInfo) Area() float64 {
	// TODO: replace with polygon area
	return float64(room.XRadius*2+1) * float64(room.YRadius*2+1)
}

func (room *RoomInfo) Overlaps(rooms []RoomInfo) bool {
	// TODO: replace with polygon intersection
	for _, prevRoom := range rooms {
		dx := abs(room.X - prevRoom.X)
		dy := abs(room.Y - prevRoom.Y)
		minDistX := room.XRadius + prevRoom.XRadius + 1
		minDistY := room.XRadius + prevRoom.YRadius + 1
		if dx <= minDistX && dy <= minDistY {
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

func MakeShape(centerX, centerY int, xRadius, yRadius int) []Point {
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

	// compute vertex positions for the shape
	dTheta := 2 * math.Pi / float64(numSpokes)
	vertices := make([]Point, numSpokes)
	for i := range numSpokes {
		theta := dTheta * float64(i)
		r := radii[i]
		x := centerX + int(r*math.Cos(theta))
		y := centerY + int(r*math.Sin(theta))
		vertices[i] = Point{X: x, Y: y}
	}

	xMin, xMax := vertices[0].X, vertices[0].X
	yMin, yMax := vertices[0].Y, vertices[0].Y
	for _, v := range vertices {
		xMin = min(xMin, v.X)
		xMax = max(xMax, v.X)
		yMin = min(yMin, v.Y)
		yMax = max(yMax, v.Y)
	}
	scaleX := float64(2*xRadius+1) / float64(xMax-xMin)
	scaleY := float64(2*yRadius+1) / float64(yMax-yMin)

	for i, v := range vertices {
		vertices[i].X = int(float64(v.X-xMin)*scaleX + float64(centerX-xRadius))
		vertices[i].Y = int(float64(v.Y-yMin)*scaleY + float64(centerY-yRadius))
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
				if PointInPolygon(point, room.Vertices) {
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
