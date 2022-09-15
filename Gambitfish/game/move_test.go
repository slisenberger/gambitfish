package game

import "testing"

func TestEfficientMove(t *testing.T) {
	testCases := []struct {
		m Move
	}{
		{
			m: Move{
				Piece: WHITEPAWN,
				Old: E2,
				Square: E4,
				TwoPawnAdvance: true,
			},
		},{
			
			m: Move{
				Piece: WHITEROOK,
				Old: E1,
				Square: E4,
				Capture: BLACKBISHOP,
			},
		},{
			m: Move{
				Piece: WHITEPAWN,
				Old: E7,
				Square: E8,
				Promotion: WHITEQUEEN,
			},
		},{
			m: Move{
				Piece: BLACKPAWN,
				Old: B5,
				Square: A6,
				Capture: WHITEPAWN,
				EnPassant: true,
			},
		},
	}

	for _, tc := range testCases {
		m := MoveToEfficientMove(tc.m)
		if tc.m.Capture != NULLPIECE {
			m.AddCapture(tc.m.Capture)
		}
		got := EfficientMoveToMove(m)
		if !got.Equals(tc.m) {
			t.Errorf("efficient move conversion failed: got %v, want %v", got, tc.m)

		}
	}
}
