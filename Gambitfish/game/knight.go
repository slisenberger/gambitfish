// An implementation of a Knight
package game

type Knight struct {
	*BasePiece
}

func (p *Knight) LegalMoves() []Move {
	return KnightMoves(p, p.Board().PieceSet[p])
}

func (p *Knight) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
}

// AttackBitboard returns the bitboard for the pieces under attack by this king.
func (p *Knight) AttackBitboard(cur Square) uint64 {
	km := LEGALKNIGHTMOVES[cur]
	switch p.Color() {
	case WHITE:
		km = km & p.Board().Position.BlackPieces
	case BLACK:
		km = km & p.Board().Position.WhitePieces
	}
	return km
}

func (p *Knight) String() string {
	return "N"
}

func (p *Knight) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♘"
	case WHITE:
		return "♞"
	}
	return ""
}

func (p *Knight) Value() float64 {
	return 3.0
}

func (p *Knight) Type() PieceType {
	return KNIGHT
}

// LegalKnightMovesDict returns a 64-indexed set of bitboards
// for the legal moves a knight can make. This is intended to be used
// as a precomputed index for detecting legal knight moves.
func LegalKnightMovesDict() [64]uint64 {
	var legalMoves [64]uint64
	// For each square on the board, prepopulate the legal moves.
	var i uint64
	var bb uint64
	var kp uint64
	for i = 0; i < 64; i++ {
		// knight position
		kp = 1 << i
		bb = 0
		s := Square(i)
		// Check for going off the left side
		if s.Col() != 1 {
			bb = bb | (kp << 15) | (kp >> 17)
		}
		if s.Col() > 2 {
			bb = bb | (kp << 6) | (kp >> 10)
		}
		// Check for going off the right side
		if s.Col() < 7 {
			bb = bb | (kp >> 6) | (kp << 10)
		}
		if s.Col() != 8 {
			bb = bb | (kp >> 15) | (kp << 17)
		}
		legalMoves[i] = bb
	}
	return legalMoves
}
