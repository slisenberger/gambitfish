package evaluate

import "../../game"

type MobilityEvaluator struct{}

// Evaluates a board by counting the number of legal moves available.
func (m MobilityEvaluator) Evaluate(b *game.Board) float64 {
	eval := .1 * float64(len(b.AllLegalMoves()))
	b.SwitchActivePlayer()
	eval -= .1 * float64(len(b.AllLegalMoves()))
	b.SwitchActivePlayer()
	return eval
}
