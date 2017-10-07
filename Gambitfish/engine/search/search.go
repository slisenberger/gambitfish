package search

import "math"
import "../../game"
import "../evaluate"

// MAX_QUIESCENCE_DEPTH is the number of extra nodes to search if
// We reach depth 0 with pending captures.
// This should eventually be controlled, but for now we max quiescence search
// another nodes.
const MAX_QUIESCENCE_DEPTH = -3

const NULL_MOVE_REDUCED_SEARCH_DEPTH = 2

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e evaluate.Evaluator, depth int, alpha, beta float64) (float64, *game.Move) {
	// Store original values for transposition table.
	alphaOrig := alpha
	// Check the transposition table for work we've already done, and either
	// return or update our cutoffs.
	if entry, ok := game.TranspositionTable[game.ZobristHash(b)]; ok && entry.Depth >= depth {
		switch entry.Precision {
		case game.EvalExact:
			return entry.Eval, &entry.BestMove
		case game.EvalLowerBound:
			if entry.Eval > alpha {
				alpha = entry.Eval
			}
		case game.EvalUpperBound:
			if entry.Eval < beta {
				beta = entry.Eval
			}
		}
		if alpha >= beta {
			return entry.Eval, &entry.BestMove
		}
	}
	// Return an eval if the game is over.
	over, winner := b.CalculateGameOver()
	if over {
		if b.Active == winner {
			return math.Inf(1), nil
		} else if b.Active == -1*winner {
			return math.Inf(-1), nil
		} else {
			return 0, nil
		}
	}

	// Evaluate any leaf nodes.
	if depth == MAX_QUIESCENCE_DEPTH || (depth <= 0 && IsQuiet(b)) {
		return e.Evaluate(b), nil
	}
	// TODO(slisenberger): I'd like to eventually ignore book moves, seeing if we can do decent
	// from the opening.
	//	if bm := BookMove(b); bm != nil {
	//		return 0.0, bm
	//	}

	var best game.Move
	var eval float64

	var moves []game.Move
	// If we are past our depth limit, we are only in quiescence search.
	// In quiescence search, only search remaining captures.
	if depth <= 0 {
		moves = game.OrderMoves(b, b.AllLegalCaptures())
	} else {
		moves = game.OrderMoves(b, b.AllLegalMoves())
	}
	bestVal := math.Inf(-1)
	for _, move := range moves {
		game.ApplyMove(b, move)
		b.SwitchActivePlayer()
		eval, _ = AlphaBetaSearch(b, e, depth-1, -beta, -alpha)
		// Undo move and restore player.
		b.SwitchActivePlayer()
		game.UndoMove(b, move)
		eval = -1.0 * eval
		// We do >= because if checkmate is inevitable, we still need to pick a move.
		if eval >= bestVal {
			bestVal = eval
			best = move
		}
		if eval > alpha {
			alpha = eval
		}
		if alpha >= beta {
			break
		}
	}
	// Store values in transposition table.
	hash := game.ZobristHash(b)
	entry := game.TTEntry{Depth: depth, Eval: bestVal, BestMove: best}
	if bestVal <= alphaOrig {
		entry.Precision = game.EvalUpperBound
	} else if bestVal >= beta {
		entry.Precision = game.EvalLowerBound
	} else {
		entry.Precision = game.EvalExact
	}
	game.TranspositionTable[hash] = entry
	return bestVal, &best
}

func IsQuiet(b *game.Board) bool {
	return len(b.AllLegalCaptures()) == 0
}
