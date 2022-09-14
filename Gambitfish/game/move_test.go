package game

import "testing"

func TestEfficientMove(t *testing.T) {
	testCases := []struct {
		m Move
	}{
		{
			m: NewMoveNoBoard(WHITEPAWN, GetSquare(4,2), GetSquare(2,2)),
		},
		{
			m: NewMoveNoBoard(BLACKKING, GetSquare(4,4), GetSquare(5,5)),
		},
	}

	for _, tc := range testCases {
		m := MoveToEfficientMove(tc.m)
		got := EfficientMoveToMove(m)
		if got != tc.m {
			t.Errorf("efficient move conversion failed: got %v, want %v", got, tc.m)

		}
	}
}
