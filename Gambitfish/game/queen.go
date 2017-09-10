// An implementation of a Queen.
package game

type Queen struct {
	*BasePiece
}

func (p *Queen) LegalMoves() []Square {
	moves := p.DiagonalMoves(p.board.PieceSet[p])
	columnMoves := p.ColumnAndRowMoves(p.board.PieceSet[p])
	return append(moves, columnMoves...)
}

func (p *Queen) String() string {
	return "Q"
}

func (p *Queen) Graphic() string {
	switch p.color {
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
