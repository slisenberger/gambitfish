// An implementation of a King.
package game

type King struct {
	*BasePiece
}

func (p *King) LegalMoves() []Move {
	moves := KingMoves(p, p.Board().PieceSet[p])
	moves = append(moves, CastlingMoves(p, p.Board().PieceSet[p])...)
	return moves
}

func (p *King) Attacking() []Square {
	moves := KingMoves(p, p.Board().PieceSet[p])
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
}

func (p *King) String() string {
	return "K"
}

func (p *King) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♔"
	case WHITE:
		return "♚"
	}
	return ""
}

func (p *King) Value() float64 {
	return 100.0
}

func (p *King) Type() PieceType {
	return KING
}

// LegalKingMovesDict returns a 64-indexed set of bitboards
// for the legal moves a king can make. This is intended to be used
// as a precomputed index for constructing valid king moves.
func LegalKingMovesDict() [64]uint64 {
	var legalMoves [64]uint64
	// For each square on the board, prepopulate the legal moves.
	var i uint64
	var bb uint64
	var kp uint64
	for i = 0; i < 64; i++ {
		// king position
		kp = 1 << i
		bb = 0
		s := Square(i)
		// Check for going off the left side
		if s.Col() != 1 {
			bb = bb | (kp << 7) | (kp >> 9) | (kp >> 1)
		}
		// Check for going off the right side
		if s.Col() != 8 {
			bb = bb | (kp >> 7) | (kp << 9) | (kp << 1)
		}
		bb = bb | (kp >> 8) | (kp << 8)
		legalMoves[i] = bb
	}
	return legalMoves
}
