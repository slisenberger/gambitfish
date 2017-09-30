// An implementation of a Rook
package game

type Rook struct {
	C  Color
	KS bool
	QS bool
}

var ROOK_DIRS = []Direction{N, S, E, W}

func RookMoves(b *Board, p Piece, s Square) []Move {
	moves := []Move{}
	for _, dir := range ROOK_DIRS {
		m := RayMoves(b, p, s, dir)
		moves = append(moves, m...)
	}
	return moves
}

func RookAttackBitboard(b *Board, cur Square) uint64 {
	var res uint64
	res = 0
	pos := b.Position
	for _, dir := range ROOK_DIRS {
		// Get the ray attacks in a direction for this square.
		ra := RayAttacks(dir, cur)
		blocker := ra & pos.Occupied
		if blocker > 1 {
			blockSquare := BitScan(blocker, dir > 0)
			ra = ra ^ RayAttacks(dir, blockSquare)
		}
		res = res | ra
	}
	return res
}
