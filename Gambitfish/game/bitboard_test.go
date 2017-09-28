package game

import "reflect"
import "testing"

func TestSquaresFromBitboard(t *testing.T) {
	testCases := []struct {
		name     string
		bitboard uint64
		want     []Square
	}{
		{
			"empty bitboard",
			0,
			[]Square{},
		},
		{
			"only a 1",
			1,
			[]Square{A1},
		},
		{
			"first two squares",
			3,
			[]Square{A1, B1},
		},
	}

	for _, tc := range testCases {
		got := SquaresFromBitBoard(tc.bitboard)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%v: got %v squares from bitboard, wanted %v", tc.name, got, tc.want)

		}

	}

}
