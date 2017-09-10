package evaluate

import "../../game"

type MaterialEvaluator struct{}

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
