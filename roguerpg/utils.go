package main

import (
	"math"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func roundEven(n int) int {
	if n%2 == 0 {
		return n
	} else {
		return n - 1
	}
}

func sign(n int) int {
	if n > 0 {
		return 1
	} else if n < 0 {
		return -1
	}
	return 0
}

func clamp(n int, min int, max int) int {
	if n < min {
		return min
	} else if n > max {
		return max
	}
	return n
}

type Point struct {
	X int
	Y int
}

type FPoint struct {
	X float64
	Y float64
}

// ====================================================================
// Polygon utils
// ====================================================================

type Polygon []Point

func (poly Polygon) area() float64 {
	n := len(poly)
	if n < 3 {
		return 0.0
	}

	// Compute the area of the polygon using the shoelace algorithm
	sum := 0.0
	for i := range n {
		vCurr := poly[i]
		vNext := poly[(i+1)%n]
		sum += float64(vCurr.X*vNext.Y - vNext.X*vCurr.Y)
	}
	return math.Abs(sum) / 2.0
}

// onSegment checks if point q lies on the line segment pr.
func onSegment(p, q, r Point) bool {
	// Check Collinearity
	if orientation(p, q, r) != 0 {
		return false // Not collinear, so Q cannot be on segment PR
	}
	// Check Bounding Box
	return q.X <= max(p.X, r.X) && q.X >= min(p.X, r.X) &&
		q.Y <= max(p.Y, r.Y) && q.Y >= min(p.Y, r.Y)
}

// PointInPolygon implements the Ray Casting Algorithm (Even-Odd Rule).
// It determines if a given point is inside the polygon defined by the vertices.
func (poly Polygon) contains(point Point) bool {
	n := len(poly)
	if n < 3 {
		return false
	}

	x := float64(point.X)
	y := float64(point.Y)
	inside := false

	// Iterate over each edge of the polygon
	// (vj, vi) represents the current edge
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		vi := poly[i]
		vj := poly[j]

		if onSegment(Point{X: vi.X, Y: vi.Y}, point, Point{X: vj.X, Y: vj.Y}) {
			return true
		}

		// Check if the horizontal ray from (x, y) intersects the edge (vi, vj).
		// The condition (vi.Y > y) != (vj.Y > y) ensures the edge crosses the line y=Y.
		if (float64(vi.Y) > y) != (float64(vj.Y) > y) {
			// Calculate the X-intercept (intersectX) of the edge with the horizontal ray y=Y.
			vxDiff := float64(vj.X - vi.X)
			vyDiff := float64(vj.Y - vi.Y)
			pyDiff := float64(point.Y - vi.Y)
			intersectX := (vxDiff*pyDiff)/vyDiff + float64(vi.X)

			// Flip 'inside' if the intersection point is strictly to the right of our test point X.
			if x < intersectX {
				inside = !inside
			}
		}
	}

	return inside
}

// orientation finds the orientation of the ordered triplet (p, q, r).
func orientation(p, q, r Point) int {
	// The value is calculated as the signed area of the triangle (p, q, r)
	// (y2 - y1) * (x3 - x2) - (x2 - x1) * (y3 - y2)
	// where p=(x1, y1), q=(x2, y2), r=(x3, y3)
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)

	if val == 0 {
		return 0 // Collinear
	}
	if val > 0 {
		return 1 // Clockwise (or right turn)
	}
	return 2 // Counter-clockwise (or left turn)
}

// lineIntersection checks if the line segment p1q1 and p2q2 intersect.
func lineIntersection(p1, q1, p2, q2 Point) bool {
	// Find the four orientations needed for the general and special cases
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	// General case: segments intersect if orientations are different for both segments
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special Cases (Collinear checks)

	// p1, q1, p2 are collinear and p2 lies on segment p1q1
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	// p1, q1, q2 are collinear and q2 lies on segment p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	// p2, q2, p1 are collinear and p1 lies on segment p2q2
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}

	// p2, q2, q1 are collinear and q1 lies on segment p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}

func (poly1 Polygon) overlaps(poly2 Polygon) bool {
	// Check for Edge Intersections
	n1 := len(poly1)
	n2 := len(poly2)

	// Iterate over every edge of poly1 (P1-P2)
	for i := 0; i < n1; i++ {
		P1 := poly1[i]
		P2 := poly1[(i+1)%n1] // Wrap around to P0

		// Iterate over every edge of poly2 (Q1-Q2)
		for j := 0; j < n2; j++ {
			Q1 := poly2[j]
			Q2 := poly2[(j+1)%n2]

			if lineIntersection(P1, P2, Q1, Q2) {
				return true
			}
		}
	}

	// Check for Point-in-Polygon (Containment)

	// Check if any vertex of poly2 is inside poly1
	for _, vertex := range poly2 {
		if poly1.contains(vertex) {
			return true
		}
	}

	// Check if any vertex of poly1 is inside poly2
	// (This covers the case where poly1 is inside poly2)
	for _, vertex := range poly1 {
		if poly2.contains(vertex) {
			return true
		}
	}

	// If no edge intersects and no vertex is contained, they do not overlap.
	return false
}
