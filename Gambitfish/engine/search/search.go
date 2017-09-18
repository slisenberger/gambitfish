package search

import "fmt"
import "math"
import "../../game"
import "../evaluate"

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) (float64, *game.Move) {
	over, winner := b.CalculateGameOver()
	if over || depth == 0 {
		if over {
			fmt.Println("in search tree found winner: %v" + winner.String())
		}
		return e.Evaluate(b), nil
	}
	_ = alpha
	_ = beta
	var best game.Move
	var eval float64
	moves := b.AllLegalMoves()
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
