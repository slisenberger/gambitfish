// An implementation of a Pawn.
package game

type Pawn struct {
	*BasePiece
}

func (p *Pawn) LegalMoves(b *Board) []Move {
	return PawnMoves(b, p, b.PieceSet[p])
}

func (p *Pawn) AttackBitboard(b *Board, cur Square) uint64 {
	var res uint64
	res = 0
	switch p.Color() {
	case WHITE:
		res = WHITEPAWNATTACKS[cur]
	case BLACK:
		res = BLACKPAWNATTACKS[cur]
	}
	return res
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
