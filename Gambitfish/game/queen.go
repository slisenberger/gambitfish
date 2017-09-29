// An implementation of a Queen.
package game

type Queen struct {
	*BasePiece
}

var QUEEN_DIRS = []Direction{N, S, E, W, NW, NE, SW, SE}

func (p *Queen) LegalMoves(b *Board) []Move {
	moves := []Move{}
	for _, dir := range QUEEN_DIRS {
		m := RayMoves(b, p, b.PieceSet[p], dir)
		moves = append(moves, m...)
	}
	return moves
	// TODO(slisenberger): testing bitboard, clean this up if it works.
	// TESTING SOMETHING NEW...
	//moves := DiagonalMoves(p, p.Board().PieceSet[p])
	//columnMoves := ColumnAndRowMoves(p, p.Board().PieceSet[p])
	//return append(moves, columnMoves...)
}

func (p *Queen) AttackBitboard(b *Board, cur Square) uint64 {
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

func (p *Queen) String() string {
	return "Q"
}

func (p *Queen) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♕"
	case WHITE:
		return "♛"
	}
	return ""
}

func (p *Queen) Value() float64 {
	return 9.0
}

func (p *Queen) Type() PieceType {
	return QUEEN
}
