package search

import "fmt"
import "math"
import "../../game"
import "../evaluate"

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) (float64, *game.Move) {
	if b.Finished() || depth == 0 {
		return e.Evaluate(b), nil
	}
	_ = alpha
	_ = beta
	var best game.Move
	var eval float64
	moves := b.AllLegalMoves()
	bestVal := math.Inf(-1)
	for _, move := range moves {
		newBoard := game.CopyBoard(b)
		newBoard.ApplyMove(move)
		newBoard.SwitchActivePlayer()
		eval, _ = AlphaBetaSearch(newBoard, e, depth-1, -beta, -alpha)
		eval = -eval
		if eval > bestVal {
			fmt.Println(fmt.Sprintf("new best eval at depth %v!: %v, %v", depth, eval, move))
			fmt.Println(fmt.Sprintf("above value is negative of value for %v", newBoard.Active))
			bestVal = eval
			best = move
		}
		//		if eval > alpha {
		//			alpha = eval
		//		}
		//		if alpha >= beta {
		//			break
		//		}
	}
	return bestVal, &best
}
