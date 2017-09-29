// Transposition manages transposition tables for avoiding redoing calculation.
package game

type EvalPrecision int

const (
	EvalExact = EvalPrecision(iota)
	EvalLowerBound
	EvalUpperBound
)

type TTEntry struct {
	Depth     int           // the depth this entry was searched to
	Eval      float64       // What this was evaluated as.
	Precision EvalPrecision // Whether we evaluated this node as an alpha/beta cutoff.
	BestMove  Move
}
