package main

import (
	"testing"
)

func TestMove(t *testing.T) {
	type TestCase struct {
		name  string
		start Fifteen
		x     int
		y     int
		want  Fifteen
	}

	start := Fifteen{
		Grid: [16]int{
			1, 2, 3, 4,
			5, 6, -1, 8,
			9, 10, 11, 12,
			13, 14, 15, 7,
		},
	}

	testCases := []TestCase{
		{
			name:  "move 3 down",
			start: start,
			x:     2,
			y:     0,
			want: Fifteen{
				Grid: [16]int{
					1, 2, -1, 4,
					5, 6, 3, 8,
					9, 10, 11, 12,
					13, 14, 15, 7,
				},
			},
		},
		{
			name:  "move 5,6 right",
			start: start,
			x:     0,
			y:     1,
			want: Fifteen{
				Grid: [16]int{
					1, 2, 3, 4,
					-1, 5, 6, 8,
					9, 10, 11, 12,
					13, 14, 15, 7,
				},
			},
		},
		{
			name:  "move 6 right",
			start: start,
			x:     1,
			y:     1,
			want: Fifteen{
				Grid: [16]int{
					1, 2, 3, 4,
					5, -1, 6, 8,
					9, 10, 11, 12,
					13, 14, 15, 7,
				},
			},
		},
		{
			name:  "move 8 left",
			start: start,
			x:     3,
			y:     1,
			want: Fifteen{
				Grid: [16]int{
					1, 2, 3, 4,
					5, 6, 8, -1,
					9, 10, 11, 12,
					13, 14, 15, 7,
				},
			},
		},
		{
			name:  "move 11 up",
			start: start,
			x:     2,
			y:     2,
			want: Fifteen{
				Grid: [16]int{
					1, 2, 3, 4,
					5, 6, 11, 8,
					9, 10, -1, 12,
					13, 14, 15, 7,
				},
			},
		},
		{
			name:  "move 11,15 up",
			start: start,
			x:     2,
			y:     3,
			want: Fifteen{
				Grid: [16]int{
					1, 2, 3, 4,
					5, 6, 11, 8,
					9, 10, 15, 12,
					13, 14, -1, 7,
				},
			},
		},
	}

	for _, tc := range testCases {
		board := tc.start
		err := board.Move(tc.x, tc.y)
		if err != nil {
			t.Errorf("%s: unexpected error %v", tc.name, err)
		}
		if board != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, board, tc.want)
		}
	}
}

func TestIsSolvable(t *testing.T) {
	type TestCase struct {
		name  string
		start [16]int
		want  bool
	}

	testCases := []TestCase{
		{
			name:  "in order",
			start: [16]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, -1},
			want:  true,
		},
		{
			name:  "last 2 numbers switched",
			start: [16]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 13, -1},
			want:  false,
		},
	}

	for _, tc := range testCases {
		got := isSolvable(tc.start)
		if got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}
