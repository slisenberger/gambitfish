// An implementation of a Rook
package game

type Rook struct {
	*BasePiece
	HasMoved bool
}

func (p *Rook) LegalMoves() []Move {
	return ColumnAndRowMoves(p, p.Board().PieceSet[p])
}

func (p *Rook) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
}

func (p *Rook) String() string {
	return "R"
}

func (p *Rook) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♖"
	case WHITE:
		return "♜"
	}
	return ""
}

func (p *Rook) ApplyMove(m Move) {
	p.HasMoved = true
}

func (p *Rook) Value() float64 {
	return 5.0
}
