// An implementation of a Pawn.
package game

type Pawn struct {
	*BasePiece
}

func (p *Pawn) LegalMoves() []Square {
	return p.PawnMoves(p.board.PieceSet[p])
}

func (p *Pawn) String() string {
	switch p.color {
	case WHITE:
		return "P"
	case BLACK:
		return "p"
	}
	return ""
}

func (p *Pawn) Value() float64 {
	return 1.0
}
