// An implementation of a King.
package game

type King struct {
	C Color
}

func KingMoves(b *Board, p Piece, s Square) []EfficientMove {
	moves := LegalKingMoves(b, p, s)
	moves = append(moves, CastlingMoves(b, p, s)...)
	return moves
}
