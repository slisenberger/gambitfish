package search

import "math"
import "../../game"
import "../evaluate"

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) (float64, *game.Move) {
	if b.Finished() || depth == 0 {
		return e.Evaluate(b), nil
	}
	var best game.Move
	moves := b.AllLegalMoves()
	bestVal := math.Inf(-1)
	for _, move := range moves {
		newBoard := b.Copy()
		newBoard.ApplyMove(move)
		newBoard.SwitchActivePlayer()
		eval, _ := AlphaBetaSearch(newBoard, e, depth-1, -beta, -alpha)
		eval = -1.0 * eval
		if eval > alpha {
			alpha = eval
		}
		if eval > bestVal {
			bestVal = eval
			best = move
		}
		if alpha >= beta {
			break
		}
	}
	return bestVal, &best
}
