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
		{
			name: "m2 from game with engine. black blundered f6.",
			fen: "r1bqkbnr/ppppp1p1/n4p1p/8/2PPP3/8/PP3PPP/RNBQKBNR w - - 1 1",
			move: "Qd1-h5",
			depth: 3,
		},
		{
			name: "m1 engine surrendered.",
			fen: "rn2kb1r/pp2pppp/5n2/Q3N3/3P4/4B3/P3qPPP/2R2RK1 w kq - 0 1",
			move: "Rc1-c8",
			depth: 1,
		},{
			name: "don't get mated!",
			fen: "4k3/p2N2pr/p5qB/2b2p2/1rp1b2P/4P3/PP3PP1/2RQ1RK1 w - - 0 1",
			move: "Bh6-g5",
			depth: 5,
		},
	}

	e := game.MaterialEvaluator{}

	for _, tc := range testCases {
           b, err := game.BoardFromFen(tc.fen)
	   if err != nil {
		   t.Error("failed to read board from fen")
	   }
	   _, move, _ := AlphaBetaSearch(b, e, tc.depth, math.Inf(-1), math.Inf(1), false, b.Active, game.NewKillerMoves())
	   if move.String() != tc.move {
		t.Errorf("Got wrong move for test %v. Want %v, got %v",tc.name, tc.move, move.String())
		b.Print()
	   }
        }
}

// Test preferring free material
func TestMaterial(t *testing.T) {
	game.InitInternalData()
	testCases := []struct {
		name  string
		fen   string
		move  string // The best move at the given depth
		depth int
	}{
		{
			name:  "free early queen for black",
			fen:   "rnbqkb1r/pppp1ppp/5n2/4p3/2B1P1Q1/8/PPPP1PPP/RNB1K1NR b KQkq - 0 1",
			move: "Nf6xg4",
			depth: 1,
		},
		{
			name: "free early queen for white",
			fen: "rnb1kbnr/pppp1ppp/8/4p1q1/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 1",
			move: "Nf3xg5",
			depth: 1,
		},
	}

	e := game.MaterialEvaluator{}

	for _, tc := range testCases {
           b, err := game.BoardFromFen(tc.fen)
	   if err != nil {
		   t.Error("failed to read board from fen")
	   }
	   eval, move, _ := AlphaBetaSearch(b, e, tc.depth, math.Inf(-1), math.Inf(1), false, b.Active, game.NewKillerMoves())
	   if move.String() != tc.move {
		t.Errorf("Got wrong move in test %v. Want %v, got %v (eval %v)", tc.name, tc.move,  move.String(), eval)
		b.Print()
	   }
        }
}
