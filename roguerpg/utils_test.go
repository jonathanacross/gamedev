package main

import (
	"math"
	"testing"
)

// Helper function for float comparison in tests due to precision issues
const floatTolerance = 1e-9

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

// ====================================================================
// Basic Helpers Tests
// ====================================================================

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for _, tt := range tests {
		if got := abs(tt.input); got != tt.want {
			t.Errorf("abs(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMin(t *testing.T) {
	if got := min(5, 10); got != 5 {
		t.Errorf("min(5, 10) = %v, want 5", got)
	}
	if got := min(10, 5); got != 5 {
		t.Errorf("min(10, 5) = %v, want 5", got)
	}
	if got := min(-1, -10); got != -10 {
		t.Errorf("min(-1, -10) = %v, want -10", got)
	}
}

func TestMax(t *testing.T) {
	if got := max(5, 10); got != 10 {
		t.Errorf("max(5, 10) = %v, want 10", got)
	}
	if got := max(10, 5); got != 10 {
		t.Errorf("max(10, 5) = %v, want 10", got)
	}
	if got := max(-1, -10); got != -1 {
		t.Errorf("max(-1, -10) = %v, want -1", got)
	}
}

// ====================================================================
// Polygon Tests
// ====================================================================

func TestArea(t *testing.T) {
	tests := []struct {
		name string
		poly Polygon
		want float64
	}{
		{
			name: "Triangle",
			poly: Polygon{{1, 1}, {3, 4}, {5, 1}}, // Base 4, Height 3 -> Area 6
			want: 6.0,
		},
		{
			name: "Rectangle",
			poly: Polygon{{0, 0}, {10, 0}, {10, 5}, {0, 5}}, // Area 50
			want: 50.0,
		},
		{
			name: "Self-Crossing (Should still calculate signed area)",
			poly: Polygon{{0, 0}, {10, 10}, {0, 10}, {10, 0}}, // Area 0 due to crossing
			want: 0.0,
		},
		{
			name: "Small Polygon",
			poly: Polygon{{1, 2}, {2, 1}, {3, 2}, {2, 3}}, // Diamond shape, Area 2
			want: 2.0,
		},
		{
			name: "Less than 3 vertices",
			poly: Polygon{{1, 1}, {2, 2}},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := area(tt.poly); !floatEquals(got, tt.want) {
				t.Errorf("area() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrientation(t *testing.T) {
	tests := []struct {
		name    string
		p, q, r Point
		want    int // 0: Collinear, 1: CW, 2: CCW
	}{
		{"Collinear", Point{0, 0}, Point{1, 1}, Point{2, 2}, 0},
		{"Collinear (Reversed)", Point{2, 2}, Point{1, 1}, Point{0, 0}, 0},
		{"Clockwise", Point{0, 0}, Point{1, 0}, Point{0, 1}, 2},
		{"CounterClockwise", Point{0, 0}, Point{0, 1}, Point{1, 0}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := orientation(tt.p, tt.q, tt.r); got != tt.want {
				t.Errorf("orientation(%v, %v, %v) = %v, want %v", tt.p, tt.q, tt.r, got, tt.want)
			}
		})
	}
}

func TestOnSegment(t *testing.T) {
	p := Point{0, 0}
	r := Point{10, 10}

	tests := []struct {
		name string
		q    Point
		want bool
	}{
		{"Midpoint", Point{5, 5}, true},
		{"Endpoint P", p, true},
		{"Endpoint R", r, true},
		{"Outside X-range", Point{11, 11}, false},
		{"Outside Y-range", Point{-1, -1}, false},
		{"Off Segment", Point{5, 6}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := onSegment(p, tt.q, r); got != tt.want {
				t.Errorf("onSegment(%v, %v, %v) = %v, want %v", p, tt.q, r, got, tt.want)
			}
		})
	}
}

func TestLineIntersection(t *testing.T) {
	tests := []struct {
		name           string
		p1, q1, p2, q2 Point // Segment 1: p1-q1, Segment 2: p2-q2
		want           bool
	}{
		{"General Case (Cross)", Point{10, 0}, Point{0, 10}, Point{0, 0}, Point{10, 10}, true},
		{"No Intersection", Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}, false},
		{"Shared Endpoint", Point{0, 0}, Point{10, 0}, Point{10, 0}, Point{20, 0}, true},
		{"T-Junction", Point{0, 0}, Point{10, 0}, Point{5, -5}, Point{5, 5}, true},
		{"Collinear Overlap", Point{0, 0}, Point{10, 0}, Point{5, 0}, Point{15, 0}, true},
		{"Collinear No Overlap", Point{0, 0}, Point{5, 0}, Point{10, 0}, Point{15, 0}, false},
		{"Collinear Just Touch", Point{0, 0}, Point{5, 0}, Point{5, 0}, Point{10, 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lineIntersection(tt.p1, tt.q1, tt.p2, tt.q2); got != tt.want {
				t.Errorf("lineIntersection(%v-%v, %v-%v) = %v, want %v", tt.p1, tt.q1, tt.p2, tt.q2, got, tt.want)
			}
		})
	}
}

func TestPointInPolygon(t *testing.T) {
	// A simple square: (1,1), (3,1), (3,3), (1,3)
	square := Polygon{{1, 1}, {3, 1}, {3, 3}, {1, 3}}

	// A concave polygon: (0,0), (4,0), (4,4), (2,2), (0,4)
	concave := Polygon{{0, 0}, {4, 0}, {4, 4}, {2, 2}, {0, 4}}

	tests := []struct {
		name  string
		point Point
		poly  Polygon
		want  bool
	}{
		// Square tests
		{"Square: Inside Center", Point{2, 2}, square, true},
		{"Square: Outside", Point{0, 0}, square, false},
		{"Square: On Vertex (Ray passes through)", Point{1, 1}, square, true},
		{"Square: On Edge (Horizontal)", Point{2, 1}, square, true},

		// Concave tests
		{"Concave: Inside", Point{1, 1}, concave, true},
		{"Concave: Inside Concave Area", Point{3, 3}, concave, true},
		{"Concave: On Boundary Concave Area", Point{3, 3}, concave, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pointInPolygon(tt.point, tt.poly); got != tt.want {
				t.Errorf("pointInPolygon(%v) = %v, want %v", tt.point, got, tt.want)
			}
		})
	}
}

func TestOverlaps(t *testing.T) {
	// Simple square
	polyA := Polygon{{0, 0}, {4, 0}, {4, 4}, {0, 4}}
	// Simple square
	polyB := Polygon{{3, 3}, {7, 3}, {7, 7}, {3, 7}}
	// Simple square
	polyC := Polygon{{5, 5}, {9, 5}, {9, 9}, {5, 9}}
	// Internal square D inside A
	polyD := Polygon{{1, 1}, {3, 1}, {3, 3}, {1, 3}}

	// Concave poly E
	polyE := Polygon{{0, 0}, {10, 0}, {10, 10}, {5, 5}, {0, 10}}
	// Small poly F that intersects the concave part of E
	polyF := Polygon{{4, 6}, {6, 6}, {6, 8}, {4, 8}}

	tests := []struct {
		name         string
		poly1, poly2 Polygon
		want         bool
	}{
		{"Overlap (Cross)", polyA, polyB, true},
		{"No Overlap", polyA, polyC, false},
		{"Containment (D in A)", polyA, polyD, true},
		{"Containment (A around D)", polyD, polyA, true},                          // Should check vertices of D in A
		{"Touch at Vertex", polyA, Polygon{{4, 4}, {8, 4}, {8, 8}, {4, 8}}, true}, // Vertex touch (4,4)
		{"No Overlap (Close)", Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}}, Polygon{{2, 2}, {3, 2}, {3, 3}, {2, 3}}, false},
		{"Concave Intersection (Edge)", polyE, polyF, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := overlaps(tt.poly1, tt.poly2); got != tt.want {
				t.Errorf("overlaps(%v, %v) = %v, want %v", tt.poly1, tt.poly2, got, tt.want)
			}
		})
	}
}
