package game

type MaterialEvaluator struct{}

// Evaluates a board by counting the material weights for all remaining pieces.
func (m MaterialEvaluator) Evaluate(b *Board) float64 {
	eval := 0.0
	for _, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() == b.Active  {
			eval += p.Value()
		} else {
			eval -= p.Value()
		}
	}
	return eval
}
