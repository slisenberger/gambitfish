package game

type CompoundEvaluator struct {
	Evaluators []Evaluator
}

func (e CompoundEvaluator) Evaluate(b *Board) float64 {
	result := 0.0
	for _, e := range e.Evaluators {
		result += e.Evaluate(b)
	}
	return result
}
