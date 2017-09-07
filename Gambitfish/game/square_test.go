package game

import "testing"

func TestAsString(t *testing.T) {
	testCases := []struct {
		square *Square
		want   string
	}{
		{
			square: &Square{row: 1, col: 1},
			want:   "a1",
		}, {
			square: &Square{row: 1, col: 8},
			want:   "h1", // 1st row is "1", 8th column is "h"
		}, {

			square: &Square{row: 8, col: 8},
			want:   "h8",
		},
	}

	for _, tc := range testCases {
		got := tc.square.String()
		if got != tc.want {
			t.Errorf("string representation for %v, %v: got %v, want %v", tc.square.row, tc.square.col, got, tc.want)

		}
	}
}

func TestIndex(t *testing.T) {
	testCases := []struct {
		square *Square
		want   int
	}{
		{
			square: &Square{1, 1},
			want:   0,
		},
		{
			square: &Square{8, 8},
			want:   63,
		},
		{
			square: &Square{6, 8},
			want:   47,
		},
	}
	for _, tc := range testCases {
		got := tc.square.Index()
		if got != tc.want {
			t.Errorf("index for %v, %v incorrect: got %v, want %v", tc.square.row, tc.square.col, got, tc.want)
		}
	}
}
