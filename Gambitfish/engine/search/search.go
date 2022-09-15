package search

import "math"
import "../../game"

// MAX_QUIESCENCE_DEPTH is the number of extra nodes to search if
// We reach depth 0 with pending captures.
// This should eventually be controlled, but for now we max quiescence search
// another nodes.
const MAX_QUIESCENCE_DEPTH = 8

const NULL_MOVE_REDUCED_SEARCH_DEPTH = 2

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e game.Evaluator, depth int, alpha, beta float64, nullMove bool, c game.Color, km game.KillerMoves) (float64, game.EfficientMove, int) {
	// The number of nodes searched.
	nodes := 0
	// Return an eval if the game is over.

	lm := b.AllLegalMoves()
	over, winner := b.CalculateGameOver(lm)
	if over {
		if winner == 0 {
			return 0.0, game.EfficientMove(0), 1
		} else {
			return math.Inf(-1), game.EfficientMove(0), 1
		}
	}
	// Store original values for transposition table.
	alphaOrig := alpha
	// Check the transposition table for work we've already done, and either
	// return or update our cutoffs.
	h := game.ZobristHash(b)
	if entry, ok := game.TranspositionTable[h]; ok && (entry.Depth >= depth) && entry.Position == b.Position {
		// Mark this entry to not be deleted.
		entry.Ancient = false
		game.TranspositionTable[h] = entry
		switch entry.Precision {
		case game.EvalExact:
			return entry.Eval, entry.BestMove, 1
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
			return entry.Eval, entry.BestMove, 1
		}
	}
	// Evaluate any leaf nodes.
	if (depth <= 0) {
		// TODO: Leaf node transposition tables are turned off for bug finding.
		// Store values in transposition table.
		hash := game.ZobristHash(b)
		eval, move, nodes := QuiescenceSearch(b, e, MAX_QUIESCENCE_DEPTH, alpha, beta)
		entry := game.TTEntry{Depth: 0, Eval: eval, BestMove: game.EfficientMove(0), Precision: game.EvalExact, Ancient: false, Position: b.Position}
		// Only store values if they are better values than we've seen before.
		_, ok := game.TranspositionTable[hash]
		if !ok {
			game.TranspositionTable[hash] = entry
		}
		return eval, move, nodes
		//return QuiescenceSearch(b, e, MAX_QUIESCENCE_DEPTH, alpha, beta) // Only store values if they are better values than we've seen before.  
	}

	// TODO(slisenberger): I'd like to eventually ignore book moves, seeing if we can do decent
	// from the opening.
	//	if bm := BookMove(b); bm != nil {
	//		return 0.0, bm
	//	}



	var best game.EfficientMove
	var eval float64
	var moves []game.EfficientMove
	// Check extensions.
	if game.IsCheck(b, b.Active) {
		depth = depth + 1
	}

	// Try a null move first. If we can prune the search tree without
	// moving, we should. We also identify threats in the position this way.
	if nullMove && !game.IsCheck(b, c) {
		// NullMoves affect en passant state, so we need to remember it.
		epSquare := b.EPSquare
		b.EPSquare = game.OFFBOARD_SQUARE
		var n int
	        b.SwitchActivePlayer()	
		eval, _, n = AlphaBetaSearch(b, e, depth-1-NULL_MOVE_REDUCED_SEARCH_DEPTH, -beta, -alpha, false, -c, km)
		// negamax
		eval = -1 * eval
	        b.SwitchActivePlayer()	
		b.EPSquare = epSquare
		if eval >= beta {
			return eval, game.EfficientMove(0), nodes + n
		}
	}

	moves = game.OrderMoves(b, lm, depth, km, false)
	bestVal := math.Inf(-1)
	for i := 0; i < len(moves); i++ {
		move := moves[i]
		// Late Move Reductions. Trim the search space for later moves in our ordering scheme if they are quiet.
		if (i >= 3) && (depth > 3) && move.Capture() == game.NULLPIECE && !game.IsCheck(b, b.Active) && move.Promotion() == game.NULLPIECE {
			// Also exclude moves that give check from reductions.
			// Need a faster way to check if moves are checking moves.
			bs := game.ApplyMove(b, move)
			if !game.IsCheck(b, c * -1) {
				depth = depth - 1
			}
			game.UndoMove(b, move, bs)

		}

		bs := game.ApplyMove(b, move)
		b.SwitchActivePlayer()
		// Temporarily turn off null move reductions.
		eval, _, n := AlphaBetaSearch(b, e, depth-1, -beta, -alpha, false, -c, km)
		// Negate eval -- it's opponent's opinion!
		eval = -1 * eval
		game.UndoMove(b, move, bs)
	        b.SwitchActivePlayer()
		nodes += n
		// Undo move and restore player.
		// We do >= because if checkmate is inevitable, we still need to pick a move.
		if eval >= bestVal {
			bestVal = eval
			best = move
		}
		if eval > alpha {
			alpha = eval
		}
		if alpha >= beta {
			// Non captures that cause beta cutoffs should be tried
			// earlier in sooner iterations.
			if move.Capture() == game.NULLPIECE {
				km.AddKillerMove(depth, move)
		        }
			break
		}
	}
	// Store values in transposition table.
	hash := game.ZobristHash(b)
	entry := game.TTEntry{Depth: depth, Eval: bestVal, BestMove: best, Ancient: false, Position: b.Position}
	if bestVal <= alphaOrig {
		entry.Precision = game.EvalUpperBound
	} else if bestVal >= beta {
		entry.Precision = game.EvalLowerBound
	} else {
		entry.Precision = game.EvalExact
	}

	// Only store values if they are better values than we've seen before, or if
	// no values have been stored, or if a collission.
	old, ok := game.TranspositionTable[hash]
	if !ok || (old.Depth < depth) || old.Position != b.Position {
		game.TranspositionTable[hash] = entry
	}

	return bestVal, best, nodes
}

func QuiescenceSearch(b *game.Board, e game.Evaluator, depth int, alpha, beta float64) (float64, game.EfficientMove, int){
	// The number of nodes searched.
	nodes := 0

	var moves []game.EfficientMove

	// Else, quiescence search only searches legal captures and checks, or check evasions

	qmoves, allmoves := b.AllQuiescenceMoves()


	// Start by making sure the game is still playable
	over, winner := b.CalculateGameOver(allmoves)
	if over {
		if winner == 0 {
			return 0.0, game.EfficientMove(0), 1
		} else {
			return math.Inf(-1), game.EfficientMove(0), 1
		}
	}
	moves = game.OrderMoves(b, qmoves, depth, nil, true)

	// evaluate the position as a stand pat baseline
	eval := e.Evaluate(b)
	// Return normal evaluation from quiet boards at max depth.
	if depth <= 0 || len(moves) == 0 {
		return eval, game.EfficientMove(0), 1
	}
	// Otherwise, use stand pat value to optimize quiescence bounds.
	if eval >= beta {
		return eval, game.EfficientMove(0), 1
	}
//	if eval > alpha {
//		alpha = eval
//	}

	var best game.EfficientMove
	bestVal := math.Inf(-1)
	for _, move := range moves {
		var eval float64
		bs := game.ApplyMove(b, move)
		b.SwitchActivePlayer()
		eval, _, n := QuiescenceSearch(b, e, depth-1, -beta, -alpha)
		// Negate eval -- it's opponent's opinion!
		eval = -1 * eval
		game.UndoMove(b, move, bs)
	        b.SwitchActivePlayer()
		nodes += n
		// Undo move and restore player.
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

	return bestVal, best, nodes
}
