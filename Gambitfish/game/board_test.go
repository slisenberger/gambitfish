package game

import "testing"

func TestDefaultBoard(t *testing.T) {
	b := DefaultBoard()

	for i, piece := range b.Squares {
		if piece == nil {
			continue
		}
		s := piece.Square()
		if s.Index() != i {
			t.Errorf("wrong index for default piece: got %v, want %v", s.Index(), i)

		}

	}

}
