// An implementation of a Queen.
package game

type Queen struct {
	*BasePiece
}

func (p *Queen) LegalMoves() []Square {
	moves := p.DiagonalMoves()
	columnMoves := p.ColumnAndRowMoves()
	return append(moves, columnMoves...)
}

func (p *Queen) String() string {
	switch p.color {
	case WHITE:
		return "Q"
	case BLACK:
		return "q"
	}
	return ""
}
