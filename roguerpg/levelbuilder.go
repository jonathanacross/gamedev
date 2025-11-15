package main

import (
	"math"
	"math/rand/v2"
)

const (
	floor    = 0
	wall     = 1
	vertical = 2
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
			switch lb.Squares[y][x] {
			case floor:
				result = result + ". "
			case wall:
				result = result + "# "
			case vertical:
				result = result + "| "
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
	margin := 2
	x := rand.IntN(mapWidth-2*(xRadius+margin)) + xRadius + margin
	y := rand.IntN(mapHeight-2*(yRadius+margin)) + yRadius + margin

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

func visitPermutation(rooms []RoomInfo) []int {
	permutation := []int{}
	seen := make([]bool, len(rooms))

	// start with first room
	currRoom := 0
	seen[currRoom] = true
	permutation = append(permutation, currRoom)

	for len(permutation) < len(rooms) {
		// find nearest unseen room
		bestDist := math.MaxInt
		closestRoom := -1
		for i, r := range rooms {
			if !seen[i] {
				dist := abs(r.X-rooms[currRoom].X) + abs(r.Y-rooms[currRoom].Y)
				if dist < bestDist {
					bestDist = dist
					closestRoom = i
				}
			}
		}

		currRoom = closestRoom
		seen[currRoom] = true
		permutation = append(permutation, currRoom)
	}

	return permutation
}

// adds a new path
func addPath(data [][]int, rooms []RoomInfo, oldx, oldy, newx, newy int) {
	dx := sign(newx - oldx)
	dy := sign(newy - oldy)

	for x := oldx; x != newx; x += dx {
		data[oldy][x] = floor
	}
	for y := oldy; y != newy; y += dy {
		data[y][newx] = floor
	}
}

func connectRoomsWithPaths(data [][]int, rooms []RoomInfo) {
	// Sort the rooms by closest path so there aren't
	// too many long paths that crisscross over the
	// dungeon.
	pi := visitPermutation(rooms)

	// Add a few extra paths so that the final result isn't
	// too circular.
	pi = append(pi, pi[0])
	pi = append(pi, pi[len(rooms)/4])
	pi = append(pi, pi[2*len(rooms)/4])
	pi = append(pi, pi[3*len(rooms)/4])

	// connect the rooms with corridors -- just go in order
	for i := 1; i < len(pi); i++ {
		oldx := roundEven(rooms[pi[i-1]].X)
		oldy := roundEven(rooms[pi[i-1]].Y)
		newx := roundEven(rooms[pi[i]].X)
		newy := roundEven(rooms[pi[i]].Y)

		addPath(data, rooms, oldx, oldy, newx, newy)
	}
}

// Remove patterns that don't draw sensibly in 2.5D.
func removeIllegalPatterns(data [][]int) {
	//	# .   and   . #
	//	. #         # .
	for y := 1; y < len(data)-1; y++ {
		for x := 1; x < len(data[0])-1; x++ {
			d11 := data[y][x]
			d12 := data[y+1][x]
			d21 := data[y][x+1]
			d22 := data[y+1][x+1]

			if d11 == wall && d12 == floor && d21 == floor && d22 == wall {
				data[y][x] = floor
			} else if d11 == floor && d12 == wall && d21 == wall && d22 == floor {
				data[y][x+1] = floor
			}
		}
	}

	//	.
	//	#
	//	.
	for y := 1; y < len(data)-2; y++ {
		for x := 1; x < len(data[0])-1; x++ {
			if data[y][x] == floor && data[y+1][x] == wall && data[y+2][x] == floor {
				data[y+1][x] = floor
			}
		}
	}
}

func ensureOuterEdgeIsWall(data [][]int) {
	for x := range data[0] {
		data[0][x] = wall
		data[len(data)-1][x] = wall
	}
	for y := range data {
		data[y][0] = wall
		data[y][len(data[0])-1] = wall
	}
}

func extendVerticalWalls(data [][]int) {
	for y := 0; y < len(data)-1; y++ {
		for x := 0; x < len(data[0]); x++ {
			if data[y][x] == wall && data[y+1][x] == floor {
				data[y][x] = vertical
			}
		}
	}
}

func BuildLevelBlueprint(width, height int) *LevelBlueprint {
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

	connectRoomsWithPaths(data, rooms)

	removeIllegalPatterns(data)

	ensureOuterEdgeIsWall(data)

	extendVerticalWalls(data)

	return &LevelBlueprint{
		Width:   width,
		Height:  height,
		Squares: data,
	}
}

func (lb *LevelBlueprint) GetSheetIdxForTerrain(x, y int) int {
	switch lb.Squares[y][x] {
	case floor:
		return rand.IntN(5) + 5
	case wall:
		return 10
	case vertical:
		return rand.IntN(3) + 0
	default:
		return 0
	}
}

// returns 1 if wall, 0 otherwise
func (lb *LevelBlueprint) getWallCoeff(x, y int) int {
	if (x < 0 || x >= lb.Width || y < 0 || y >= lb.Height) ||
		lb.Squares[y][x] == wall {
		return 1
	}
	return 0
}

type Offset struct {
	X     int
	Y     int
	value int
}

func (lb *LevelBlueprint) GetOffsetValue(x, y int, offset Offset) int {
	nx := x + offset.X
	ny := y + offset.Y
	coeff := lb.getWallCoeff(nx, ny)
	return coeff * offset.value
}

func (lb *LevelBlueprint) GetBlobValueForWallEdge(x, y int) int {
	N := Offset{0, -1, 1 << 0}
	NE := Offset{1, -1, 1 << 1}
	E := Offset{1, 0, 1 << 2}
	SE := Offset{1, 1, 1 << 3}
	S := Offset{0, 1, 1 << 4}
	SW := Offset{-1, 1, 1 << 5}
	W := Offset{-1, 0, 1 << 6}
	NW := Offset{-1, -1, 1 << 7}

	cardinals := []Offset{N, E, S, W}

	lookup := 0
	for _, offset := range cardinals {
		lookup += lb.GetOffsetValue(x, y, offset)
	}

	if (lookup&N.value > 0) && (lookup&E.value > 0) {
		lookup += lb.GetOffsetValue(x, y, NE)
	}
	if (lookup&S.value > 0) && (lookup&E.value > 0) {
		lookup += lb.GetOffsetValue(x, y, SE)
	}
	if (lookup&N.value > 0) && (lookup&W.value > 0) {
		lookup += lb.GetOffsetValue(x, y, NW)
	}
	if (lookup&S.value > 0) && (lookup&W.value > 0) {
		lookup += lb.GetOffsetValue(x, y, SW)
	}

	return lookup
}

func BuildLevel(width, height int) [][]*Tile {
	tileSize := 16
	blueprint := BuildLevelBlueprint(width, height)
	terrainTiles := [][]*Tile{}

	terrainSpriteSheet := NewSpriteSheet(tileSize, tileSize, 5, 3)
	wallEdgeBlobSpriteSheet := NewBlobSpriteSheet(tileSize, tileSize)

	for y := 0; y < blueprint.Height; y++ {
		row := []*Tile{}
		for x := 0; x < blueprint.Width; x++ {
			squareType := blueprint.Squares[y][x]
			location := Location{
				X: float64(x * tileSize),
				Y: float64(y * tileSize),
			}

			var tile *Tile = nil
			if squareType == wall {
				blobValue := blueprint.GetBlobValueForWallEdge(x, y)
				tile = NewTile(location, WallBlobTileset, wallEdgeBlobSpriteSheet.Rect(blobValue), true)
			} else {
				tileIdx := blueprint.GetSheetIdxForTerrain(x, y)
				solid := blueprint.Squares[y][x] != floor
				tile = NewTile(location, TerrainTileset, terrainSpriteSheet.Rect(tileIdx), solid)
			}
			row = append(row, tile)
		}
		terrainTiles = append(terrainTiles, row)
	}

	return terrainTiles
}

/*
func main() {
	lb := BuildLevel(70, 50)
	fmt.Printf("%v", lb)
}
*/
