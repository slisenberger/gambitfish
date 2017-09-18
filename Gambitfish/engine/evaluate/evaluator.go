// package evaluate provides an interface for evaluating a chess board's value.
// Positive evaluations mean WHITE has an advantage, and negative evaluations mean BLACK has an advantage.
package evaluate

import "../../game"

type Evaluator interface {
	Evaluate(*game.Board) float64
}
