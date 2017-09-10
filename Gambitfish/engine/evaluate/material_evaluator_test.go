package evaluate

import "testing"
import "../../game"

func TestMaterialEvaluation(t *testing.T) {
	b := game.DefaultBoard()
	e := MaterialEvaluator{}
	if e.Evaluate(b) != 0 {
		t.Error("material evaluation for default board unequal: ")
	}
}
