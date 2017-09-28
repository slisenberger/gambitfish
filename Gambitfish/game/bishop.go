// An implementation of a Bishop
package game

type Bishop struct {
	*BasePiece
}

func (p *Bishop) LegalMoves() []Move {
	moves := []Move{}
	dirs := []Direction{NE, NW, SE, SW}
	for _, dir := range dirs {
		moves = append(moves, RayMoves(p, p.Board().PieceSet[p], dir)...)
	}
	return moves
	// TESTING A NEW WAY OF IMPLEMENTING THIS.
	// return DiagonalMoves(p, p.Board().PieceSet[p])
}

func (p *Bishop) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
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
