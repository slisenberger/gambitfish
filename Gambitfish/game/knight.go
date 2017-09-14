// An implementation of a Knight
package game

type Knight struct {
	*BasePiece
}

func (p *Knight) LegalMoves() []Move {
	return KnightMoves(p, p.Board().PieceSet[p])
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
