package game

type MobilityEvaluator struct{}

// Evaluates a board by counting the number of legal moves available.
func (m MobilityEvaluator) Evaluate(b *Board) float64 {
	eval := .01 * float64(len(b.AllLegalMoves()))
	b.SwitchActivePlayer()
	eval -= .01 * float64(len(b.AllLegalMoves()))
	b.SwitchActivePlayer()
	return eval
}
