package game

import "fmt"

import "math/rand"

type Capture struct {
	Piece  Piece
	Square Square
}

// Moves have Pieces and squares
type Move struct {
	Piece           Piece
	Square          Square
	Old             Square
	Capture         *Capture
	EnPassant       bool
	KSCastle        bool
	QSCastle        bool
	Promotion       Piece // Applicable for only Pawn moves
	PrevWQSCastling bool
	PrevWKSCastling bool
	PrevBQSCastling bool
	PrevBKSCastling bool
	PrevCheck       bool
	PrevEPSquare    Square
	TwoPawnAdvance  bool // For En Passant Management.
}

func (m Move) String() string {
	if m.QSCastle {
		return "O-O-O"
	} else if m.KSCastle {
		return "O-O"
	}
	mv := fmt.Sprintf("%v%v", m.Piece, m.Old)
	if m.Capture != nil {
		mv += "x"
	} else {
		mv += "-"
	}
	mv += m.Square.String()
	if m.Promotion != NULLPIECE {
		mv = fmt.Sprintf("%v=%v", mv, m.Promotion)
	}
	return mv
}

func NewMove(p Piece, square Square, old Square, b *Board) Move {
	return Move{
		Piece:           p,
		Square:          square,
		Old:             old,
		EnPassant:       false,
		KSCastle:        false,
		QSCastle:        false,
		PrevBQSCastling: b.BQSCastling,
		PrevWQSCastling: b.WQSCastling,
		PrevBKSCastling: b.BKSCastling,
		PrevWKSCastling: b.WKSCastling,
		Promotion:       NULLPIECE,
		PrevEPSquare:    b.EPSquare,
		TwoPawnAdvance:  false,
	}
}

// Order the moves in an intelligent way for alpha beta pruning.
func OrderMoves(b *Board, moves []Move) []Move {
	// Just to get things off the ground, we'll shuffle the moves, just to get some variety
	// in the AI vs AI games.
	for i := range moves {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	//This is our result array and map of seen moves
	res := make([]Move, len(moves))
	seen := make(map[string]bool, len(moves))

	i := 0

	// Start with what we already believe the best move is.
	if entry, ok := TranspositionTable[ZobristHash(b)]; ok {
		res[0] = entry.BestMove
		seen[entry.BestMove.String()] = true
		i = 1
	}

	// Loop through the move list the rest of the times for other orderings.
	for {
		if i >= len(moves) {
			break
		}
		// Find MVV/LVA captures
		mvv := 0.0    // Most valuable victim
		lva := 1000.0 // least valuable attacker seeing that victim so far.
		var best Move
		var bestNonCapture Move
		for _, m := range moves {
			// Skip moves we've already ordered
			if seen[m.String()] {
				continue
			}
			if m.Capture != nil {
				// If it's our new best mvv, and also least valuable attacker,
				// it's the best move so far.
				if m.Capture.Piece.Value() >= mvv && m.Piece.Value() < lva {
					mvv = m.Capture.Piece.Value()
					lva = m.Piece.Value()
					best = m
				}
			} else {
				bestNonCapture = m
			}
		}
		// Add to results, and don't loop through this move again.
		// If we found a victim at all.
		if mvv > 0.0 {
			res[i] = best
			seen[best.String()] = true
		} else {
			res[i] = bestNonCapture
			seen[bestNonCapture.String()] = true
		}
		i++
	}
	return res
}
