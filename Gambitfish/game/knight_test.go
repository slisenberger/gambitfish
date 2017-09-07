package game

import "testing"

func initCentralKnightBoard() *Board {
	b := &Board{active: WHITE}
	b.Squares[36] = &Knight{&BasePiece{color: WHITE, square: &Square{4, 5}, board: b}}
	return b
}

func TestLegalKnightMoves(t *testing.T) {
	testCases := []struct {
		name       string
		board      *Board
		pieceIndex int
		want       int
	}{{
		name:       "starting board knights should have two moves.",
		board:      DefaultBoard(),
		pieceIndex: 1,
		want:       2,
	},
		{
			name:       "central knights should have 8 squares",
			board:      initCentralKnightBoard(),
			pieceIndex: 36,
			want:       8,
		},
	}
	for _, tc := range testCases {
		p := tc.board.Squares[tc.pieceIndex]
		got := len(p.LegalMoves())

		if got != tc.want {
			t.Errorf("%v: wrong number of moves: got %v want %v", tc.name, got, tc.want)
		}
	}
}
