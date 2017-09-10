// An implementation of a Knight
package game

type Knight struct {
	*BasePiece
}

func (p *Knight) LegalMoves() []Square {
	return p.KnightMoves(p.board.PieceSet[p])
}

func (p *Knight) String() string {
	switch p.color {
	case WHITE:
		return "N"
	case BLACK:
		return "n"
	}
	return ""
}

func (p *Knight) Value() float64 {
	return 3.0
}
