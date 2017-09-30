// An implementation of a Bishop
package game

type Bishop struct {
	C Color
}

var BISHOP_DIRS = []Direction{NE, NW, SE, SW}

func BishopMoves(b *Board, p Piece, s Square) []Move {
	return RayMoves(b, p, s, BISHOP_DIRS)
}

func BishopAttackBitboard(b *Board, cur Square) uint64 {
	return RayAttackBitboard(b, cur, BISHOP_DIRS)
}
func (p *Bishop) Value() float64 {
	return 3.0
}
