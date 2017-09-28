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

func TestRayAttacksDict(t *testing.T) {
	InitInternalData()
	testCases := []struct {
		s    Square
		dir  Direction
		want []Square
	}{
		{
			A1,
			SW,
			[]Square{},
		}, {
			A1,
			NE,
			[]Square{B2, C3, D4, E5, F6, G7, H8},
		}, {
			H8,
			NE,
			[]Square{},
		}, {
			H8,
			SW,
			[]Square{A1, B2, C3, D4, E5, F6, G7},
		}, {
			D5,
			N,
			[]Square{D6, D7, D8},
		}, {
			D5,
			S,
			[]Square{D1, D2, D3, D4},
		}, {
			C1,
			NE,
			[]Square{D2, E3, F4, G5, H6},
		}, {
			C1,
			NW,
			[]Square{B2, A3},
		}, {
			C1,
			SE,
			[]Square{},
		}, {
			C1,
			SW,
			[]Square{},
		}, {
			C1,
			E,
			[]Square{D1, E1, F1, G1, H1},
		}, {
			C1,
			W,
			[]Square{A1, B1},
		},
	}

	for _, tc := range testCases {
		got := SquaresFromBitBoard(RAY_ATTACKS[tc.dir][tc.s])
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("wrong legal ray moves in initialized dict for %v in dir %v: got %v, want %v", tc.s, tc.dir, got, tc.want)
		}

	}
}
