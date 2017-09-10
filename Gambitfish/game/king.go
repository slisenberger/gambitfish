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
	return "K"
}

func (p *King) Graphic() string {
	switch p.color {
	case BLACK:
		return "♔"
	case WHITE:
		return "♚"
	}
	return ""
}

func (p *King) ApplyMove(m Move) {
	p.HasMoved = true
}

func (p *King) Value() float64 {
	return 100.0
}
