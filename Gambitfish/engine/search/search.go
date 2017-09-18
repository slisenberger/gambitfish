package search

import "math"
import "../../game"
import "../evaluate"

// MAX_QUIESCENCE_DEPTH is the number of extra nodes to search if
// We reach depth 0 with pending captures.
// This should eventually be controlled, but for now we max quiescence search
// another nodes.
const MAX_QUIESCENCE_DEPTH = -1

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) (float64, *game.Move) {
	over, _ := b.CalculateGameOver()
	if over || depth == MAX_QUIESCENCE_DEPTH || (depth <= 0 && IsQuiet(b)) {
		return e.Evaluate(b), nil
	}
	if bm := BookMove(b); bm != nil {
		return 0.0, bm
	}
	_ = alpha
	_ = beta
	var best game.Move
	var eval float64

	moves := b.AllLegalMoves()
	// If we are past our depth limit, we are only in quiescence search.
	// In quiescence search, only search remaining captures.
	if depth <= 0 {
		captures := []game.Move{}
		for _, move := range moves {
			if move.Capture != nil {
				captures = append(captures, move)
			}
		}
		moves = captures
	}
	bestVal := math.Inf(-1)
	for _, move := range moves {
		game.ApplyMove(b, move)
		b.SwitchActivePlayer()
		eval, _ = AlphaBetaSearch(b, e, depth-1, -beta, -alpha)
		// Undo move and restore player.
		b.SwitchActivePlayer()
		game.UndoMove(b, move)
		eval = -1.0 * eval
		// We do >= because if checkmate is inevitable, we still need to pick a move.
		if eval >= bestVal {
			bestVal = eval
			best = move
		}
		if eval > alpha {
			alpha = eval
		}
		if alpha >= beta {
			break
		}
	}
	return bestVal, &best
}

func IsQuiet(b *game.Board) bool {
	for _, m := range b.AllLegalMoves() {
		if m.Capture != nil {
			return false
		}
	}
	return true
}
