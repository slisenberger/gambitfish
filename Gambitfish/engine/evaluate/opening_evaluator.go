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
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, .1, .2, .2, .1, 0, 0,
	0, 0, .1, .2, .2, .1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// Bitboards representing starting configurations.
var bStartBishops = uint64(0xFF00000000000000)
var wStartBishops = uint64(0x00000000000000FF)
var bStartKnights = uint64(0xFF00000000000000)
var wStartKnights = uint64(0x00000000000000FF)

// Evaluates a board in a couple ways that encourage good opening play:
// 1. Gives a slight bonus for the number of pieces attacking the center.
// 2. Gives a slight negative for the number of knights/bishops on their
// starting squares.
func (m OpeningEvaluator) Evaluate(b *game.Board) float64 {
	res := 0.0
	// For every piece, calculate its attacks, and add the values.
	for cur, p := range b.Squares {
		if p == game.NULLPIECE {
			continue
		}
		for _, s := range game.SquaresFromBitBoard(game.AttackBitboard(b, p, game.Square(cur))) {
			if p.Color() == game.WHITE {
				res += centerAttackingWeights[s]
			} else {
				res -= centerAttackingWeights[s]
			}
		}
	}
	// Penalize knights and bishops on starting squares.
	var startPenalty = .35
	if (wStartBishops & b.Position.WhiteBishops) > 0 {
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
	return float64(b.Active) * res
}
