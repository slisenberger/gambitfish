// An implementation of a Pawn.
package game

type Pawn struct {
	*BasePiece
}

func (p *Pawn) LegalMoves() []Move {
	return PawnMoves(p, p.Board().PieceSet[p])
}

func (p *Pawn) Attacking() []Square {
	return PawnAttackingSquares(p, p.Board().PieceSet[p])
}

func (p *Pawn) AttackBitboard(cur Square) uint64 {
	//TODO(slisenberger): THIS IS NOT IMPLEMENTED.
	return 0
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

func (p *Pawn) Type() PieceType {
	return PAWN
}
