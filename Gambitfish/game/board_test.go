package game

import "testing"

func TestDefaultBoard(t *testing.T) {
	b := DefaultBoard()

	for i, piece := range b.Squares {
		if piece == nil {
			continue
		}
		if piece.Square().Index() != i {
			t.Errorf("wrong index for default piece: got %v, want %v", piece.Square().Index(), i)

		}

	}

}
