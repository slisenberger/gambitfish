package game

import "testing"

func TestMaterialEvaluation(t *testing.T) {
	b := DefaultBoard()
	e := MaterialEvaluator{}
	if e.Evaluate(b) != 0.0 {
		t.Errorf("material evaluation for default board unequal: %v", e.Evaluate(b))
	}
}
