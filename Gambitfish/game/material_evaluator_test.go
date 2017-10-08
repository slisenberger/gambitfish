package game

import "testing"

func TestMaterialEvaluation(t *testing.T) {
	b := DefaultBoard()
	e := MaterialEvaluator{}
	if e.Evaluate(b) != 0 {
		t.Error("material evaluation for default board unequal: ")
	}
}
