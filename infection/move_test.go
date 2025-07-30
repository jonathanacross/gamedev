package main

import "testing"

func TestIndexToHumanCoords(t *testing.T) {
	tests := []struct {
		index int
		want  string
	}{
		{0, "A7"},
		{6, "G7"},
		{7, "A6"},
		{42, "A1"},
		{43, "B1"},
		{48, "G1"},
	}

	for _, tt := range tests {
		got := IndexToHumanCoords(tt.index)
		if got != tt.want {
			t.Errorf("IndexToHumanCoords(%d) = %s; want %s", tt.index, got, tt.want)
		}
	}
}

func TestHumanCoordsToIndex(t *testing.T) {
	tests := []struct {
		coords string
		want   int
	}{
		{"A7", 0},
		{"G7", 6},
		{"A6", 7},
		{"A1", 42},
		{"B1", 43},
		{"G1", 48},
	}

	for _, tt := range tests {
		got := HumanCoordsToIndex(tt.coords)
		if got != tt.want {
			t.Errorf("HumanCoordsToIndex(%s) = %d; want %d", tt.coords, got, tt.want)
		}
	}
}
