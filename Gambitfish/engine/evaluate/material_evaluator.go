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
	if over, winner := b.CalculateGameOver(); over {
		if winner == b.Active {
			return math.Inf(1)
		} else if winner == -1*b.Active {
			return math.Inf(-1)
		} else {
			return 0.0
		}
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
