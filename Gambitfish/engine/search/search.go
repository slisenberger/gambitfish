package search

import "../../game"
import "../evaluate"

func AlphaBetaSearch(b *game.Board, e *evaluate.Evaluator, depth int, alpha, beta float64) *game.Move {
	if b.Finished() || depth == 0 {
		return nil
	}
	var best *game.Move = nil
	score := 0
	moves := b.AllLegalMoves()
	// THIS IS NOT ALPHA BETA, but let's see how it works!
	for _, move := range moves {
		newBoard := &game.Board{*b}
		newBoard.ApplyMove(move)
		eval := e.Evaluate(b)
		if eval >= score {
			score = eval
			best = move
		}
	}
	_ = alpha
	_ = beta
	return &best
}
