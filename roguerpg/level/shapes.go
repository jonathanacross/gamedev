package level

import (
	"math"
	"math/rand/v2"
	"roguerpg/core"
)

// MakeShape generates a specialized polygon used for room generation.
func MakeShape(centerX, centerY int, xRadius, yRadius int) core.Polygon {
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
	fVertices := make([]core.FPoint, numSpokes)
	for i := range numSpokes {
		theta := dTheta * float64(i)
		r := radii[i]
		// Calculations remain in float64
		x := float64(centerX) + r*math.Cos(theta)
		y := float64(centerY) + r*math.Sin(theta)
		fVertices[i] = core.FPoint{X: x, Y: y}
	}

	// Calculate bounds
	xMin, xMax := fVertices[0].X, fVertices[0].X
	yMin, yMax := fVertices[0].Y, fVertices[0].Y
	for _, v := range fVertices {
		xMin = math.Min(xMin, v.X)
		xMax = math.Max(xMax, v.X)
		yMin = math.Min(yMin, v.Y)
		yMax = math.Max(yMax, v.Y)
	}

	// Check for zero-division risk
	var scaleX, scaleY float64 = 1.0, 1.0
	if xMax-xMin > 0 {
		scaleX = float64(2*xRadius+1) / (xMax - xMin)
	}
	if yMax-yMin > 0 {
		scaleY = float64(2*yRadius+1) / (yMax - yMin)
	}

	// Final conversion and scaling to integer Point
	vertices := make(core.Polygon, numSpokes)
	for i, v := range fVertices {
		// Scaling and offset calculation using floats
		// Note used core.Point which has X,Y int.
		scaledX := (v.X-xMin)*scaleX + float64(centerX-xRadius)
		scaledY := (v.Y-yMin)*scaleY + float64(centerY-yRadius)

		vertices[i].X = int(scaledX)
		vertices[i].Y = int(scaledY)
	}

	return vertices
}
