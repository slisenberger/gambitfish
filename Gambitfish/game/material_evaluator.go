package game

type MaterialEvaluator struct{}

// Utility for printing in debugging.
func ColorToString(c Color) string {
	if c == 1 {
		return "WHITE"
	} else {
		return "BLACK"
	}
}

// Evaluates a board by counting the material weights for all remaining pieces.
func (m MaterialEvaluator) Evaluate(b *Board) float64 {
	eval := 0.0
	for _, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() == BLACK {
			eval -= p.Value()
		} else {
			eval += p.Value()
		}
	}
	return float64(b.Active) * eval
}
