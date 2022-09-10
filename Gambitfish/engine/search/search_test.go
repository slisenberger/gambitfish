package search

import "math"
import "testing"
import "../../game"


// Test finding winning checkmates.
func TestMate(t *testing.T) {
	game.InitInternalData()
	testCases := []struct {
		name  string
		fen   string
		move  string // The best move at the given depth
		depth int
	}{
		{
			name:  "blitzkrieg m4",
			fen:   "r1bqkbnr/p1pp1ppp/1pn5/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 0 1",
			move: "Qh5xf7",
			depth: 1,
		},
		{
			name:  "blitzkrieg m4 as black",
			fen:   "rnb1k1nr/pppp1ppp/8/2b1p3/2B1P2q/P1N5/1PPP1PPP/R1BQK1NR b KQkq - 0 1",
			move: "Qh4xf2",
			depth: 1,
		},
		{
			name: "mate in 2, queen sac",
			fen: "2q5/pR6/1p3pnk/1P4pp/8/5QPP/P2r2BK/8 w - - 0 1",
			move: "Qf3xh5",
			depth: 3,
		},
		{
			name: "m2 trapped king zugswang",
			fen: "B7/K1B1p1Q1/5r2/7p/1P1kp1bR/3P3R/1P1NP3/2n5 w - - 0 1",
			move: "Ba8-c6",
			depth: 3,
		},
		{
			name: "m2 queen prep",
			fen: "r1bq2r1/b4pk1/p1pp1p2/1p2pP2/1P2P1PB/3P4/1PPQ2P1/R3K2R w - - 0 1",
			move: "Qd2-h6",
			depth: 3,
		},
	}

	e := game.MaterialEvaluator{}

	for _, tc := range testCases {
           b, err := game.BoardFromFen(tc.fen)
	   if err != nil {
		   t.Error("failed to read board from fen")
	   }
	   _, move, _ := AlphaBetaSearch(b, e, tc.depth, math.Inf(-1), math.Inf(1), false, b.Active)
	   if move.String() != tc.move {
		t.Errorf("Got wrong move. Want %v, got %v", tc.move, move.String())
		b.Print()
	   }
        }
}
