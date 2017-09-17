// An implementation of a Knight
package game

type Knight struct {
	*BasePiece
}

func (p *Knight) LegalMoves() []Move {
	return KnightMoves(p, p.Board().PieceSet[p])
}

func (p *Knight) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
}

func (p *Knight) String() string {
	return "N"
}

func (p *Knight) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♘"
	case WHITE:
		return "♞"
	}
	return ""
}

func (p *Knight) Value() float64 {
	return 3.0
}
