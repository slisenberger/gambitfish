// An implementation of a King.
package game

type King struct {
	*BasePiece
	HasMoved bool
}

func (p *King) LegalMoves() []Move {
	return KingMoves(p, p.Board().PieceSet[p])
}

func (p *King) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
}

func (p *King) String() string {
	return "K"
}

func (p *King) Graphic() string {
	switch p.Color() {
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
