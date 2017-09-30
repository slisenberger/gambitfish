// An implementation of a Rook
package game

type Rook struct {
	C  Color
	KS bool
	QS bool
}

var ROOK_DIRS = []Direction{N, S, E, W}

func RookMoves(b *Board, p Piece, s Square) []Move {
	return RayMoves(b, p, s, ROOK_DIRS)
}

func RookAttackBitboard(b *Board, cur Square) uint64 {
	return RayAttackBitboard(b, cur, ROOK_DIRS)
}
