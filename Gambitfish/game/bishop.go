// An implementation of a Bishop
package game

type Bishop struct {
	*BasePiece
}

func (p *Bishop) LegalMoves() []Square {
	return p.DiagonalMoves(p.board.PieceSet[p])
}

func (p *Bishop) String() string {
	return "B"
}

func (p *Bishop) Graphic() string {
	switch p.color {
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
