package game

// PieceSquareEvaluator consults piece square tables
type PieceSquareEvaluator struct{}

// Table Definitions
// Note that because Square A1 is index 0, these tables if accessed
// by square index are a flipped image of what the board looks like
// from white's perspective.

var PAWN_VALUE_TABLE = [64]float64{
	0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
	-.31, 0.08, -.07, -.37, -.36, -.14, .03, -.31,
	-.22, 0.09, .05, -.11, -.10, -.02, .03, -.19,
	-.26, .03, .1, 0.09, .06, .01, 0.0, -.23,
	-.17, .16, -.02, .15, .15, 0.0, .15, -.13,
	.07, .29, .21, .44, .4, .31, .44, .07,
	.78, .83, .86, .73, 1.02, .82, .85, .90,
	0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
}

var KNIGHT_VALUE_TABLE = [64]float64{
	-.74, -.23, -.26, -.24, -.19, -.35, -.22, -.69,
	-.23, -.15, .02, 0.0, .02, 0.0, -.23, -.2,
	-.18, .1, .13, .22, .18, .15, .11, -.14,
	-.01, .05, .31, .21, .22, .35, .02, 0.0,
	.24, .24, .45, .37, .33, .41, .25, .17,
	.1, .67, .01, .74, .73, .27, .62, -.02,
	-.03, -.06, 1, -.36, .04, .62, -.04, -.14,
	-.66, -.53, -.75, -.75, -.1, -.55, -.58, -.7,
}

var ROOK_VALUE_TABLE = [64]float64{
	-0.3, -0.24, -0.18, 0.05, -.02, -0.18, -0.31, -0.32,
	-0.53, -0.38, -0.31, -0.26, -0.29, -0.43, -0.44, -0.53,
	-0.42, -.28, -0.42, -0.25, -0.25, -0.35, -0.26, -0.46,
	-0.28, -.35, -0.16, -0.21, -0.13, -0.29, -0.46, -0.3,
	0.0, 0.05, 0.16, 0.13, 0.18, -0.04, -0.09, -0.06,
	0.19, 0.35, 0.28, 0.33, 0.45, 0.27, 0.25, 0.15,
	0.55, 0.29, 0.56, 0.67, 0.55, 0.62, 0.34, 0.60,
	0.35, 0.29, 0.33, 0.04, 0.37, 0.33, 0.56, 0.5,
}

var BISHOP_VALUE_TABLE = [64]float64{
	-.07, .02, -.15, -.12, -.14, -.15, -.1, -.1,
	.19, .20, .11, .06, .07, .06, .20, .16,
	.14, .25, .24, .15, 0.08, .25, .20, .15,
	.13, .10, .17, .23, .17, .16, 0.0, .07,
	.25, .17, .2, .34, .26, .25, .15, .1,
	-0.09, .39, -3.2, .41, .52, -.1, .28, -.14,
	-.11, .2, .35, -.42, -.39, .31, .02, -.22,
	-.59, -.78, -.82, -.76, -.23, -1.07, -.37, -.5,
}

var QUEEN_VALUE_TABLE = [64]float64{
	-.39, -.3, -.31, -.13, -.31, -.36, -.34, -.42,
	-.36, -.18, 0.0, -.19, -.15, -.15, -.21, -.38,
	-.3, -.06, -.13, -.11, -.16, -.11, -.16, -.27,
	-.14, -.15, -.02, -.05, -.01, -.1, -.2, -.22,
	.01, -.16, .22, .17, .25, .2, -.13, -.06,
	-.02, .43, .32, .60, .72, .63, .43, .02,
	.14, .32, .6, -.1, .2, .76, -.57, .24,
	.06, .01, -0.08, -1.04, .69, .24, .88, .26,
}

var KING_VALUE_TABLE = [64]float64{
	.17, .3, -.02, -.14, .06, -.01, .4, .18,
	-.04, .03, -.14, -.5, -.57, -.18, .13, .04,
	-.47, -.42, -.43, -.79, -.64, -.32, -.29, -.32,
	-.55, -.43, -.52, -.28, -.51, -.47, -0.08, -.5,
	-.55, .5, .11, -.04, -.19, .13, 0.0, -.49,
	-.62, .12, -.57, .44, -.67, .28, .37, -.31,
	-.31, -1, .55, .56, .56, .55, .1, .03,
	.04, .54, .47, -.99, -.99, .6, .83, -.62,
}

// Evaluates a board based on where pieces are located, as referenced
// in the piece value tables.
func (m PieceSquareEvaluator) Evaluate(b *Board) float64 {
	eval := 0.0
	for i, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		sq := Square(i)
		if p.Color() == BLACK {
			sq = GetSquare(9-sq.Row(), sq.Col())
		}

		// We need to change our index for black since their board
		// is mirrored.
		switch p {
		case WHITEPAWN:
				eval += PAWN_VALUE_TABLE[sq]
		case BLACKPAWN:
				eval -= PAWN_VALUE_TABLE[sq]
		case WHITEKNIGHT:
				eval += KNIGHT_VALUE_TABLE[sq]
		case BLACKKNIGHT:
				eval -= KNIGHT_VALUE_TABLE[sq]
		case WHITEROOK:
				eval += ROOK_VALUE_TABLE[sq]
		case BLACKROOK:
				eval -= ROOK_VALUE_TABLE[sq]
		case WHITEBISHOP:
				eval += BISHOP_VALUE_TABLE[sq]
		case BLACKBISHOP:
				eval -= BISHOP_VALUE_TABLE[sq]
		case WHITEQUEEN:
				eval += QUEEN_VALUE_TABLE[sq]
		case BLACKQUEEN:
				eval -= QUEEN_VALUE_TABLE[sq]
		case WHITEKING:
				eval += KING_VALUE_TABLE[sq]
		case BLACKKING:
				eval -= KING_VALUE_TABLE[sq]
		}
	}
	return float64(b.Active) * eval
}
