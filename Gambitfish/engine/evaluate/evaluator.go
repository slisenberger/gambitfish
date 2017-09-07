// package evaluate provides an interface for evaluating a chess board's value
// for the active player.
package evaluate

import "../../game"

type Evaluator interface {
	Evaluate(*game.Board) float
}
