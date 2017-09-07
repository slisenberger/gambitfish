// An implementation of a King.
package game

type King struct {
	*BasePiece
	HasMoved bool
}

func (p *King) LegalMoves() []Square {
	return p.KingMoves()
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

func (p *King) ApplyMove(m Move) Piece {
	p.square = &m.square
	p.HasMoved = true
	return p
}
