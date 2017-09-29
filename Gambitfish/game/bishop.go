// An implementation of a Bishop
package game

type Bishop struct {
	*BasePiece
}

var BISHOP_DIRS = []Direction{NE, NW, SE, SW}

func (p *Bishop) LegalMoves(b *Board) []Move {
	var moves []Move
	for _, dir := range BISHOP_DIRS {
		m := RayMoves(b, p, b.PieceSet[p], dir)
		moves = append(moves, m...)
	}
	return moves
	// TESTING A NEW WAY OF IMPLEMENTING THIS.
	//return DiagonalMoves(p, p.Board().PieceSet[p])
}

func (p *Bishop) AttackBitboard(b *Board, cur Square) uint64 {
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

func (p *Bishop) String() string {
	return "B"
}

func (p *Bishop) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♗"
	case WHITE:
		return "♝"
	}
	return ""
}

func (p *Bishop) Value() float64 {
	return 3.0
}

func (p *Bishop) Type() PieceType {
	return BISHOP
}
