// An implementation of a Rook
package game

type Rook struct {
	*BasePiece
	KS bool
	QS bool
}

var ROOK_DIRS = []Direction{N, S, E, W}

func (p *Rook) LegalMoves() []Move {
	moves := []Move{}
	for _, dir := range ROOK_DIRS {
		m := RayMoves(p, p.Board().PieceSet[p], dir)
		moves = append(moves, m...)
	}
	return moves
	// TODO(slisenberger): trying something new with bitboards,
	// clean this up or revert..
	// return ColumnAndRowMoves(p, p.Board().PieceSet[p])
}

func (p *Rook) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
}

func (p *Rook) AttackBitboard(cur Square) uint64 {
	var res uint64
	res = 0
	pos := p.Board().Position
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
	switch p.Color() {
	case WHITE:
		res = res & pos.BlackPieces
	case BLACK:
		res = res & pos.WhitePieces
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
