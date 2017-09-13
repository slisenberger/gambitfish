package search

import "math"
import "testing"
import "../../game"
import "../evaluate"

func initGreedyQueenBoard() *game.Board {
	b := &game.Board{Active: game.WHITE}
	b.Squares[1] = &game.Queen{&game.BasePiece{C: game.WHITE, B: b}}
	b.Squares[0] = &game.King{&game.BasePiece{C: game.WHITE, B: b}, false}
	b.Squares[56] = &game.Pawn{&game.BasePiece{C: game.BLACK, B: b}}
	b.Squares[57] = &game.King{&game.BasePiece{C: game.BLACK, B: b}, false}
	b.Squares[49] = &game.Pawn{&game.BasePiece{C: game.BLACK, B: b}}
	b.InitPieceSet()
	return b
}

// Test that at depth=2 with a material evaluator, the white player doesn't
// greedily grab a guarded pawn.
func TestSacrificeQueen(t *testing.T) {
	b := initGreedyQueenBoard()
	e := evaluate.MaterialEvaluator{}
	// With one ply, the greedy option is take the pawn!
	_, move := AlphaBetaSearch(b, e, 1, math.Inf(-1), math.Inf(1))
	if move.String() != "Qb1 to b7" {
		t.Errorf("queen did not take greedy option qb7, but instead %v", move)
	}
	// With two plies, it should be clear that taking the pawn is a bad
	// sacrifice!
	_, move = AlphaBetaSearch(b, e, 2, math.Inf(-1), math.Inf(1))
	if move.String() == "Qb1 to b7" {
		t.Error("Queen sacrificed on position with depth=2, which is losing")
	}
}
