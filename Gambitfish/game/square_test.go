package game

import "testing"

func TestAsString(t *testing.T) {
	testCases := []struct {
		square Square
		want   string
	}{
		{
			square: GetSquare(1, 1),
			want:   "a1",
		}, {
			square: GetSquare(1, 8),
			want:   "h1", // 1st Row is "1", 8th Column is "h"
		}, {

			square: GetSquare(8, 8),
			want:   "h8",
		},
	}

	for _, tc := range testCases {
		got := tc.square.String()
		if got != tc.want {
			t.Errorf("string representation for %v, %v: got %v, want %v", tc.square.Row(), tc.square.Col(), got, tc.want)

		}
	}
}

func TestIndex(t *testing.T) {
	testCases := []struct {
		square Square
		want   int
	}{
		{
			square: GetSquare(1, 1),
			want:   0,
		},
		{
			square: GetSquare(8, 8),
			want:   63,
		},
		{
			square: GetSquare(6, 8),
			want:   47,
		},
	}
	for _, tc := range testCases {
		got := int(tc.square)
		if got != tc.want {
			t.Errorf("index for %v, %v incorrect: got %v, want %v", tc.square.Row(), tc.square.Col(), got, tc.want)
		}
	}
}
