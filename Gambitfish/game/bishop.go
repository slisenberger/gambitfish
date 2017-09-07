// An implementation of a Bishop
package game

type Bishop struct {
	*BasePiece
}

func (p *Bishop) LegalMoves() []Square {
	return p.DiagonalMoves()
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

func (p *Bishop) ApplyMove(m Move) Piece {
	p.square = &m.square
	return p
}
