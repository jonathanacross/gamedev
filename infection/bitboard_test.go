package main

import (
	"reflect"
	"testing"
)

func TestBitBoardSet(t *testing.T) {
	var b BitBoard
	b = b.Set(0)
	if b != 1 {
		t.Errorf("Expected 1, got %d", b)
	}
	b = b.Set(63)
	if b != (1<<63)|1 {
		t.Errorf("Expected (1<<63)|1, got %d", b)
	}
}

func TestBitBoardClear(t *testing.T) {
	var b BitBoard = (1 << 63) | 1
	b = b.Clear(0)
	if b != (1 << 63) {
		t.Errorf("Expected (1<<63), got %d", b)
	}
	b = b.Clear(63)
	if b != 0 {
		t.Errorf("Expected 0, got %d", b)
	}
}

func TestBitBoardGet(t *testing.T) {
	var b BitBoard = (1 << 10) | (1 << 20)
	if !b.Get(10) {
		t.Errorf("Expected bit 10 to be set")
	}
	if !b.Get(20) {
		t.Errorf("Expected bit 20 to be set")
	}
	if b.Get(0) {
		t.Errorf("Expected bit 0 to be clear")
	}
}

func TestBitBoardGetSetBitIndices(t *testing.T) {
	tests := []struct {
		name string
		bb   BitBoard
		want []SquareIndex
	}{
		{"empty", 0, []SquareIndex{}},
		{"single bit", 1 << 5, []SquareIndex{5}},
		{"multiple bits", (1 << 1) | (1 << 3) | (1 << 7), []SquareIndex{1, 3, 7}},
		{"high bits", (1 << 62) | (1 << 63), []SquareIndex{62, 63}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.bb.GetSetBitIndices()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("%s: %v want %v", tc.name, got, tc.want)
			}
		})
	}
}
