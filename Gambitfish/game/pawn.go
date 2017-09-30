// An implementation of a Pawn.
package game

// CURRENTLY UNCALLED
func PawnAttackBitboard(b *Board, p Piece, cur Square) uint64 {
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
