// An implementation of a Knight
package game

type Knight struct{}

func KnightMoves(b *Board, p Piece, s Square) []EfficientMove {
	return LegalKnightMoves(b, p, s)
}

// AttackBitboard returns the bitboard for the pieces under attack by this king.
func (p *Knight) AttackBitboard(b *Board, cur Square) uint64 {
	km := LEGALKNIGHTMOVES[cur]
	return km
}

func (p *Knight) Value() float64 {
	return 3.0
}
