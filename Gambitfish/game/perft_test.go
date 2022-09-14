package game

import "fmt"
import "testing"
import "time"

func TestPerft(t *testing.T) {
	InitInternalData()
	testCases := []struct {
		name  string
		fen   string
		moves []int // The number of legal moves at each successive depth.
	}{
		{
			name:  "starting board",
			fen:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []int{20, 400, 8902, 197281, 4865609},
		}, {
			name:  "kiwipete",
			fen:   "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0",
			moves: []int{48, 2039, 97862, 4085603},
		}, {
			name:  "endgame 1",
			fen:   "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0",
			moves: []int{14, 191, 2812, 43238, 674624},
		}, {
			name:  "promotion-antics",
			fen:   "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			moves: []int{44, 1486, 62379, 2103487},
		}, {
			name:  "most legal moves",
			fen:   "R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - - 0 1",
			moves: []int{218},
		},
	}

	for _, tc := range testCases {
		b, err := BoardFromFen(tc.fen)
		if err != nil {
			t.Errorf("error reading fen string %v: %v", tc.fen, err)
		}
		b.Print()
		fmt.Println(b.AllLegalMoves())

		for i, want := range tc.moves {
			depth := i + 1
			start := time.Now()
			got := PerftHelper(b, depth)
			total := time.Since(start)
			if got != want {
				t.Errorf("%v got wrong result for depth %v: got %v, want %v", tc.name, depth, got, want)
			}
			fmt.Println(fmt.Sprintf("calculated perft for depth %v in time %v", depth, total))

		}
	}
}

// PerftHelper handles the move generation.
func PerftHelper(b *Board, depth int) int {
	legalMoves := b.AllLegalMoves()
	if depth == 1 {
		return len(legalMoves)
	}
	result := 0
	for _, move := range legalMoves {
		bs := ApplyMove(b, move)
		b.SwitchActivePlayer()
		result += PerftHelper(b, depth-1)
		b.SwitchActivePlayer()
		UndoMove(b, move, bs)
	}
	return result
}
