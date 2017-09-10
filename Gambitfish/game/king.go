// An implementation of a King.
package game

type King struct {
	*BasePiece
	HasMoved bool
}

func (p *King) LegalMoves() []Square {
	return p.KingMoves(p.board.PieceSet[p])
}

func (p *King) String() string {
	switch p.color {
	case WHITE:
		return "K"
	case BLACK:
		return "k"
	}
	return ""
}

func (p *King) ApplyMove(m Move) {
	p.HasMoved = true
}

func (p *King) Value() float64 {
	return 100.0
}
