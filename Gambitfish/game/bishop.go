// An implementation of a Bishop
package game

type Bishop struct {
	C Color
}

var BISHOP_DIRS = []Direction{NE, NW, SE, SW}

func BishopMoves(b *Board, p Piece, s Square) []Move {
	var moves []Move
	for _, dir := range BISHOP_DIRS {
		m := RayMoves(b, p, s, dir)
		moves = append(moves, m...)
	}
	return moves
}

func BishopAttackBitboard(b *Board, cur Square) uint64 {
	var res uint64
	res = 0
	pos := b.Position
	for _, dir := range BISHOP_DIRS {
		// Get the ray attacks in a direction for this square.
		ra := RAY_ATTACKS[dir][cur]
		blocker := ra & pos.Occupied
		if blocker > 1 {
			blockSquare := BitScan(blocker, dir > 0)
			ra = ra ^ RAY_ATTACKS[dir][blockSquare]
		}
		res = res | ra
	}
	return res
}
func (p *Bishop) Value() float64 {
	return 3.0
}
