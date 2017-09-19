package game

import "testing"

func TestDefaultBoard(t *testing.T) {
	b := DefaultBoard()

	for i, piece := range b.Squares {
		if piece == nil {
			continue
		}
		s := b.PieceSet[piece]
		if int(s) != i {
			t.Errorf("wrong index for default piece %v: got %v, want %v", piece, s, i)

		}
	}

}
