package evaluate

import "math"
import "../../game"

type MaterialEvaluator struct{}

// Utility for printing in debugging.
func ColorToString(c game.Color) string {
	if c == 1 {
		return "WHITE"
	} else {
		return "BLACK"
	}
}

// Evaluates a board by counting the material weights for all remaining pieces.
func (m MaterialEvaluator) Evaluate(b *game.Board) float64 {
	eval := 0.0
	if b.Winner == b.Active {
		return math.Inf(1)
	}
	for piece, _ := range b.PieceSet {
		if piece.Color() == b.Active {
			eval += piece.Value()
		} else {
			eval -= piece.Value()
		}
	}
	return eval
}
