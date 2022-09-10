package search

import "fmt"
import "math"
import "../../game"

// MAX_QUIESCENCE_DEPTH is the number of extra nodes to search if
// We reach depth 0 with pending captures.
// This should eventually be controlled, but for now we max quiescence search
// another nodes.
// TODO: Turning off quiescene to fix material evaluation search bugs.
const MAX_QUIESCENCE_DEPTH = -2

const NULL_MOVE_REDUCED_SEARCH_DEPTH = 2

// An Alpha Beta Negamax implementation. Function stolen from here:
// https://en.wikipedia.org/wiki/Negamax#Negamax_with_alpha_beta_pruning
func AlphaBetaSearch(b *game.Board, e game.Evaluator, depth int, alpha, beta float64, nullMove bool, c game.Color) (float64, *game.Move, int) {
	// The number of nodes searched.
	nodes := 0
	// Store original values for transposition table.
	alphaOrig := alpha
	// Check the transposition table for work we've already done, and either
	// return or update our cutoffs.
	if entry, ok := game.TranspositionTable[game.ZobristHash(b)]; ok && entry.Depth >= depth {
		switch entry.Precision {
		case game.EvalExact:
			return entry.Eval, &entry.BestMove, 1
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
			return entry.Eval, &entry.BestMove, 1
		}
	}
	// Return an eval if the game is over.
	over, _ := b.CalculateGameOver()
	if over {
			return math.Inf(-1), nil, 1
		}

	// Evaluate any leaf nodes.
	if (depth <= MAX_QUIESCENCE_DEPTH) || (depth <= 0 && IsQuiet(b)) {
		return e.Evaluate(b), nil, 1
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
		if game.IsCheck(b, b.Active) {
			moves = game.OrderMoves(b, b.AllLegalMoves())
		} else {
			moves = game.OrderMoves(b, b.AllLegalCaptures())
		}
	} else {
		moves = game.OrderMoves(b, b.AllLegalMoves())
	}

	// Try a null move first. If we can prune the search tree without
	// moving, we should.
	if nullMove && !game.IsCheck(b, c) {
		// NullMoves affect en passant state, so we need to remember it.
		epSquare := b.EPSquare
		b.EPSquare = game.OFFBOARD_SQUARE
		var n int
	        b.SwitchActivePlayer()	
		eval, _, n = AlphaBetaSearch(b, e, depth-1-NULL_MOVE_REDUCED_SEARCH_DEPTH, -beta, -alpha, false, -c)
		// negamax
		eval = -1 * eval
	        b.SwitchActivePlayer()	
		b.EPSquare = epSquare
		if eval >= beta {
			return eval, nil, nodes + n
		}
	}
	bestVal := math.Inf(-1)
	for _, move := range moves {
		if move.Piece == game.NULLPIECE {
			b.Print()
			fmt.Println(b.AllLegalMoves())
			fmt.Println(b.AllLegalCaptures())
			fmt.Println(b.AllLegalChecks())
			fmt.Println(depth)
			fmt.Println(moves)
			panic("somehow am making a null move..")
		}
		game.ApplyMove(b, move)
		b.SwitchActivePlayer()
		eval, _, n := AlphaBetaSearch(b, e, depth-1, -beta, -alpha, false, -c)
		// Negate eval -- it's opponent's opinion!
		eval = -1 * eval
	        b.SwitchActivePlayer()	
		nodes += n
		// Undo move and restore player.
		game.UndoMove(b, move)
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
	return bestVal, &best, nodes
}

func IsQuiet(b *game.Board) bool {
	return (len(b.AllLegalCaptures()) == 0) && (len(b.AllLegalChecks()) == 0) && !game.IsCheck(b, b.Active)
}
