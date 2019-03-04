package game

import "reflect"
import "testing"

func TestLegalKingMovesDict(t *testing.T) {
	InitInternalData()
	testCases := []struct {
		s    Square
		want []Square
	}{
		{
			A1,
			[]Square{B1, A2, B2},
		}, {
			D2,
			[]Square{C1, D1, E1, C2, E2, C3, D3, E3},
		}, {
			H1,
			[]Square{G1, G2, H2},
		}, {
			H8,
			[]Square{G7, H7, G8},
		},
	}

	for _, tc := range testCases {
		got := SquaresFromBitBoard(LEGALKINGMOVES[tc.s])
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("wrong legal king moves in initialized dict for %v: got %v, want %v", tc.s, got, tc.want)
		}
	}
}
