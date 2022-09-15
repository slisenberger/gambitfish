// Transposition manages transposition tables for avoiding redoing calculation.
package game

type EvalPrecision int

// The transposition table holding a list of previously seen positions and
// their evaluation.
var TranspositionTable = map[uint64]TTEntry{}

const (
	EvalExact = EvalPrecision(iota)
	EvalLowerBound
	EvalUpperBound
)

type TTEntry struct {
	Depth     int           // the depth this entry was searched to
	Eval      float64       // What this was evaluated as.
	Precision EvalPrecision // Whether we evaluated this node as an alpha/beta cutoff.
	BestMove  EfficientMove
	Ancient   bool          // Whether this has gone a full evaluation without being accessed.
}

func EraseOldTableEntries() {
	var tt = map[uint64]TTEntry{}
	for k, v := range TranspositionTable {
		if !v.Ancient {
			v.Ancient = true
			tt[k] = v
		}
	}
	TranspositionTable = tt
}
