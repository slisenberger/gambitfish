package search

import "fmt"
import "math"
import "../../game"
import "../evaluate"

func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) *game.Move {
	if b.Finished() || depth == 0 {
		return nil
	}
	var best game.Move
	score := math.Inf(-1)
	moves := b.AllLegalMoves()
	// THIS IS NOT ALPHA BETA, but let's see how it works!
	for _, move := range moves {
		fmt.Println("legal move: " + move.String())
		newBoard := b.Copy()
		newBoard.ApplyMove(move)
		eval := e.Evaluate(newBoard)
		if eval > score {
			fmt.Println(fmt.Sprintf("new best move with score of %v: %v", eval, move.String()))
			score = eval
			best = move
		}
	}
	_ = alpha
	_ = beta
	return &best
}
