package game

import "reflect"
import "testing"

func TestSquaresFromBitboard(t *testing.T) {
	InitInternalData()
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

func TestSetBitOnBoard(t *testing.T) {
	testCases := []struct {
		name     string
		bitboard uint64
		s        Square
		want     uint64
	}{
		{
			"successfully unsets square on h8",
			uint64(0x8000000000000000),
			H8,
			uint64(0),
		},
		{
			"successfully unsets square on h8, full board",
			uint64(0xFFFFFFFFFFFFFFFF),
			H8,
			uint64(0x7FFFFFFFFFFFFFFF),
		},
	}

	for _, tc := range testCases {
		got := UnSetBitOnBoard(tc.bitboard, tc.s)
		if got != tc.want {
			t.Errorf("%v: got %v unsetting square %v, wanted %v", tc.name, got, tc.s, tc.want)

		}
		got = SetBitOnBoard(got, tc.s)
		if got != tc.bitboard {
			t.Errorf("%v: got %v setting square %v, wanted %v", tc.name, got, tc.s, tc.bitboard)

		}
	}
}
