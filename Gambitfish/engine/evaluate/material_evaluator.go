package evaluate

import "../../game"

type MaterialEvaluator struct{}

// Evaluates a board by counting the material weights for all remaining pieces.
func (m MaterialEvaluator) Evaluate(b *game.Board) float64 {
	eval := 0.0
	for _, piece := range b.Squares {
		if piece == nil {
			continue
		}
		if piece.Color() == b.Active {
			eval += piece.Value()
		} else {
			eval -= piece.Value()
		}
	}
	return eval

}
