// An implementation of a Queen.
package game

type Queen struct {
	C Color
}

var QUEEN_DIRS = []Direction{N, S, E, W, NW, NE, SW, SE}

func QueenMoves(b *Board, p Piece, s Square) []EfficientMove {
	return RayMoves(b, p, s, true, true)
}

func QueenAttackBitboard(b *Board, cur Square) uint64 {
	return RayAttackBitboard(b, cur, true, true)
}
