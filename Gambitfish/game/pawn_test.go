package game

import "testing"

func initBlockedPawnBoard() *Board {
	b := &Board{Active: WHITE}
	b.Squares[8] = &Pawn{&BasePiece{C: WHITE, B: b}}
	b.Squares[16] = &Pawn{&BasePiece{C: WHITE, B: b}}
	b.InitPieceSet()
	return b
}

func initLonePawnThirdRank() *Board {
	b := &Board{Active: WHITE}
	b.Squares[16] = &Pawn{&BasePiece{C: WHITE, B: b}}
	b.InitPieceSet()
	return b
}

func initFullPawnCaptures() *Board {
	b := &Board{Active: WHITE}
	b.Squares[9] = &Pawn{&BasePiece{C: WHITE, B: b}}
	b.Squares[16] = &Pawn{&BasePiece{C: BLACK, B: b}}
	b.Squares[18] = &Pawn{&BasePiece{C: BLACK, B: b}}
	b.InitPieceSet()
	return b
}
func initPawnBlockedByOtherTeam() *Board {
	b := &Board{Active: WHITE}
	b.Squares[8] = &Pawn{&BasePiece{C: WHITE, B: b}}
	b.Squares[16] = &Pawn{&BasePiece{C: BLACK, B: b}}
	b.InitPieceSet()
	return b
}

func TestLegalPawnMoves(t *testing.T) {
	testCases := []struct {
		name       string
		board      *Board
		pieceIndex int
		want       int
	}{{
		name:       "starting board pawns should have two moves",
		board:      DefaultBoard(),
		pieceIndex: 9,
		want:       2,
	}, {
		name:       "pawn with a pawn in front should not have moves",
		board:      initBlockedPawnBoard(),
		pieceIndex: 8,
		want:       0,
	}, {
		name:       "pawn on 3rd rank should have one move",
		board:      initLonePawnThirdRank(),
		pieceIndex: 16,
		want:       1,
	}, {
		name:       "starting pawn with captures on both sides",
		board:      initFullPawnCaptures(),
		pieceIndex: 9,
		want:       4,
	}, {
		name:       "pawn blocked by other team",
		board:      initPawnBlockedByOtherTeam(),
		pieceIndex: 8,
		want:       0,
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
