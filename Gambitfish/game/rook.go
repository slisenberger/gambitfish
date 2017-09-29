// An implementation of a Rook
package game

type Rook struct {
	C  Color
	KS bool
	QS bool
}

var ROOK_DIRS = []Direction{N, S, E, W}

func (p *Rook) LegalMoves(b *Board) []Move {
	moves := []Move{}
	for _, dir := range ROOK_DIRS {
		m := RayMoves(b, p, b.PieceSet[p], dir)
		moves = append(moves, m...)
	}
	return moves
}

func (p *Rook) AttackBitboard(b *Board, cur Square) uint64 {
	var res uint64
	res = 0
	pos := b.Position
	for _, dir := range ROOK_DIRS {
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

func (p *Rook) String() string {
	return "R"
}

func (p *Rook) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♖"
	case WHITE:
		return "♜"
	}
	return ""
}

func (p *Rook) Value() float64 {
	return 5.0
}

func (p *Rook) Type() PieceType {
	return ROOK
}
