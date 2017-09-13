// An implementation of a Queen.
package game

type Queen struct {
	*BasePiece
}

func (p *Queen) LegalMoves() []Square {
	moves := p.DiagonalMoves(p.Board().PieceSet[p])
	columnMoves := p.ColumnAndRowMoves(p.Board().PieceSet[p])
	return append(moves, columnMoves...)
}

func (p *Queen) String() string {
	return "Q"
}

func (p *Queen) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♕"
	case WHITE:
		return "♛"
	}
	return ""
}

func (p *Queen) Value() float64 {
	return 9.0
}
