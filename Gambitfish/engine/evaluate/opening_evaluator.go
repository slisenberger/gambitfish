// Opening evaluator evaluates a position in a way that encourages
// development in the opening.
package evaluate

import "../../game"

type OpeningEvaluator struct{}

// centerAttackingWeights is an array of the bonus an attacked square
// provides from being in the center.
var centerAttackingWeights = []float64{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, .1, .1, 0, 0, 0,
	0, 0, .1, .35, .35, .1, 0, 0,
	0, 0, .1, .35, .35, .1, 0, 0,
	0, 0, 0, .1, .1, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// Bitboards representing starting configurations.
var wStartBishops = uint64(0x4200000000000000)
var bStartBishops = uint64(0x0000000000000042)
var wStartKnights = uint64(0x2400000000000000)
var bStartKnights = uint64(0x0000000000000024)

// Evaluates a board in a couple ways that encourage good opening play:
// 1. Gives a slight bonus for the number of pieces attacking the center.
// 2. Gives a slight negative for the number of knights/bishops on their
// starting squares.
func (m OpeningEvaluator) Evaluate(b *game.Board) float64 {
	res := 0.0
	// For every piece, calculate its attacks, and add the values.
	for p, cur := range b.PieceSet {
		for _, s := range game.SquaresFromBitBoard(p.AttackBitboard(b, cur)) {
			if p.Color() == b.Active {
				res += centerAttackingWeights[s]
			} else {
				res -= centerAttackingWeights[s]
			}
		}
	}
	// Penalize knights and bishops on starting squares.
	var startPenalty = .3
	switch b.Active {
	case game.WHITE:
		if wStartBishops&b.Position.WhiteBishops > 0 {
			res -= startPenalty * float64(len(game.SquaresFromBitBoard(wStartBishops&b.Position.WhiteBishops)))

		}
		if wStartKnights&b.Position.WhiteKnights > 0 {
			res -= startPenalty * float64(len(game.SquaresFromBitBoard(wStartKnights&b.Position.WhiteKnights)))

		}
		if bStartBishops&b.Position.BlackBishops > 0 {
			res += startPenalty * float64(len(game.SquaresFromBitBoard(bStartBishops&b.Position.BlackBishops)))
		}
		if bStartKnights&b.Position.BlackKnights > 0 {
			res += startPenalty * float64(len(game.SquaresFromBitBoard(bStartKnights&b.Position.BlackKnights)))
		}

	case game.BLACK:
		if wStartBishops&b.Position.WhiteBishops > 0 {
			res += startPenalty * float64(len(game.SquaresFromBitBoard(wStartBishops&b.Position.WhiteBishops)))

		}
		if wStartKnights&b.Position.WhiteKnights > 0 {
			res += startPenalty * float64(len(game.SquaresFromBitBoard(wStartKnights&b.Position.WhiteKnights)))

		}
		if bStartBishops&b.Position.BlackBishops > 0 {
			res -= startPenalty * float64(len(game.SquaresFromBitBoard(bStartBishops&b.Position.BlackBishops)))
		}
		if bStartKnights&b.Position.BlackKnights > 0 {
			res -= startPenalty * float64(len(game.SquaresFromBitBoard(bStartKnights&b.Position.BlackKnights)))
		}
	}
	return res
}
