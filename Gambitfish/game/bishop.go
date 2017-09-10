// An implementation of a Bishop
package game

type Bishop struct {
	*BasePiece
}

func (p *Bishop) LegalMoves() []Square {
	return p.DiagonalMoves(p.board.PieceSet[p])
}

func (p *Bishop) String() string {
	switch p.color {
	case WHITE:
		return "B"
	case BLACK:
		return "b"
	}
	return ""
}

func (p *Bishop) Value() float64 {
	return 3.0
}
