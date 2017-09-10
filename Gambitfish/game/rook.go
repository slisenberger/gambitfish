// An implementation of a Rook
package game

type Rook struct {
	*BasePiece
	HasMoved bool
}

func (p *Rook) LegalMoves() []Square {
	return p.ColumnAndRowMoves(p.board.PieceSet[p])
}

func (p *Rook) String() string {
	return "R"
}

func (p *Rook) Graphic() string {
	switch p.color {
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
