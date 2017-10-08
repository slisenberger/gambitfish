package game

// PieceSquareEvaluator consults piece square tables
type PieceSquareEvaluator struct{}

// Table Definitions
// Note that because Square A1 is index 0, these tables if accessed
// by square index are a flipped image of what the board looks like
// from white's perspective.

var PAWN_VALUE_TABLE = [64]float64{
	0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
	.05, .10, .10, 0.0, 0.0, .10, .10, .05,
	.05, .05, -.1, .15, .15, -.1, .05, .05,
	.00, 0.0, 0.0, .30, .30, 0.0, 0.0, 0.0,
	.05, .05, .05, .33, .33, .05, .05, .05,
	.10, .10, .20, .40, .40, .10, .10, .10,
	.40, .60, .60, .60, .60, .60, .60, .40,
	0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
}

var KNIGHT_VALUE_TABLE = [64]float64{
	-.6, -.35, -.2, -.2, -.2, -.2, -.35, -.6,
	-.4, 0.0, 0.05, 0.05, 0.05, 0.05, 0.0, -.4,
	-.3, 0.1, 0.2, 0.2, 0.2, 0.2, 0.1, -.3,
	-.1, 0.1, 0.2, 0.4, 0.4, 0.2, 0.1, -.1,
	-.1, 0.1, 0.2, 0.35, 0.35, 0.2, 0.1, -.1,
	-.2, 0.0, 0.15, 0.2, 0.2, 0.15, 0.0, -.2,
	-.2, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, -.2,
	-.4, -.2, -.2, -.2, -.2, -.2, -.2, -.4,
}

// Evaluates a board based on where pieces are located, as referenced
// in the piece value tables.
func (m PieceSquareEvaluator) Evaluate(b *Board) float64 {
	eval := 0.0
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		sq := Square(s)

		c := p.Color()
		// We need to change our index for black since their board
		// is mirrored.
		if c == BLACK {
			sq = GetSquare(9-sq.Row(), sq.Col())
		}
		switch p.Type() {
		case PAWN:
			if c == WHITE {
				eval += PAWN_VALUE_TABLE[sq]
			} else {
				eval -= PAWN_VALUE_TABLE[sq]
			}
		case KNIGHT:
			if c == WHITE {
				eval += KNIGHT_VALUE_TABLE[sq]
			} else {
				eval -= KNIGHT_VALUE_TABLE[sq]
			}
		}
	}
	return float64(b.Active) * eval
}
