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
				Capture: &Capture{
					Piece: BLACKBISHOP,
					Square: E4,
				},
			},
		},{
			m: Move{
				Piece: WHITEPAWN,
				Old: E7,
				Square: E8,
				Promotion: WHITEQUEEN,
			},
		},
	}

	for _, tc := range testCases {
		m := MoveToEfficientMove(tc.m)
		if tc.m.Capture != nil {
			m.AddCapture(tc.m.Capture.Piece, tc.m.Capture.Square)
		}
		got := EfficientMoveToMove(m)
		if !got.Equals(tc.m) {
			t.Errorf("efficient move conversion failed: got %v, want %v", got, tc.m)

		}
	}
}
