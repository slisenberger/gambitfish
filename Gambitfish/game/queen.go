// An implementation of a Queen.
package game

type Queen struct {
	*BasePiece
}

func (p *Queen) LegalMoves() []Move {
	moves := DiagonalMoves(p, p.Board().PieceSet[p])
	columnMoves := ColumnAndRowMoves(p, p.Board().PieceSet[p])
	return append(moves, columnMoves...)
}

func (p *Queen) Attacking() []Square {
	moves := p.LegalMoves()
	squares := make([]Square, len(moves))
	for i, move := range moves {
		squares[i] = move.Square
	}
	return squares
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
