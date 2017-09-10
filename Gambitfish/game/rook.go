// An implementation of a Rook
package game

type Rook struct {
	*BasePiece
	HasMoved bool
}

func (p *Rook) LegalMoves() []Square {
	return p.ColumnAndRowMoves()
}

func (p *Rook) String() string {
	switch p.color {
	case WHITE:
		return "R"
	case BLACK:
		return "r"
	}
	return ""
}

func (p *Rook) ApplyMove(m Move) {
	p.HasMoved = true
}
