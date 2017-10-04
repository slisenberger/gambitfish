package evaluate

import "../../game"

const FIRST_ROW_PAWN_SHIELD_VALUE = .35
const SECOND_ROW_PAWN_SHIELD_VALUE = .15

type KingSafetyEvaluator struct{}

// Evaluate returns an estimate of the positional safety of a king
// according to its pawns.
func (k KingSafetyEvaluator) Evaluate(b *game.Board) float64 {
	eval := 0.0
	// Find the king.
	wkbb := b.Position.WhiteKing
	wkS := game.SquaresFromBitBoard(wkbb)[0]
	bkbb := b.Position.BlackKing
	if bkbb == 0 {
		b.Print()
		panic("help")
	}
	bkS := game.SquaresFromBitBoard(bkbb)[0]
	var wkShield uint64
	var bkShield uint64
	// Identify the king's pawn shields.
	if wkS.Col() == 1 {
		wkShield = wkbb<<8 | wkbb<<9
	} else if wkS.Col() == 8 {
		wkShield = wkbb<<7 | wkbb<<8

	} else {
		wkShield = wkbb<<7 | wkbb<<8 | wkbb<<9
	}
	if bkS.Col() == 1 {
		bkShield = bkbb>>8 | bkbb>>7
	} else if bkS.Col() == 8 {
		bkShield = bkbb>>9 | bkbb>>8
	} else {
		bkShield = bkbb>>7 | bkbb>>8 | bkbb>>9
	}

	// Find and evaluate the pawns in the pawn shields.
	wp := game.SquaresFromBitBoard(wkShield & b.Position.WhitePawns)
	bp := game.SquaresFromBitBoard(bkShield & b.Position.BlackPawns)

	eval += FIRST_ROW_PAWN_SHIELD_VALUE * float64(len(wp))
	eval -= FIRST_ROW_PAWN_SHIELD_VALUE * float64(len(bp))

	// Now find the second rank above the king
	if wkS.Col() == 1 {
		wkShield = wkbb<<16 | wkbb<<17
	} else if wkS.Col() == 8 {
		wkShield = wkbb<<15 | wkbb<<16
	} else {
		wkShield = wkbb<<15 | wkbb<<16 | wkbb<<17
	}
	if bkS.Col() == 1 {
		bkShield = bkbb>>16 | bkbb>>15
	} else if bkS.Col() == 8 {
		bkShield = bkbb>>17 | bkbb>>16
	} else {
		bkShield = bkbb>>15 | bkbb>>16 | bkbb>>17
	}
	// Find and evaluate the pawns one rank in front of the king.
	wp = game.SquaresFromBitBoard(wkShield & b.Position.WhitePawns)
	bp = game.SquaresFromBitBoard(bkShield & b.Position.BlackPawns)

	eval += SECOND_ROW_PAWN_SHIELD_VALUE * float64(len(wp))
	eval -= SECOND_ROW_PAWN_SHIELD_VALUE * float64(len(bp))

	// Find and evaluate king column
	// This is really crude. Basically trying to promote castling if at all possible.
	switch wkS.Col() {
	case 1, 2, 3, 7, 8:
		eval += .5
	case 4, 5:
		eval -= .2
	}
	switch bkS.Col() {
	case 1, 2, 3, 7, 8:
		eval -= .5
	case 4, 5:
		eval += .2
	}

	// Find the pawn shields two rows in front of the king.
	// TODO

	// Return, making sure we are color appropriate.
	return float64(b.Active) * eval
}
