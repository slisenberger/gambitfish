// An implementation of a Bishop
package game

type Bishop struct {
	*BasePiece
}

func (p *Bishop) LegalMoves() []Move {
	return DiagonalMoves(p, p.Board().PieceSet[p])
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
