package search

import "fmt"
import "math"
import "../../game"
import "../evaluate"

func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int) *game.Move {
	if b.Finished() || depth == 0 {
		return nil
	}
	var best game.Move
	score := math.Inf(-1)
	moves := b.AllLegalMoves()
	for _, move := range moves {
		newBoard := b.Copy()
		newBoard.ApplyMove(move)
		eval := Max(newBoard, e, depth, math.Inf(-1), math.Inf(1))
		if eval > score {
			fmt.Println(fmt.Sprintf("new best move with score of %v: %v", eval, move.String()))
			score = eval
			best = move
		}
	}
	return &best
}

func Max(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) float64 {
	if b.Finished() || depth == 0 {
		return e.Evaluate(b)
	}
	moves := b.AllLegalMoves()
	for _, move := range moves {
		newBoard := b.Copy()
		newBoard.ApplyMove(move)
		newBoard.SwitchActivePlayer()
		score := Min(newBoard, e, depth-1, alpha, beta)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}

	}
	return alpha
}

func Min(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) float64 {
	if b.Finished() || depth == 0 {
		return -e.Evaluate(b)
	}
	moves := b.AllLegalMoves()
	for _, move := range moves {
		newBoard := b.Copy()
		newBoard.ApplyMove(move)
		// TODO(slisenberger): I think this is buggy. My guess is switching the active player and
		// running the min eval is causing some messed up moves.
		newBoard.SwitchActivePlayer()
		score := Max(newBoard, e, depth-1, alpha, beta)
		if score <= alpha {
			return alpha
		}
		if score < beta {
			beta = score
		}
	}
	return beta
}
