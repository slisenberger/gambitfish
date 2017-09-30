// An implementation of a Queen.
package game

type Queen struct {
	C Color
}

var QUEEN_DIRS = []Direction{N, S, E, W, NW, NE, SW, SE}

func QueenMoves(b *Board, p Piece, s Square) []Move {
	moves := []Move{}
	for _, dir := range QUEEN_DIRS {
		m := RayMoves(b, p, s, dir)
		moves = append(moves, m...)
	}
	return moves
}

func QueenAttackBitboard(b *Board, cur Square) uint64 {
	var res uint64
	res = 0
	pos := b.Position
	for _, dir := range QUEEN_DIRS {
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
