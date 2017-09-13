// An implementation of a Pawn.
package game

type Pawn struct {
	*BasePiece
}

func (p *Pawn) LegalMoves() []Square {
	return p.PawnMoves(p.Board().PieceSet[p])
}

func (p *Pawn) String() string {
	return "P"
}

func (p *Pawn) Graphic() string {
	switch p.Color() {
	case BLACK:
		return "♙"
	case WHITE:
		return "♟"
	}
	return ""
}

func (p *Pawn) Value() float64 {
	return 1.0
}
